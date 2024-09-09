package safelock

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mrf345/safelock-cli/slErrs"
)

type InputReader interface {
	io.Reader
	io.Seeker
}

type safelockReader struct {
	io.Reader
	*safelockReaderWriterBase
	reader   InputReader
	overflow []byte
}

func newReader(
	pwd string,
	reader InputReader,
	start float64,
	cancel context.CancelFunc,
	config EncryptionConfig,
	errs chan<- error,
) safelockReader {
	return safelockReader{
		reader: reader,
		safelockReaderWriterBase: &safelockReaderWriterBase{
			pwd:    pwd,
			errs:   errs,
			cancel: cancel,
			config: config,
			start:  start,
			end:    100.0,
		},
	}
}

func (sr *safelockReader) setInputSize() (err error) {
	size, err := sr.reader.Seek(0, io.SeekEnd)
	sr.inputSize = int(size)
	return
}

func (sr *safelockReader) ReadHeader() (err error) {
	headerSize := sr.config.getHeaderSizeOut(sr.inputSize)
	sizeDiff := int64(sr.inputSize - headerSize)
	headerBytes := make([]byte, headerSize)

	if _, err = sr.reader.Seek(sizeDiff, io.SeekStart); err != nil {
		err = fmt.Errorf("can't seek header > %w", err)
		return sr.handleErr(err)
	}

	if _, err = sr.reader.Read(headerBytes); err != nil {
		err = fmt.Errorf("can't read header > %w", err)
		return sr.handleErr(err)
	}

	header := string(bytes.Trim(headerBytes, "\x00"))
	sr.blocks = strings.Split(header, ";")[1:]

	if len(sr.blocks) == 0 {
		err = &slErrs.ErrFailedToAuthenticate{Msg: "missing header content"}
		return
	}

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
		err = fmt.Errorf("invalid header block size > %w", err)
		return 0, sr.handleErr(err)
	}

	encrypted := make([]byte, blockSize)

	if read, err = sr.reader.Read(encrypted); err != nil && err != io.EOF {
		err = fmt.Errorf("cant't read encrypted chunk > %w", err)
		return read, sr.handleErr(err)
	}

	if decrypted, err = decryptChunk(encrypted, sr.pwd, read, sr.config); err != nil {
		err = fmt.Errorf("can't decrypt chunk > %w", err)
		return read, sr.handleErr(err)
	}

	sr.outputSize += len(decrypted)

	return sr.handleOverflowOut(&chunk, decrypted), nil
}

func (sr *safelockReader) handleOverflowIn(chunk *[]byte) (over int, done bool) {
	if len(sr.overflow) > 0 {
		over = copy(*chunk, sr.overflow)
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

	return
}
