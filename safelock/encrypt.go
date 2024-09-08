package safelock

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mholt/archiver/v4"
	slErrs "github.com/mrf345/safelock-cli/slErrs"
	"github.com/mrf345/safelock-cli/utils"
)

// encrypts `inputPaths` which can be either a slice of file or directory paths and then
// outputs into an object `output` that implements [io.Writer] such as [io.File]
//
// NOTE: `ctx` context is optional you can pass `nil` and the method will handle it
func (sl Safelock) Encrypt(ctx context.Context, inputPaths []string, output io.Writer, password string) (err error) {
	errs := make(chan error)
	signals, closeSignals := utils.GetExitSignals()

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
		if err = sl.validateEncryptionInputs(inputPaths, password); err != nil {
			errs <- fmt.Errorf("invalid encryption input > %w", err)
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		config := sl.EncryptionConfig
		gcm := newAsyncGcm(password, config, errs)
		writer := newWriter(password, output, 20.0, cancel, gcm, config, errs)

		if err = sl.encryptFiles(ctx, inputPaths, writer); err != nil {
			errs <- err
			return
		}

		if err = writer.WriteHeader(); err != nil {
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

func (sl Safelock) validateEncryptionInputs(inputPaths []string, pwd string) (err error) {
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

func (sl Safelock) encryptFiles(
	ctx context.Context,
	inputPaths []string,
	slWriter safelockWriter,
) (err error) {
	var files []archiver.File
	var filesMap = make(map[string]string, len(inputPaths))
	var cancelListingStatus = sl.updateListingStatus(ctx, 1.0, slWriter.start)

	for _, path := range inputPaths {
		filesMap[path] = ""
	}

	if files, err = archiver.FilesFromDisk(nil, filesMap); err != nil {
		err = fmt.Errorf("failed to read and list input paths > %w", err)
		return
	}

	go func() {
		for _, file := range files {
			slWriter.increaseInputSize(int(file.Size()))
		}

		cancelListingStatus()
		sl.updateProgressStatus(ctx, "Encrypting", slWriter)
	}()

	if err = sl.archive(ctx, &slWriter, files); err != nil {
		err = fmt.Errorf("failed to create encrypted archive file > %w", err)
		return
	}

	slWriter.cancel()
	return
}

func (sl Safelock) updateListingStatus(ctx context.Context, start, end float64) (cancel context.CancelFunc) {
	ctx, cancel = context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if start >= end {
					return
				}

				start += 1
				sl.updateStatus("Listing and preparing files", start)
				time.Sleep(time.Second / 2)
			}
		}
	}()

	return
}

func (sl Safelock) updateProgressStatus(ctx context.Context, act string, rw getPercent) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			sl.updateStatus(fmt.Sprintf("%s files", act), rw.getCompletedPercent())
			time.Sleep(time.Second / 5)
		}
	}
}
