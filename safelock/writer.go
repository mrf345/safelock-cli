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
	writer io.Writer
}

func newWriter(
	pwd string,
	writer io.Writer,
	start float64,
	cancel context.CancelFunc,
	aead *aeadWrapper,
) safelockWriter {
	return safelockWriter{
		writer: writer,
		safelockReaderWriterBase: &safelockReaderWriterBase{
			aead:   aead,
			pwd:    pwd,
			cancel: cancel,
			start:  start,
			end:    100.0,
		},
	}
}

func (sw *safelockWriter) Write(chunk []byte) (written int, err error) {
	if written, err = sw.writer.Write(sw.aead.encrypt(chunk)); err != nil {
		err = fmt.Errorf("can't write encrypted chunk > %w", err)
		return written, sw.handleErr(err)
	}

	sw.outputSize += written
	sw.blocks = append(sw.blocks, fmt.Sprintf("%d", written))

	return
}

func (sw *safelockWriter) WriteHeader() (err error) {
	if 0 >= sw.outputSize {
		return
	}

	sw.setHeaderSize()

	header := "BS;" + strings.Join(sw.blocks, ";")
	headerBytes := make([]byte, sw.headerSize)
	headerBytes = append([]byte(header), headerBytes[len(header):]...)

	if _, err = sw.writer.Write(headerBytes); err != nil {
		err = fmt.Errorf("can't write header bytes > %w", err)
		return sw.handleErr(err)
	}

	return
}

func (sw *safelockWriter) setHeaderSize() {
	ratio := sw.aead.config.HeaderRatio
	size := sw.outputSize / ratio

	if ratio > size {
		sw.headerSize = ratio
		return
	}

	sw.headerSize = (sw.outputSize + size) / ratio
}
