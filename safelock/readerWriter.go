package safelock

import (
	"context"

	"github.com/mrf345/safelock-cli/utils"
)

type safelockReaderWriterBase struct {
	config EncryptionConfig
	pwd    string
	errs   chan<- error
	cancel context.CancelFunc
	calc   *utils.PercentCalculator

	size       int
	headerSize int
	blocks     []string
	err        error
}

func (srw *safelockReaderWriterBase) diffSize() int64 {
	return int64(srw.size - srw.headerSize)
}

func (srw *safelockReaderWriterBase) handleErr(err error) error {
	srw.err = err
	srw.errs <- err
	srw.cancel()
	return err
}
