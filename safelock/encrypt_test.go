package safelock_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"testing"

	myErrs "github.com/mrf345/safelock-cli/errors"
	"github.com/stretchr/testify/assert"
)

func TestEncryptWithInvalidInputPath(t *testing.T) {
	assert := assert.New(t)
	inputPath := "wrong_input.txt"
	password := "testing123456"
	sl := GetQuietSafelock()
	outputFile, _ := sl.TempStore.NewFile("", "output_file")

	defer outputFile.RemoveQuietly()
	err := sl.Encrypt(context.TODO(), inputPath, outputFile.Name(), password)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidFile)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestEncryptWithInvalidOutputPath(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	sl := GetQuietSafelock()
	inputFile, _ := sl.TempStore.NewFile("", "input_file")
	outputFile, _ := sl.TempStore.NewFile("", "output_file")

	defer inputFile.RemoveQuietly()
	defer outputFile.RemoveQuietly()
	err := sl.Encrypt(context.TODO(), inputFile.Name(), outputFile.Name(), password)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidOutputPath)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestEncryptFile(t *testing.T) {
	assert := assert.New(t)
	password := "testing123456"
	encSl := GetQuietSafelock()
	decSl := GetQuietSafelock()
	inputFile, _ := encSl.TempStore.NewFile("", "input_file")
	outputDir, _ := encSl.TempStore.NewDir("", "output_dir")
	outputPath := path.Join(outputDir.Path, "output_file.sla")
	content := "Hello World!"
	decryptedPath := path.Join(outputDir.Path, path.Base(inputFile.Name()))

	defer inputFile.RemoveQuietly()
	defer outputDir.RemoveQuietly()
	_, _ = inputFile.WriteString(content)
	inputFile.Close()

	inErr := encSl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	outErr := decSl.Decrypt(context.TODO(), outputPath, outputDir.Path, password)
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
	inputFile, _ := encSl.TempStore.NewFile("", "input_file")
	outputDir, _ := encSl.TempStore.NewDir("", "output_dir")
	outputPath := path.Join(outputDir.Path, "output_file.sla")
	decryptedPath := path.Join(outputDir.Path, path.Base(inputFile.Name()))

	defer inputFile.RemoveQuietly()
	defer outputDir.RemoveQuietly()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	inErr := encSl.Encrypt(context.TODO(), inputFile.Name(), outputPath, password)
	outErr := decSl.Decrypt(context.TODO(), outputPath, outputDir.Path, password)
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
	inputFile, _ := sl.TempStore.NewFile("", "input_file")
	outputDir, _ := sl.TempStore.NewDir("", "output_dir")
	outputPath := path.Join(outputDir.Path, "output_file.sla")

	defer cancel()
	defer inputFile.RemoveQuietly()
	defer outputDir.RemoveQuietly()
	_, _ = inputFile.WriteString(content)
	_, _ = inputFile.Seek(0, io.SeekStart)

	err := sl.Encrypt(ctx, inputFile.Name(), outputPath, password)
	_, isExpectedErr := err.(*myErrs.ErrContextExpired)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}
