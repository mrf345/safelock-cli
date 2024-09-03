package safelock_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	slErrs "github.com/mrf345/safelock-cli/slErrs"
	"github.com/stretchr/testify/assert"
)

func TestDecryptWithInvalidOutputPath(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	sl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputFile, _ := os.CreateTemp("", "output_file")

	defer os.Remove(inputFile.Name())
	defer os.Remove(outputFile.Name())

	err := sl.Decrypt(context.TODO(), inputFile, outputFile.Name(), password)

	assert.NotNil(err)
	assert.True(slErrs.Is[*slErrs.ErrInvalidOutputPath](err))
}

func TestDecryptFileWithTimeout(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	sl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDirPath, _ := os.MkdirTemp("", "output_dir")
	outputFile, _ := os.CreateTemp(outputDirPath, "output_file.sla")
	inputPaths := []string{inputFile.Name()}

	cancel()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	encErr := sl.Encrypt(context.TODO(), inputPaths, outputFile, password)
	decErr := sl.Decrypt(ctx, outputFile, outputDirPath, password)

	assert.Nil(encErr)
	assert.NotNil(decErr)
	assert.True(errors.Is(decErr, context.DeadlineExceeded))

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
	outputFile, _ := os.CreateTemp(outputDirPath, "output_file.sla")
	inputPaths := []string{inputFile.Name()}

	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	encErr := sl.Encrypt(context.TODO(), inputPaths, outputFile, password)
	decErr := sl.Decrypt(context.TODO(), outputFile, outputDirPath, "wrong")

	assert.Nil(encErr)
	assert.NotNil(decErr)
	assert.True(slErrs.Is[*slErrs.ErrFailedToAuthenticate](decErr))

	// XXX: don't defer (temp files won't be deleted)
	os.Remove(inputFile.Name())
	os.RemoveAll(outputDirPath)
}
