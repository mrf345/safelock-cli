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

	// Prepare files to encrypt (should ignore this)
	inputFilePath, encryptedFilePath, clean := getUnencryptedFilePaths()
	defer clean()

	// Encrypt `inputFilePath` with the assigned settings
	if err := lock.Encrypt(ctx, inputFilePath, encryptedFilePath, password); err != nil {
		fmt.Println("failed!")
	}

	// Output:
}

func getUnencryptedFilePaths() (
	inputFilePath,
	encryptedFilePath string,
	clean func(),
) {
	inputFile, _ := os.CreateTemp("", "test_input")
	inputFilePath = inputFile.Name()
	outputDir, _ := os.MkdirTemp("", "test_output")
	encryptedFilePath = filepath.Join(outputDir, "encrypted.sla")
	clean = func() {
		os.Remove(inputFile.Name())
		os.RemoveAll(outputDir)
	}
	return
}

func ExampleSafelock_Decrypt() {
	lock := safelock.New()
	password := "testing123456"
	ctx := context.Background()

	// Disable logs and output
	lock.Quiet = true

	// Prepare files to decrypt (should ignore this)
	encryptedFilePath, clean := getEncryptedFilePath()
	outputDirPath := filepath.Dir(encryptedFilePath)
	defer clean()

	// Decrypt `encryptedFilePath` with the assigned settings
	if err := lock.Decrypt(ctx, encryptedFilePath, outputDirPath, password); err != nil {
		fmt.Println("failed!")
	}

	// Output:
}

func getEncryptedFilePath() (encryptedFilePath string, clean func()) {
	lock := safelock.New()
	lock.Quiet = true
	password := "testing123456"
	ctx := context.Background()
	file, _ := os.CreateTemp("", "test_input")
	dirPath, _ := os.MkdirTemp("", "test_output")
	filePath := file.Name()
	encryptedFilePath = filepath.Join(dirPath, "encrypted.sla")
	_ = lock.Encrypt(ctx, filePath, encryptedFilePath, password)
	clean = func() {
		os.Remove(file.Name())
		os.RemoveAll(dirPath)
	}
	return
}
