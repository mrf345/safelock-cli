package safelock_test

import (
	"context"
	"errors"
	"io"
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
	outputFile, _ := sl.TempStore.NewFile("", "output_file")

	defer outputFile.RemoveQuietly()

	err := sl.Decrypt(context.TODO(), inputPath, outputFile.Name(), password)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidFile)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestDecryptWithInvalidOutputPath(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	sl := GetQuietSafelock()
	inputFile, _ := sl.TempStore.NewFile("", "input_file")
	outputFile, _ := sl.TempStore.NewFile("", "output_file")

	defer inputFile.RemoveQuietly()
	defer outputFile.RemoveQuietly()

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
	inputFile, _ := sl.TempStore.NewFile("", "input_file")
	outputDir, _ := sl.TempStore.NewDir("", "output_dir")
	outputPath := filepath.Join(outputDir.Path, "output_file.sla")

	defer cancel()
	defer inputFile.RemoveQuietly()
	defer outputDir.RemoveQuietly()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	encErr := sl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	decErr := sl.Decrypt(ctx, inputFile.Name(), outputPath, password)
	_, isExpectedErr := decErr.(*myErrs.ErrContextExpired)

	assert.Nil(encErr)
	assert.NotNil(decErr)
	assert.True(isExpectedErr)
}

func TestDecryptFileWithWrongPassword(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	sl := GetQuietSafelock()
	inputFile, _ := sl.TempStore.NewFile("", "input_file")
	outputDir, _ := sl.TempStore.NewDir("", "output_dir")
	outputPath := filepath.Join(outputDir.Path, "output_file.sla")

	defer inputFile.RemoveQuietly()
	defer outputDir.RemoveQuietly()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	encErr := sl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	decErr := sl.Decrypt(context.TODO(), outputPath, outputDir.Path, "wrong")
	_, isExpectedErr := errors.Unwrap(decErr).(*myErrs.ErrFailedToAuthenticate)

	assert.Nil(encErr)
	assert.NotNil(decErr)
	assert.True(isExpectedErr)
}
