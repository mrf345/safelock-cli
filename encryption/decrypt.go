package encryption

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/mholt/archiver/v4"
	myErrs "github.com/mrf345/safelock-cli/errors"
	myTypes "github.com/mrf345/safelock-cli/types"
	"github.com/mrf345/safelock-cli/utils"
	"golang.org/x/crypto/pbkdf2"
)

// TODO: add docs
// TODO: add tests
func Decrypt(inputPath, outputPath, password string, options myTypes.EncryptionOptions) (err error) {
	var channel = make(chan error)
	var percent, status string

	go func() {
		var archiveFile *os.File
		var archivePath string

		options.Status <- "Validating input and output"
		options.ProgressPercent <- "1%"
		if err = validateDecryptionPaths(inputPath, outputPath); err != nil {
			channel <- fmt.Errorf("invalid decryption input/output paths > %w", err)
			return
		}

		options.Status <- "Decrypting input file"
		options.ProgressPercent <- "5%"
		if archivePath, archiveFile, err = decryptArchiveFileInChunks(inputPath, outputPath, password, options); err != nil {
			channel <- fmt.Errorf("failed to decrypt archive file > %w", err)
			return
		}
		archiveFile.Seek(0, io.SeekStart)

		options.Status <- "Extracting files"
		options.ProgressPercent <- "95%"
		if err = extractArchiveFile(outputPath, archiveFile, options); err != nil {
			err = fmt.Errorf("failed to extract archive file > %w", err)
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
			status = fmt.Sprintf("Decrypted %s", outputPath)
			percent = "100%"
			options.Status <- status
			options.ProgressPercent <- percent
			utils.ClearAndPrint(options, "%s (%s)\n", status, "100%")
			close(channel)
			return
		case status = <-options.Status:
			utils.ClearAndPrint(options, "%s (%s)\n", status, percent)
		case percent = <-options.ProgressPercent:
			utils.ClearAndPrint(options, "%s (%s)\n", status, percent)
		}
	}
}

func validateDecryptionPaths(iPath, oPath string) (err error) {
	if _, err = utils.IsValidFile(iPath); err != nil {
		return
	}

	if _, err = utils.IsValidDir(oPath); err != nil {
		return
	}

	return
}

func decryptArchiveFileInChunks(
	iPath, oPath, password string,
	options myTypes.EncryptionOptions,
) (path string, file *os.File, err error) {
	var iFile *os.File
	var sought int
	var percentGetter utils.PercentGetter

	if iFile, err = os.Open(iPath); err != nil {
		return
	}

	if path, file, err = utils.CreateTempFile(iPath); err != nil {
		err = fmt.Errorf("failed to create temporary file > %w", err)
		return
	}

	bufferSize := options.BufferSize + options.NonceLength
	buffer := make([]byte, bufferSize)

	if percentGetter, err = utils.GetChunkPercentGetter(iFile, bufferSize); err != nil {
		return
	}

	for idx := 0; ; idx++ {
		var block cipher.Block
		var gcm cipher.AEAD
		var decrypted []byte

		options.ProgressPercent <- percentGetter(idx)

		if sought, err = iFile.Read(buffer); err != nil && err != io.EOF {
			err = fmt.Errorf("failed to read encrypted file > %w", err)
			return
		} else if err == io.EOF {
			err = nil
			break
		}

		encrypted := buffer[:sought-(options.NonceLength)]
		nonce := buffer[sought-(options.NonceLength) : sought]
		key := pbkdf2.Key([]byte(password), nonce, options.KeyLength, options.IterationCount, options.Hash)

		if block, err = aes.NewCipher(key); err != nil {
			err = fmt.Errorf("failed to create new cipher > %w", err)
			return
		}

		if gcm, err = cipher.NewGCM(block); err != nil {
			err = fmt.Errorf("failed to create new GCM > %w", err)
			return
		}

		if decrypted, err = gcm.Open(nil, nonce, encrypted, nil); err != nil {
			err = fmt.Errorf("failed to decrypt chunk > %w", err)
			return
		}

		if _, err = io.Copy(file, bytes.NewReader(decrypted)); err != nil {
			err = fmt.Errorf("failed to copy to temporary file > %w", err)
			return
		}
	}

	return
}

func extractArchiveFile(outputPath string, archiveFile *os.File, options EncryptionOptions) (err error) {
	var reader io.ReadCloser
	var fileHandler = getArchiveFileHandler(outputPath, options)

	if reader, err = options.Compression.OpenReader(archiveFile); err != nil {
		return fmt.Errorf("cannot read archive file > %w", err)
	}

	if err = options.Archival.Extract(options.Context, reader, nil, fileHandler); err != nil {
		return fmt.Errorf("cannot extract archive file > %w", err)
	}

	return
}

func getArchiveFileHandler(outputPath string, options EncryptionOptions) archiver.FileHandler {
	return func(ctx context.Context, file archiver.File) (err error) {
		var outputFile *os.File
		var reader io.ReadCloser
		var fullPath = path.Join(outputPath, file.NameInArchive)

		if file.IsDir() {
			os.MkdirAll(fullPath, file.Mode().Perm())
			return
		} else {
			os.MkdirAll(path.Dir(fullPath), file.Mode().Perm())
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
