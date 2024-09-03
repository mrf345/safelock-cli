package safelock

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mholt/archiver/v4"
	slErrs "github.com/mrf345/safelock-cli/slErrs"
	"github.com/mrf345/safelock-cli/utils"
)

// encrypts `inputPaths` which can be either a slice of file or directory paths and then
// outputs into an object `output` that implements [io.Writer] such as [io.File]
//
// NOTE: `ctx` context is optional you can pass `nil` and the method will handle it
func (sl *Safelock) Encrypt(ctx context.Context, inputPaths []string, output io.Writer, password string) (err error) {
	errs := make(chan error)
	signals := sl.getExitSignalsChannel()

	if ctx == nil {
		ctx = context.Background()
	}

	sl.StatusObs.
		On(StatusUpdate.Str(), sl.logStatus).
		Trigger(StatusStart.Str())

	defer sl.StatusObs.
		Off(StatusUpdate.Str(), sl.logStatus).
		Trigger(StatusEnd.Str())

	go func() {
		var calc *utils.PercentCalculator

		if err = sl.validateEncryptionInputs(inputPaths, password); err != nil {
			errs <- fmt.Errorf("invalid encryption input > %w", err)
			return
		}

		if calc, err = utils.NewPathsCalculator(inputPaths, 20.0); err != nil {
			errs <- fmt.Errorf("failed to read input paths > %w", err)
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		rw := newWriter(password, output, cancel, calc, sl.EncryptionConfig, errs)
		rw.asyncGcm = newAsyncGcm(password, sl.EncryptionConfig, errs)

		if err = sl.encryptFiles(ctx, inputPaths, rw, calc); err != nil {
			errs <- fmt.Errorf("failed to create encrypted archive file > %w", err)
			return
		}

		if err = rw.WriteHeader(); err != nil {
			errs <- fmt.Errorf("failed to create encrypted file header > %w", err)
			return
		}

		sl.updateStatus("All set and encrypted!", 100.0)
		close(signals)
		close(errs)
	}()

	for {
		select {
		case <-ctx.Done():
			err = context.DeadlineExceeded
			return
		case err = <-errs:
			sl.StatusObs.Trigger(StatusError.Str(), err)
			return
		case <-signals:
			return
		}
	}
}

func (sl *Safelock) getExitSignalsChannel() chan os.Signal {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	return signals
}

func (sl *Safelock) validateEncryptionInputs(inputPaths []string, pwd string) (err error) {
	sl.updateStatus("Validating inputs", 0.0)

	for _, path := range inputPaths {
		if _, err = os.Stat(path); err != nil {
			return &slErrs.ErrInvalidInputPath{Path: path, Err: err}
		}
	}

	if len(pwd) < sl.MinPasswordLength {
		return &slErrs.ErrInvalidPassword{Len: len(pwd), Need: sl.MinPasswordLength}
	}

	return
}

func (sl *Safelock) encryptFiles(
	ctx context.Context,
	inputPaths []string,
	slWriter *safelockWriter,
	calc *utils.PercentCalculator,
) (err error) {
	var files []archiver.File
	var filesMap = make(map[string]string, len(inputPaths))

	for _, path := range inputPaths {
		filesMap[path] = ""
	}

	if files, err = archiver.FilesFromDisk(nil, filesMap); err != nil {
		err = fmt.Errorf("failed to list archive files > %w", err)
		return
	}

	go sl.updateProgressStatus(ctx, "Encrypting", calc)
	defer slWriter.cancel()

	return sl.archive(ctx, slWriter, files)
}

func (sl *Safelock) updateProgressStatus(ctx context.Context, act string, calc *utils.PercentCalculator) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			sl.updateStatus(fmt.Sprintf("%s files", act), calc.GetPercent())
			time.Sleep(time.Second / 5)
		}
	}
}
