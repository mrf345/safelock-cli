package safelock

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type safelockWriter struct {
	io.Writer
	*safelockReaderWriterBase
	writer   io.Writer
	asyncGcm asyncGcm
}

func newWriter(
	pwd string,
	writer io.Writer,
	start float64,
	cancel context.CancelFunc,
	asyncGcm asyncGcm,
	config EncryptionConfig,
	errs chan<- error,
) safelockWriter {
	return safelockWriter{
		writer:   writer,
		asyncGcm: asyncGcm,
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

func (sw *safelockWriter) Write(chunk []byte) (written int, err error) {
	encrypted := sw.asyncGcm.encryptChunk(chunk)

	if written, err = sw.writer.Write(encrypted); err != nil {
		err = fmt.Errorf("can't write encrypted chunk > %w", err)
		return written, sw.handleErr(err)
	}

	sw.outputSize += written
	sw.blocks = append(sw.blocks, fmt.Sprintf("%d", written))

	return
}

func (sw *safelockWriter) WriteHeader() (err error) {
	sw.asyncGcm.done <- true

	if 0 >= sw.outputSize {
		return
	}

	header := "BS;" + strings.Join(sw.blocks, ";")
	headerSize := sw.config.getHeaderSizeIn(sw.outputSize)
	headerBytes := make([]byte, headerSize)
	headerBytes = append([]byte(header), headerBytes[len(header):]...)

	if _, err = sw.writer.Write(headerBytes); err != nil {
		err = fmt.Errorf("can't write header bytes > %w", err)
		return sw.handleErr(err)
	}

	return
}
