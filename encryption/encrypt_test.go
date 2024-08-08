package encryption_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/mholt/archiver/v4"
	"github.com/mrf345/safelock-cli/encryption"
	myErrs "github.com/mrf345/safelock-cli/errors"
	"github.com/stretchr/testify/assert"
)

func TestEncryptWithInvalidInputPath(t *testing.T) {
	assert := assert.New(t)
	inputPath := "wrong_input.txt"
	outputFile, _ := ioutil.TempFile("", "output_file")
	password := "testing123456"
	options := encryption.GetDefaultEncryptionOptions()

	err := encryption.Encrypt(inputPath, outputFile.Name(), password, options)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidFile)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestEncryptWithInvalidOutputPath(t *testing.T) {
	assert := assert.New(t)
	inputFile, _ := ioutil.TempFile("", "1_input_file")
	outputFile, _ := ioutil.TempFile("", "output_file")
	password := "testing123456"
	options := encryption.GetDefaultEncryptionOptions()

	err := encryption.Encrypt(inputFile.Name(), outputFile.Name(), password, options)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidOutputPath)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}

func TestEncryptFile(t *testing.T) {
	assert := assert.New(t)
	inputFile, _ := ioutil.TempFile("", "input_file")
	outputDir, _ := ioutil.TempDir("", "output_dir")
	outputPath := path.Join(outputDir, "output_file.sla")
	password := "testing123456"
	options := encryption.GetDefaultEncryptionOptions()
	content := "Hello World!"
	decryptedPath := path.Join(outputDir, path.Base(inputFile.Name()))

	inputFile.WriteString(content)
	inputFile.Seek(0, io.SeekStart)

	inErr := encryption.Encrypt(inputFile.Name(), outputPath, password, options)
	outErr := encryption.Decrypt(outputPath, outputDir, password, options)
	reader, _ := os.Open(decryptedPath)
	decrypted, _ := ioutil.ReadAll(reader)

	assert.Nil(inErr)
	assert.Nil(outErr)
	assert.Equal(content, fmt.Sprintf("%s", decrypted))
}

func TestEncryptFileWithSha256AndGzip(t *testing.T) {
	assert := assert.New(t)
	inputFile, _ := ioutil.TempFile("", "input_file")
	outputDir, _ := ioutil.TempDir("", "output_dir")
	outputPath := path.Join(outputDir, "output_file.sla")
	password := "testing123456"
	content := "Hello World!"
	decryptedPath := path.Join(outputDir, path.Base(inputFile.Name()))
	options := encryption.GetDefaultEncryptionOptions()
	options.Compression = archiver.Gz{}
	options.Hash = sha256.New
	options.KeyLength = 32

	inputFile.WriteString(content)
	inputFile.Seek(0, io.SeekStart)

	inErr := encryption.Encrypt(inputFile.Name(), outputPath, password, options)
	outErr := encryption.Decrypt(outputPath, outputDir, password, options)
	reader, _ := os.Open(decryptedPath)
	decrypted, _ := ioutil.ReadAll(reader)

	assert.Nil(inErr)
	assert.Nil(outErr)
	assert.Equal(content, fmt.Sprintf("%s", decrypted))
}

func TestEncryptFileWithTimeout(t *testing.T) {
	assert := assert.New(t)
	inputFile, _ := ioutil.TempFile("", "input_file")
	outputDir, _ := ioutil.TempDir("", "output_dir")
	outputPath := path.Join(outputDir, "output_file.sla")
	password := "testing123456"
	content := "Hello World!"
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	options := encryption.GetDefaultEncryptionOptions()
	options.Context = ctx

	defer cancel()
	inputFile.WriteString(content)
	inputFile.Seek(0, io.SeekStart)

	err := encryption.Encrypt(inputFile.Name(), outputPath, password, options)
	_, isExpectedErr := err.(*myErrs.ErrContextExpired)

	assert.NotNil(err)
	assert.True(isExpectedErr)
}
