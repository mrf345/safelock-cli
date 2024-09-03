package safelock

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
	slErrs "github.com/mrf345/safelock-cli/slErrs"
	"github.com/mrf345/safelock-cli/utils"
)

// decrypts `input` which must be an object that implements [io.Reader] and [io.Seeker] such as [os.File]
// and then outputs the content into `outputPath` which must be a valid path to an existing directory
//
// NOTE: `ctx` context is optional you can pass `nil` and the method will handle it
func (sl *Safelock) Decrypt(ctx context.Context, input InputReader, outputPath, password string) (err error) {
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

		if err = sl.validateDecryptionPaths(outputPath); err != nil {
			errs <- fmt.Errorf("invalid decryption input > %w", err)
			return
		}

		if calc, err = utils.NewSeekerCalculator(input, 1.0); err != nil {
			errs <- fmt.Errorf("failed to read input paths > %w", err)
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		rw := newReader(password, input, cancel, calc, sl.EncryptionConfig, errs)

		if err = rw.ReadHeader(); err != nil {
			errs <- fmt.Errorf("failed to read input file header > %w", err)
			return
		}

		if err = sl.decryptFiles(ctx, outputPath, rw, calc); err != nil {
			errs <- fmt.Errorf("failed to extract archive file > %w", err)
			return
		}

		sl.updateStatus("All set and decrypted!", 100.0)
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

func (sl *Safelock) validateDecryptionPaths(outputPath string) (err error) {
	sl.updateStatus("Validating inputs", 0.0)

	if info, err := os.Stat(outputPath); err != nil || !info.IsDir() {
		return &slErrs.ErrInvalidOutputPath{Path: outputPath, Err: err}
	}

	return
}

func (sl *Safelock) decryptFiles(
	ctx context.Context,
	outputPath string,
	slReader *safelockReader,
	calc *utils.PercentCalculator,
) (err error) {
	var reader io.ReadCloser

	if reader, err = sl.Compression.OpenReader(slReader); err != nil {
		return fmt.Errorf("cannot read archive file > %w", err)
	}

	go sl.updateProgressStatus(ctx, "Decrypting", calc)
	defer slReader.cancel()

	fileHandler := getExtractFileHandler(outputPath)

	if err = sl.Archival.Extract(ctx, reader, nil, fileHandler); err != nil {
		return fmt.Errorf("cannot extract archive file > %w", err)
	}

	return
}

func getExtractFileHandler(outputPath string) archiver.FileHandler {
	return func(ctx context.Context, file archiver.File) (err error) {
		var outputFile *os.File
		var reader io.ReadCloser
		var fullPath = filepath.Join(outputPath, file.NameInArchive)

		if file.IsDir() {
			err = os.MkdirAll(fullPath, file.Mode().Perm())
			return
		} else {
			if err = os.MkdirAll(filepath.Dir(fullPath), file.Mode().Perm()); err != nil {
				return
			}
		}

		if reader, err = file.Open(); err != nil {
			err = fmt.Errorf("failed to open within archive file > %w", err)
			return
		}
		defer reader.Close()

		if outputFile, err = os.Create(fullPath); err != nil {
			err = fmt.Errorf("failed to create decrypted file > %w", err)
			return
		}
		defer outputFile.Close()

		if _, err = io.Copy(outputFile, reader); err != nil {
			err = fmt.Errorf("failed to write decrypted file > %w", err)
			return
		}

		return
	}
}
