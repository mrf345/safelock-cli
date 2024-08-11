package safelock

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mholt/archiver/v4"
	myErrs "github.com/mrf345/safelock-cli/errors"
	"github.com/mrf345/safelock-cli/internal/utils"
	"golang.org/x/crypto/pbkdf2"
)

// encrypts `inputPath` which can be either a file or directory and output into the `outputPath`
// which must be a nonexisting file filepath.
//
// NOTE: `ctx` context is optional you can pass `nil` and the method will handle it
func (sl *Safelock) Encrypt(ctx context.Context, inputPath, outputPath, password string) (err error) {
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
		var archiveFile *utils.RegFile
		var outputFile *os.File

		sl.updateStatus("Validating input and output", "0%")
		if err = validateEncryptionPaths(inputPath, outputPath); err != nil {
			errs <- fmt.Errorf("invalid encryption input/output paths > %w", err)
			return
		}

		if len(password) < sl.MinPasswordLength {
			errs <- &myErrs.ErrInvalidPassword{Len: len(password), Need: sl.MinPasswordLength}
			return
		}

		sl.updateStatus("Creating compressed archive file", "1%")
		if archiveFile, err = sl.createArchiveFile(ctx, inputPath); err != nil {
			errs <- fmt.Errorf("failed to create archive file > %w", err)
			return
		}

		if outputFile, err = os.Create(outputPath); err != nil {
			errs <- fmt.Errorf("filed to create output file > %w", err)
			return
		}

		unRegister := sl.Registry.Register(outputFile)

		sl.updateStatus("Encrypting compressed archive file", "30%")
		if err = sl.encryptAndWriteInChunks(password, archiveFile.File, outputFile); err != nil {
			errs <- fmt.Errorf("failed to encrypt file > %w", err)
			return
		}

		unRegister()
		sl.updateStatus(fmt.Sprintf("Encrypted %s", outputPath), "100%")
		sl.StatusObs.Trigger(EventStatusEnd)
		close(signals)
		close(errs)
	}()

	for {
		select {
		case <-ctx.Done():
			sl.Registry.RemoveAll()
			err = &myErrs.ErrContextExpired{}
			return
		case err = <-errs:
			return
		case <-signals:
			sl.Registry.RemoveAll()
			return
		}
	}
}

func (sl *Safelock) getExitSignalsChannel() chan os.Signal {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	return signals
}

func validateEncryptionPaths(inputPath, outputPath string) (err error) {
	inputIsFile, inputErrFile := utils.IsValidFile(inputPath)
	inputIsDir, inputErrDir := utils.IsValidDir(inputPath)

	if !inputIsFile && !inputIsDir {
		if inputErrFile != nil {
			return inputErrFile
		} else {
			return inputErrDir
		}
	}

	if _, err = os.Stat(outputPath); err == nil {
		err = &myErrs.ErrInvalidOutputPath{Path: outputPath}
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}

	return
}

func (sl *Safelock) createArchiveFile(ctx context.Context, inputPath string) (file *utils.RegFile, err error) {
	var files []archiver.File

	statusCtx, cancelStatus := context.WithCancel(ctx)
	defer cancelStatus()

	if files, err = archiver.FilesFromDisk(nil, map[string]string{inputPath: ""}); err != nil {
		err = fmt.Errorf("failed to list archive files > %w", err)
		return
	}

	if file, err = sl.Registry.NewFile("", "e_output_temp"); err != nil {
		err = fmt.Errorf("failed to create temporary file > %w", err)
		return
	}

	format := archiver.CompressedArchive{
		Compression: sl.Compression,
		Archival:    sl.Archival,
	}

	go sl.updateArchiveFileStatus(statusCtx, inputPath, file.Name(), "Creating", 1.0)

	if err = format.Archive(ctx, file, files); err != nil {
		err = fmt.Errorf("failed to create archive file > %w", err)
		return
	}

	_, err = file.Seek(0, io.SeekStart)

	return
}

