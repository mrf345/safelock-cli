package safelock

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mrf345/safelock-cli/utils"
)

type safelockWriter struct {
	io.Writer
	safelockReaderWriterBase
	writer   io.Writer
	asyncGcm *asyncGcm
}

func newWriter(
	pwd string,
	writer io.Writer,
	cancel context.CancelFunc,
	calc *utils.PercentCalculator,
	config EncryptionConfig,
	errs chan<- error,
) *safelockWriter {
	return &safelockWriter{
		writer: writer,
		safelockReaderWriterBase: safelockReaderWriterBase{
			pwd:    pwd,
			calc:   calc,
			errs:   errs,
			cancel: cancel,
			config: config,
		},
	}
}

func (sw *safelockWriter) Write(chunk []byte) (written int, err error) {
	encrypted := sw.asyncGcm.encryptChunk(chunk)

	if written, err = sw.writer.Write(encrypted); err != nil {
		err = fmt.Errorf("can't write encrypted chunk > %w", err)
		return written, sw.handleErr(err)
	}

	sw.calc.OutputSize += written
	sw.blocks = append(sw.blocks, fmt.Sprintf("%d", written))

	return
}

func (sw *safelockWriter) WriteHeader() (err error) {
	sw.asyncGcm.done <- true
	sw.setSize()

	if 0 >= sw.size {
		return
	}

	header := "BS;" + strings.Join(sw.blocks, ";")
	headerBytes := make([]byte, sw.headerSize)
	headerBytes = append([]byte(header), headerBytes[len(header):]...)

	if _, err = sw.writer.Write(headerBytes); err != nil {
		err = fmt.Errorf("can't write header bytes > %w", err)
		return sw.handleErr(err)
	}

	return
}

func (sw *safelockWriter) setSize() {
	sw.size = sw.calc.OutputSize
	sw.headerSize = sw.config.getHeaderSizeIn(sw.size)
}
