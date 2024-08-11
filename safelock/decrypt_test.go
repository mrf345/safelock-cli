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

func TestDecryptWithInvalidInputPath(t *testing.T) {
	assert := assert.New(t)
	inputPath := "wrong_input.txt"
	password := "testing123456"
	sl := GetQuietSafelock()
	outputFile, _ := os.CreateTemp("", "output_file")

	defer os.Remove(outputFile.Name())

	err := sl.Decrypt(context.TODO(), inputPath, outputFile.Name(), password)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidFile)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestDecryptWithInvalidOutputPath(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	sl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputFile, _ := os.CreateTemp("", "output_file")

	defer os.Remove(inputFile.Name())
	defer os.Remove(outputFile.Name())

	err := sl.Decrypt(context.TODO(), inputFile.Name(), outputFile.Name(), password)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidDirectory)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestDecryptFileWithTimeout(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	sl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDirPath, _ := os.MkdirTemp("", "output_dir")
	outputPath := filepath.Join(outputDirPath, "output_file.sla")

	cancel()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	encErr := sl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	decErr := sl.Decrypt(ctx, inputFile.Name(), outputPath, password)
	_, isExpectedErr := decErr.(*myErrs.ErrContextExpired)

	assert.Nil(encErr)
	assert.NotNil(decErr)
	assert.True(isExpectedErr)

	// XXX: don't defer (temp files won't be deleted)
	os.Remove(inputFile.Name())
	os.RemoveAll(outputDirPath)
}

func TestDecryptFileWithWrongPassword(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	sl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDirPath, _ := os.MkdirTemp("", "output_dir")
	outputPath := filepath.Join(outputDirPath, "output_file.sla")

	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	encErr := sl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	decErr := sl.Decrypt(context.TODO(), outputPath, outputDirPath, "wrong")
	_, isExpectedErr := errors.Unwrap(decErr).(*myErrs.ErrFailedToAuthenticate)

	assert.Nil(encErr)
	assert.NotNil(decErr)
	assert.True(isExpectedErr)

	// XXX: don't defer (temp files won't be deleted)
	os.Remove(inputFile.Name())
	os.RemoveAll(outputDirPath)
}
