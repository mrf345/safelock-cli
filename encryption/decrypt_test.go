package encryption_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/mrf345/safelock-cli/encryption"
	myErrs "github.com/mrf345/safelock-cli/errors"
	"github.com/stretchr/testify/assert"
)

func TestDecryptWithInvalidInputPath(t *testing.T) {
	assert := assert.New(t)
	inputPath := "wrong_input.txt"
	outputFile, _ := ioutil.TempFile("", "output_file")
	password := "testing123456"
	options := encryption.GetDefaultEncryptionOptions()

	err := encryption.Decrypt(inputPath, outputFile.Name(), password, options)
	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidFile)

	assert.Nil(err)
	assert.True(isExpectedErr)
}

// func TestDecryptWithInvalidOutputPath(t *testing.T) {
// 	assert := assert.New(t)
// 	inputFile, _ := ioutil.TempFile("", "1_input_file")
// 	outputFile, _ := ioutil.TempFile("", "output_file")
// 	password := "testing123456"
// 	options := encryption.GetDefaultEncryptionOptions()

// 	err := encryption.Decrypt(inputFile.Name(), outputFile.Name(), password, options)
// 	_, isExpectedErr := errors.Unwrap(err).(*myErrs.ErrInvalidOutputPath)

// 	assert.NotNil(err)
// 	assert.True(isExpectedErr)
// }

// func TestDecryptFile(t *testing.T) {
// 	assert := assert.New(t)
// 	inputFile, _ := ioutil.TempFile("", "input_file")
// 	outputDir, _ := ioutil.TempDir("", "output_dir")
// 	outputPath := path.Join(outputDir, "output_file.sla")
// 	password := "testing123456"
// 	options := encryption.GetDefaultEncryptionOptions()
// 	content := "Hello World!"
// 	decryptedPath := path.Join(outputDir, path.Base(inputFile.Name()))

// 	inputFile.WriteString(content)
// 	inputFile.Seek(0, io.SeekStart)

// 	inErr := encryption.Decrypt(inputFile.Name(), outputPath, password, options)
// 	outErr := encryption.Decrypt(outputPath, outputDir, password, options)
// 	reader, _ := os.Open(decryptedPath)
// 	decrypted, _ := ioutil.ReadAll(reader)

// 	assert.Nil(inErr)
// 	assert.Nil(outErr)
// 	assert.Equal(content, fmt.Sprintf("%s", decrypted))
// }

// func TestDecryptFileWithSha256AndGzip(t *testing.T) {
// 	assert := assert.New(t)
// 	inputFile, _ := ioutil.TempFile("", "input_file")
// 	outputDir, _ := ioutil.TempDir("", "output_dir")
// 	outputPath := path.Join(outputDir, "output_file.sla")
// 	password := "testing123456"
// 	content := "Hello World!"
// 	decryptedPath := path.Join(outputDir, path.Base(inputFile.Name()))
// 	options := encryption.GetDefaultEncryptionOptions()
// 	options.Compression = archiver.Gz{}
// 	options.Hash = sha256.New
// 	options.KeyLength = 32

// 	inputFile.WriteString(content)
// 	inputFile.Seek(0, io.SeekStart)

// 	inErr := encryption.Decrypt(inputFile.Name(), outputPath, password, options)
// 	outErr := encryption.Decrypt(outputPath, outputDir, password, options)
// 	reader, _ := os.Open(decryptedPath)
// 	decrypted, _ := ioutil.ReadAll(reader)

// 	assert.Nil(inErr)
// 	assert.Nil(outErr)
// 	assert.Equal(content, fmt.Sprintf("%s", decrypted))
// }

// func TestDecryptFileWithTimeout(t *testing.T) {
// 	assert := assert.New(t)
// 	inputFile, _ := ioutil.TempFile("", "input_file")
// 	outputDir, _ := ioutil.TempDir("", "output_dir")
// 	outputPath := path.Join(outputDir, "output_file.sla")
// 	password := "testing123456"
// 	content := "Hello World!"
// 	ctx, cancel := context.WithTimeout(context.Background(), 0)
// 	options := encryption.GetDefaultEncryptionOptions()
// 	options.Context = ctx

// 	defer cancel()
// 	inputFile.WriteString(content)
// 	inputFile.Seek(0, io.SeekStart)

// 	err := encryption.Decrypt(inputFile.Name(), outputPath, password, options)
// 	_, isExpectedErr := err.(*myErrs.ErrContextExpired)

// 	assert.NotNil(err)
// 	assert.True(isExpectedErr)
// }
