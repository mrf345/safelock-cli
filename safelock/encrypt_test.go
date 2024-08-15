package safelock_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	myErrs "github.com/mrf345/safelock-cli/errors"
	"github.com/stretchr/testify/assert"
)

func TestEncryptWithInvalidInputPath(t *testing.T) {
	assert := assert.New(t)
	inputPath := "wrong_input.txt"
	password := "testing123456"
	sl := GetQuietSafelock()
	outputFile, _ := os.CreateTemp("", "output_file")

	defer os.Remove(outputFile.Name())
	err := sl.Encrypt(context.TODO(), inputPath, outputFile.Name(), password)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidFile)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestEncryptFile(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	encSl := GetQuietSafelock()
	decSl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDir, _ := os.MkdirTemp("", "output_dir")
	outputPath := filepath.Join(outputDir, "output_file.sla")
	content := "Hello World!"
	decryptedPath := filepath.Join(outputDir, filepath.Base(inputFile.Name()))

	defer os.Remove(inputFile.Name())
	defer os.RemoveAll(outputDir)
	_, _ = inputFile.WriteString(content)
	inputFile.Close()

	inErr := encSl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	outErr := decSl.Decrypt(context.TODO(), outputPath, outputDir, password)
	reader, _ := os.Open(decryptedPath)
	decrypted, _ := io.ReadAll(reader)
	reader.Close()

	assert.Nil(inErr)
	assert.Nil(outErr)
	assert.Equal(content, string(decrypted[:]))
}

func TestEncryptFileWithSha256AndGzip(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	encSl := GetQuietSha256GzipSafelock()
	decSl := GetQuietSha256GzipSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDir, _ := os.MkdirTemp("", "output_dir")
	outputPath := filepath.Join(outputDir, "output_file.sla")
	decryptedPath := filepath.Join(outputDir, filepath.Base(inputFile.Name()))

	defer os.Remove(inputFile.Name())
	defer os.RemoveAll(outputDir)
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	inErr := encSl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	outErr := decSl.Decrypt(context.TODO(), outputPath, outputDir, password)
	reader, _ := os.Open(decryptedPath)
	decrypted, _ := io.ReadAll(reader)

	assert.Nil(inErr)
	assert.Nil(outErr)
	assert.Equal(content, string(decrypted[:]))
}

func TestEncryptFileWithTimeout(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	sl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDir, _ := os.MkdirTemp("", "output_dir")
	outputPath := filepath.Join(outputDir, "output_file.sla")

	cancel()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	err := sl.Encrypt(ctx, inputFile.Name(), outputPath, password)
	_, isExpectedErr := err.(*myErrs.ErrContextExpired)

	assert.NotNil(err)
	assert.True(isExpectedErr)

	// XXX: don't defer (temp files won't be deleted)
	os.Remove(inputFile.Name())
	os.RemoveAll(outputDir)
}
