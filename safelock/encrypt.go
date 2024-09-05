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
	signals, closeSignals := sl.getExitSignals()
	start := 20.0

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
		var inputFiles []archiver.File

		if err = sl.validateEncryptionInputs(inputPaths, password); err != nil {
			errs <- fmt.Errorf("invalid encryption input > %w", err)
			return
		}

		if inputFiles, err = sl.getInputFiles(ctx, inputPaths, 5.0, start); err != nil {
			errs <- fmt.Errorf("failed to read and list input paths > %w", err)
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		calc := utils.NewPathsCalculator(start, inputFiles)
		rw := newWriter(password, output, cancel, calc, sl.EncryptionConfig, errs)
		rw.asyncGcm = newAsyncGcm(password, sl.EncryptionConfig, errs)

		if err = sl.encryptFiles(ctx, inputFiles, rw, calc); err != nil {
			errs <- fmt.Errorf("failed to create encrypted archive file > %w", err)
			return
		}

		if err = rw.WriteHeader(); err != nil {
			errs <- fmt.Errorf("failed to create encrypted file header > %w", err)
			return
		}

		sl.updateStatus("All set and encrypted!", 100.0)
		close(errs)
		closeSignals()
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

func (sl *Safelock) getExitSignals() (<-chan os.Signal, func()) {
	signals := make(chan os.Signal, 2)
	close := func() {
		signal.Stop(signals)
		close(signals)
	}

	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	return signals, close
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

func (sl *Safelock) getInputFiles(
	ctx context.Context,
	paths []string,
	start, end float64,
) (files []archiver.File, err error) {
	sl.updateStatus("Listing and preparing files ", start)

	filesMap := make(map[string]string, len(paths))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if end >= start {
					start += 1.0
					time.Sleep(time.Second / 5)
				}
			}
		}
	}()

	for _, path := range paths {
		filesMap[path] = ""
	}

	if files, err = archiver.FilesFromDisk(nil, filesMap); err != nil {
		return
	}

	return
}

func (sl *Safelock) encryptFiles(
	ctx context.Context,
	inputFiles []archiver.File,
	slWriter *safelockWriter,
	calc *utils.PercentCalculator,
) (err error) {
	go sl.updateProgressStatus(ctx, "Encrypting", calc)

	if err = sl.archive(ctx, slWriter, inputFiles); err != nil {
		return
	}

	slWriter.cancel()
	return
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