func (sl *Safelock) updateArchiveFileStatus(ctx context.Context, inputPath, archivePath, act string, start float64) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if p, err := utils.GetPathsPercent(inputPath, archivePath, start, 30.0); err != nil {
				return
			} else {
				sl.updateStatus(fmt.Sprintf("%s compressed archive file", act), p)
				time.Sleep(time.Second / 4)
			}
		}
	}
}

func (sl *Safelock) encryptAndWriteInChunks(pwd string, inputFile *os.File, outputFile *os.File) (err error) {
	for err = range sl.runFilesEncryptionPipe(inputFile, outputFile, pwd) {
		if err != nil {
			os.Remove(outputFile.Name())
			return
		}
	}
	return
}

func (sl *Safelock) runFilesEncryptionPipe(inputFile *os.File, outputFile *os.File, pwd string) <-chan error {
	errs := make(chan error)

	go func() {
		size := sl.encryptionBufferSize()
		calc := &utils.ChunkPercentCalculator{File: inputFile, ChunkSize: size, Start: 30.0, Portion: 70.0}
		chunks := sl.getFileChunksChannel(inputFile, size, errs)
		encrypted := sl.getEncryptedChunksChannel(pwd, chunks, errs)
		sl.writeChunks(outputFile, "Encrypting", calc, encrypted, errs)
	}()

	return errs
}

func (sl *Safelock) getFileChunksChannel(file *os.File, chunkSize int, errs chan error) <-chan *fileChunk {
	chunks := make(chan *fileChunk, sl.ChannelSize)

	go func() {
		for {
			var sought int
			var err error

			chunk := make([]byte, chunkSize)

			if sought, err = file.Read(chunk); err != nil && err != io.EOF {
				errs <- fmt.Errorf("failed to read input file > %w", err)
				return
			} else if err == io.EOF {
				break
			}

			chunks <- &fileChunk{
				Chunk:  chunk,
				Sought: sought,
			}
		}

		close(chunks)
	}()

	return chunks
}

func (sl *Safelock) getEncryptedChunksChannel(pwd string, chunks <-chan *fileChunk, errs chan error) <-chan []byte {
	encrypted := make(chan []byte, sl.ChannelSize)

	go func() {
		for chunk := range chunks {
			var block cipher.Block
			var gcm cipher.AEAD
			var err error

			salt := make([]byte, sl.SaltLength)

			if _, err = io.ReadFull(rand.Reader, salt); err != nil {
				errs <- fmt.Errorf("failed to create random salt > %w", err)
				return
			}

			key := pbkdf2.Key([]byte(pwd), salt, sl.KeyLength, sl.IterationCount, sl.Hash)

			if block, err = aes.NewCipher(key); err != nil {
				errs <- fmt.Errorf("failed to create new cipher > %w", err)
				return
			}

			if gcm, err = cipher.NewGCM(block); err != nil {
				errs <- fmt.Errorf("failed to create new GCM > %w", err)
				return
			}

			encrypted <- append(gcm.Seal(nil, salt, chunk.Chunk[:chunk.Sought], nil), salt...)
		}

		close(encrypted)
	}()

	return encrypted
}

func (sl *Safelock) writeChunks(
	file *os.File,
	action string,
	calc *utils.ChunkPercentCalculator,
	chunks <-chan []byte,
	errs chan error,
) {
	go func() {
		for chunk := range chunks {
			if _, err := io.Copy(file, bytes.NewReader(chunk)); err != nil {
				errs <- fmt.Errorf("failed to copy to temporary file > %w", err)
				return
			}

			if percent, err := calc.GetPercent(); err != nil {
				errs <- fmt.Errorf("failed to read input file > %w", err)
				return
			} else {
				sl.updateStatus(fmt.Sprintf("%s compressed archive file", action), percent)
			}
		}

		close(errs)
	}()
}
