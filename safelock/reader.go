package safelock

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/mrf345/safelock-cli/utils"
)

type InputReader interface {
	io.Reader
	io.Seeker
}

type safelockReader struct {
	io.Reader
	safelockReaderWriterBase
	reader   InputReader
	overflow []byte
}

func newReader(
	pwd string,
	reader InputReader,
	cancel context.CancelFunc,
	calc *utils.PercentCalculator,
	config EncryptionConfig,
	errs chan<- error,
) *safelockReader {
	return &safelockReader{
		reader: reader,
		safelockReaderWriterBase: safelockReaderWriterBase{
			pwd:    pwd,
			calc:   calc,
			errs:   errs,
			cancel: cancel,
			config: config,
		},
	}
}

func (sr *safelockReader) setSize() {
	sr.size = sr.calc.InputSize
	sr.setHeaderSize()
}

func (sr *safelockReader) ReadHeader() (err error) {
	sr.setSize()

	headerBytes := make([]byte, sr.headerSize)

	if _, err = sr.reader.Seek(sr.diffSize(), io.SeekStart); err != nil {
		return sr.handleErr(err)
	}

	if _, err = sr.reader.Read(headerBytes); err != nil {
		return sr.handleErr(err)
	}

	header := string(bytes.Trim(headerBytes, "\x00")[:])
	sr.blocks = strings.Split(header, ";")[1:]

	if _, err = sr.reader.Seek(0, io.SeekStart); err != nil {
		return sr.handleErr(err)
	}

	return
}

func (sr *safelockReader) Read(chunk []byte) (read int, err error) {
	var blockSize int
	var block string
	var decrypted []byte

	if over, done := sr.handleOverflowIn(&chunk); done {
		return over, nil
	}

	if len(sr.blocks) == 0 {
		if sr.err != nil {
			return 0, sr.err
		}

		return 0, io.EOF
	}

	block, sr.blocks = sr.blocks[0], sr.blocks[1:]

	if blockSize, err = strconv.Atoi(block); err != nil {
		return 0, sr.handleErr(err)
	}

	encrypted := make([]byte, blockSize)

	if read, err = sr.reader.Read(encrypted); err != nil && err != io.EOF {
		return read, sr.handleErr(err)
	}

	if decrypted, err = decryptChunk(encrypted, sr.pwd, read, sr.config); err != nil {
		return read, sr.handleErr(err)
	}

	sr.calc.OutputSize += len(decrypted)

	return sr.handleOverflowOut(&chunk, decrypted), nil
}

func (sr *safelockReader) handleOverflowIn(chunk *[]byte) (over int, done bool) {
	if len(sr.overflow) > 0 {
		over = copy(*chunk, sr.overflow[:])
		sr.overflow = sr.overflow[over:]
		left := len(*chunk) - over

		if left == 0 {
			done = true
		}
	}

	return
}

func (sr *safelockReader) handleOverflowOut(chunk *[]byte, decrypted []byte) (copied int) {
	var chunked []byte

	if len(sr.overflow) > 0 {
		chunked = append(sr.overflow, decrypted...)
	} else {
		chunked = decrypted
	}

	copied = copy(*chunk, chunked)
	sr.overflow = chunked[copied:]
	copy(*chunk, chunked[:])

	return
}
