package safelock

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/mholt/archiver/v4"
	myErrs "github.com/mrf345/safelock-cli/errors"
	"github.com/mrf345/safelock-cli/internal/utils"
	"golang.org/x/crypto/pbkdf2"
)

// decrypts `inputPath` which must be a file and outputs the content into `outputPath`
// which must be a valid path to an existing directory.
//
// NOTE: `ctx` context is optional you can pass `nil` and the method will handle it
func (sl *Safelock) Decrypt(ctx context.Context, inputPath, outputPath, password string) (err error) {
	errs := make(chan error)
	signals := sl.getExitSignalsChannel()

	if ctx == nil {
		ctx = context.Background()
	}

	sl.StatusObs.
		On(EventStatusUpdate, sl.logStatus).
		Trigger(EventStatusStart)

	defer sl.StatusObs.
		Off(EventStatusUpdate, sl.logStatus).
		Trigger(EventStatusError, err)

	go func() {
		var archiveFile *utils.TempFile

		sl.updateStatus("Validating input and output", "0%")
		if err = validateDecryptionPaths(inputPath, outputPath); err != nil {
			errs <- fmt.Errorf("invalid decryption input/output paths > %w", err)
			return
		}

		sl.updateStatus("Decrypting compressed archive file", "1%")
		if archiveFile, err = sl.decryptArchiveFileInChunks(inputPath, password); err != nil {
			errs <- fmt.Errorf("failed to decrypt archive file > %w", err)
			return
		}

		if _, err = archiveFile.Seek(0, io.SeekStart); err != nil {
			errs <- fmt.Errorf("failed to read archive file > %w", err)
			return
		}

		sl.updateStatus("Extracting compressed archive file", "70%")
		if err = sl.extractArchiveFile(ctx, outputPath, archiveFile); err != nil {
			errs <- fmt.Errorf("failed to extract archive file > %w", err)
			return
		}

		if err = archiveFile.Close(); err != nil {
			errs <- fmt.Errorf("failed to close archive file > %w", err)
			return
		}

		if err = archiveFile.Remove(); err != nil {
			errs <- fmt.Errorf("failed to remove archive file > %w", err)
			return
		}

		sl.updateStatus(fmt.Sprintf("Decrypted into %s", outputPath), "100%")
		sl.StatusObs.Trigger(EventStatusEnd)
		close(signals)
		close(errs)
	}()

	for {
		select {
		case <-ctx.Done():
			err = &myErrs.ErrContextExpired{}
			return
		case err = <-errs:
			return
		case <-signals:
			sl.TempStore.RemoveAll()
			return
		}
	}
}

func validateDecryptionPaths(inputPath, outputPath string) (err error) {
	if _, err = utils.IsValidFile(inputPath); err != nil {
		return
	}

	if _, err = utils.IsValidDir(outputPath); err != nil {
		return
	}

	return
}

func (sl *Safelock) decryptArchiveFileInChunks(inputPath, password string) (outputFile *utils.TempFile, err error) {
	var inputFile *os.File

	if inputFile, err = os.Open(inputPath); err != nil {
		err = fmt.Errorf("failed to open input file > %w", err)
		return
	}

	if outputFile, err = sl.TempStore.NewFile("", "d_output_temp"); err != nil {
		err = fmt.Errorf("failed to create temporary file > %w", err)
		return
	}

	for err = range sl.runFilesDecryptionPipe(inputFile, outputFile, password) {
		if err != nil {
			return
		}
	}

	return
}

func (sl *Safelock) runFilesDecryptionPipe(inputFile *os.File, outputFile *utils.TempFile, pwd string) <-chan error {
	errs := make(chan error)

	go func() {
		size := sl.decryptionBufferSize()
		calc := &utils.ChunkPercentCalculator{File: inputFile, ChunkSize: size, Start: 1.0, Portion: 70.0}
		chunks := sl.getFileChunksChannel(inputFile, size, errs)
		decrypted := sl.getDecryptedChunksChannel(pwd, chunks, errs)
		sl.writeChunks(outputFile.File, "Decrypting", calc, decrypted, errs)
	}()

	return errs
}

func (sl *Safelock) getDecryptedChunksChannel(pwd string, chunks <-chan *fileChunk, errs chan error) <-chan []byte {
	decrypted := make(chan []byte, sl.ChannelSize)

	go func() {
		for chunk := range chunks {
			var block cipher.Block
			var gcm cipher.AEAD
			var err error
			var data []byte

			encrypted := chunk.Chunk[:chunk.Sought-(sl.SaltLength)]
			salt := chunk.Chunk[chunk.Sought-(sl.SaltLength) : chunk.Sought]
			key := pbkdf2.Key([]byte(pwd), salt, sl.KeyLength, sl.IterationCount, sl.Hash)

			if block, err = aes.NewCipher(key); err != nil {
				errs <- fmt.Errorf("failed to create new cipher > %w", err)
				return
			}

			if gcm, err = cipher.NewGCM(block); err != nil {
				errs <- fmt.Errorf("failed to create new GCM > %w", err)
				return
			}

			if data, err = gcm.Open(nil, salt, encrypted, nil); err != nil {
				errs <- &myErrs.ErrFailedToAuthenticate{Msg: err.Error()}
				return
			}

			decrypted <- data
		}

		close(decrypted)
	}()

	return decrypted
}

func (sl *Safelock) extractArchiveFile(ctx context.Context, outputPath string, archive *utils.TempFile) (err error) {
	var reader io.ReadCloser
	var fileHandler = getArchiveFileHandler(outputPath)

	statusCtx, cancelStatus := context.WithCancel(ctx)
	defer cancelStatus()

	if reader, err = sl.Compression.OpenReader(archive); err != nil {
		return fmt.Errorf("cannot read archive file > %w", err)
	}

	go sl.updateArchiveFileStatus(statusCtx, archive.Name(), outputPath, "Extracting", 70.0)

	if err = sl.Archival.Extract(ctx, reader, nil, fileHandler); err != nil {
		return fmt.Errorf("cannot extract archive file > %w", err)
	}

	return
}

func getArchiveFileHandler(outputPath string) archiver.FileHandler {
	return func(ctx context.Context, file archiver.File) (err error) {
		var outputFile *os.File
		var reader io.ReadCloser
		var fullPath = path.Join(outputPath, file.NameInArchive)

		if file.IsDir() {
			err = os.MkdirAll(fullPath, file.Mode().Perm())
			return
		} else {
			if err = os.MkdirAll(path.Dir(fullPath), file.Mode().Perm()); err != nil {
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
