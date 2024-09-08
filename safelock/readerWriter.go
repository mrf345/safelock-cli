package safelock

import (
	"context"
)

type getPercent interface {
	getCompletedPercent() float64
}

type safelockReaderWriterBase struct {
	config                EncryptionConfig
	pwd                   string
	errs                  chan<- error
	cancel                context.CancelFunc
	blocks                []string
	err                   error
	start, end            float64
	inputSize, outputSize int
}

func (srw *safelockReaderWriterBase) handleErr(err error) error {
	// panic(err)
	srw.err = err
	srw.errs <- err
	srw.cancel()
	return err
}

func (srw safelockReaderWriterBase) getCompletedPercent() float64 {
	percent := srw.start + (float64(srw.outputSize) / float64(srw.inputSize) * srw.end)

	if srw.end > percent {
		return percent
	}

	return srw.end
}

func (srw *safelockReaderWriterBase) increaseInputSize(increment int) {
	srw.inputSize += increment
}
