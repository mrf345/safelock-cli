package safelock_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrf345/safelock-cli/safelock"
)

func ExampleSafelock_Encrypt() {
	lock := safelock.New()
	password := "testing123456"
	ctx := context.Background()

	// Disable logs and output
	lock.Quiet = true

	// Increase minimum password length requirement
	lock.MinPasswordLength = 12

	// Prepare files to encrypt and clean up after test
	inputFile, _ := os.CreateTemp("", "test_input")
	outputFile, _ := os.CreateTemp("", "test_output")
	inputPaths := []string{inputFile.Name()}
	defer os.Remove(outputFile.Name())
	defer os.Remove(inputFile.Name())

	// This will encrypt all `inputPaths` and write into `outputFile`
	if err := lock.Encrypt(ctx, inputPaths, outputFile, password); err != nil {
		fmt.Println("failed!")
	}

	// Output:
}

func ExampleSafelock_Decrypt() {
	lock := safelock.New()
	password := "testing123456"
	ctx := context.Background()

	// Disable logs and output
	lock.Quiet = true

	// Prepare files to decrypt and clean up after test
	encryptedFile := getEncryptedFile(password)
	outputPath := filepath.Dir(encryptedFile.Name())
	defer os.Remove(encryptedFile.Name())
	defer os.RemoveAll(outputPath)

	// This will decrypt `encryptedFile` and extract files into `outputFile`
	if err := lock.Decrypt(ctx, encryptedFile, outputPath, password); err != nil {
		fmt.Println("failed!")
	}

	// Output:
}

func getEncryptedFile(password string) (outputFile *os.File) {
	lock := safelock.New()
	lock.Quiet = true
	ctx := context.TODO()

	inputFile, _ := os.CreateTemp("", "test_input")
	filePaths := []string{inputFile.Name()}
	outputFile, _ = os.CreateTemp("", "test_output")

	_ = lock.Encrypt(ctx, filePaths, outputFile, password)

	return
}
