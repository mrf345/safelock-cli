package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mholt/archiver/v4"
	myErrs "github.com/mrf345/safelock-cli/errors"
	"github.com/mrf345/safelock-cli/utils"
	"golang.org/x/crypto/pbkdf2"
)

// TODO: add docs
// TODO: add tests
func Encrypt(inputPath, outputPath, password string, options EncryptionOptions) (err error) {
	var channel = make(chan error)
	var percent, status string

	go func() {
		var archiveFile, outputFile *os.File
		var archivePath string

		options.Status <- "Validating input and output"
		options.ProgressPercent <- "1%"
		if err = validateEncryptionPaths(inputPath, outputPath); err != nil {
			channel <- fmt.Errorf("invalid encryption input/output paths > %w", err)
			return
		}

		options.Status <- "Create archive file"
		options.ProgressPercent <- "5%"
		if archivePath, archiveFile, err = createArchiveFile(inputPath, outputPath, options); err != nil {
			channel <- fmt.Errorf("failed to create archive file > %w", err)
			return
		}

		if outputFile, err = os.Create(outputPath); err != nil {
			channel <- fmt.Errorf("filed to create output file > %w", err)
			return
		}

		options.Status <- "Encrypting input file"
		options.ProgressPercent <- "10%"
		if err = encryptAndWriteInChunks(password, archiveFile, outputFile, options); err != nil {
			err = fmt.Errorf("failed to encrypt file > %w", err)
			return
		}

		if err = outputFile.Close(); err != nil {
			return
		}

		if err = archiveFile.Close(); err != nil {
			return
		}

		if err = os.Remove(archivePath); err != nil {
			return
		}

		channel <- err
	}()

	for {
		select {
		case <-options.Context.Done():
			return &myErrs.ErrContextExpired{}
		case err = <-channel:
			status = fmt.Sprintf("Encrypted %s", outputPath)
			percent = "100%"
			options.Status <- status
			options.ProgressPercent <- percent
			utils.ClearAndPrint("%s (%s)\n", status, "100%")
			close(channel)
			return
		case status = <-options.Status:
			utils.ClearAndPrint("%s (%s)\n", status, percent)
		case percent = <-options.ProgressPercent:
			utils.ClearAndPrint("%s (%s)\n", status, percent)
		}
	}
}

func validateEncryptionPaths(iPath, oPath string) (err error) {
	inputIsFile, inputErrFile := utils.IsValidFile(iPath)
	inputIsDir, inputErrDir := utils.IsValidDir(iPath)

	if !inputIsFile && !inputIsDir {
		if inputErrFile != nil {
			return inputErrFile
		} else {
			return inputErrDir
		}
	}

	if _, err = os.Stat(oPath); err == nil {
		err = &myErrs.ErrInvalidOutputPath{Path: oPath}
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}

	return
}

func createArchiveFile(iPath, oPath string, options EncryptionOptions) (path string, file *os.File, err error) {
	var files []archiver.File

	if files, err = archiver.FilesFromDisk(nil, map[string]string{iPath: ""}); err != nil {
		err = fmt.Errorf("failed to list archive files > %w", err)
		return
	}

	if path, file, err = utils.CreateTempFile(oPath); err != nil {
		err = fmt.Errorf("failed to create temporary file > %w", err)
		return
	}

	format := archiver.CompressedArchive{
		Compression: options.Compression,
		Archival:    options.Archival,
	}

	if err = format.Archive(options.Context, file, files); err != nil {
		err = fmt.Errorf("failed to create archive file > %w", err)
		return
	}

	file.Seek(0, io.SeekStart)

	return
}

func encryptAndWriteInChunks(password string, iFile, oFile *os.File, options EncryptionOptions) (err error) {
	var sought int
	var block cipher.Block
	var gcm cipher.AEAD
	var percentGetter utils.PercentGetter

	bufferSize := options.BufferSize - (options.NonceLength + 4)

	if percentGetter, err = utils.GetChunkPercentGetter(iFile, bufferSize); err != nil {
		return
	}

	buffer := make([]byte, bufferSize)
	nonce := make([]byte, options.NonceLength)

	for idx := 0; ; idx++ {

		if sought, err = iFile.Read(buffer); err != nil && err != io.EOF {
			err = fmt.Errorf("failed to read input file > %w", err)
			return
		} else if err == io.EOF {
			err = nil
			break
		}

		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			err = fmt.Errorf("failed to create random nonce > %w", err)
			return
		}

		key := pbkdf2.Key([]byte(password), nonce, options.KeyLength, options.IterationCount, options.Hash)

		if block, err = aes.NewCipher(key); err != nil {
			err = fmt.Errorf("failed to create new cipher > %w", err)
			return
		}

		if gcm, err = cipher.NewGCM(block); err != nil {
			err = fmt.Errorf("failed to create new GCM > %w", err)
			return
		}

		encrypted := append(gcm.Seal(nil, nonce, buffer[:sought], nil), nonce...)

		if _, err = io.Copy(oFile, bytes.NewReader(encrypted)); err != nil {
			err = fmt.Errorf("failed to copy to temporary file > %w", err)
			return
		}

		options.ProgressPercent <- percentGetter(idx)
	}

	return
}
