package safelock_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	slErrs "github.com/mrf345/safelock-cli/slErrs"
	"github.com/stretchr/testify/assert"
)

func TestEncryptWithInvalidInputPath(t *testing.T) {
	assert := assert.New(t)
	inputPath := []string{"wrong_input.txt"}
	password := "testing123456"
	sl := GetQuietSafelock()
	outputFile, _ := os.CreateTemp("", "output_file")

	defer os.Remove(outputFile.Name())
	err := sl.Encrypt(context.TODO(), inputPath, outputFile, password)

	assert.NotNil(err)
	assert.True(slErrs.Is[*slErrs.ErrInvalidInputPath](err))
}

func TestEncryptFile(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	encSl := GetQuietSafelock()
	decSl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDir, _ := os.MkdirTemp("", "output_dir")
	outputFile, _ := os.CreateTemp(outputDir, "output_file.sla")
	content := "Hello World!"
	decryptedPath := filepath.Join(outputDir, filepath.Base(inputFile.Name()))
	inputPaths := []string{inputFile.Name()}

	defer os.Remove(inputFile.Name())
	defer os.RemoveAll(outputDir)
	_, _ = inputFile.WriteString(content)
	inputFile.Close()

	inErr := encSl.Encrypt(context.TODO(), inputPaths, outputFile, password)
	outErr := decSl.Decrypt(context.TODO(), outputFile, outputDir, password)
	reader, _ := os.Open(decryptedPath)
	decrypted, _ := io.ReadAll(reader)
	reader.Close()

	assert.Nil(inErr)
	assert.Nil(outErr)
	assert.Equal(content, string(decrypted))
}

func TestEncryptFileWithGzip(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	encSl := GetQuietGzipSafelock()
	decSl := GetQuietGzipSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDir, _ := os.MkdirTemp("", "output_dir")
	outputFile, _ := os.CreateTemp(outputDir, "output_file.sla")
	decryptedPath := filepath.Join(outputDir, filepath.Base(inputFile.Name()))
	inputPaths := []string{inputFile.Name()}

	defer os.Remove(inputFile.Name())
	defer os.RemoveAll(outputDir)
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)
	inputFile.Close()

	inErr := encSl.Encrypt(context.TODO(), inputPaths, outputFile, password)
	outErr := decSl.Decrypt(context.TODO(), outputFile, outputDir, password)
	reader, _ := os.Open(decryptedPath)
	decrypted, _ := io.ReadAll(reader)

	assert.Nil(inErr)
	assert.Nil(outErr)
	assert.Equal(content, string(decrypted))
}

func TestEncryptFileWithTimeout(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	content := "Hello World!"
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	sl := GetQuietSafelock()
	inputFile, _ := os.CreateTemp("", "input_file")
	outputDir, _ := os.MkdirTemp("", "output_dir")
	outputFile, _ := os.CreateTemp(outputDir, "output_file.sla")
	inputPaths := []string{inputFile.Name()}

	cancel()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)
	err := sl.Encrypt(ctx, inputPaths, outputFile, password)

	assert.NotNil(err)
	assert.True(errors.Is(err, context.DeadlineExceeded))

	// XXX: don't defer (temp files won't be deleted)
	os.Remove(inputFile.Name())
	os.RemoveAll(outputDir)
}
