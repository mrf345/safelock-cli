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

	// prepare files to encrypt
	file, _ := os.CreateTemp("", "test_input")
	dirPath, _ := os.MkdirTemp("", "test_output")
	filePath := file.Name()
	encryptedFilePath := filepath.Join(dirPath, "encrypted.sla")

	// Encrypt `filePath` with the assigned settings
	if err := lock.Encrypt(ctx, filePath, encryptedFilePath, password); err != nil {
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

	// prepare files to decrypt
	encryptedFilePath, _ := getEncryptedFilePath()
	outputDir := filepath.Dir(encryptedFilePath)

	// Decrypt `encryptedFilePath` with the assigned settings
	if err := lock.Decrypt(ctx, encryptedFilePath, outputDir, password); err != nil {
		fmt.Println("failed!")
	}

	// Output:
}

func getEncryptedFilePath() (encryptedFilePath string, err error) {
	lock := safelock.New()
	lock.Quiet = true
	password := "testing123456"
	ctx := context.Background()
	file, _ := os.CreateTemp("", "test_input")
	dirPath, _ := os.MkdirTemp("", "test_output")
	filePath := file.Name()
	encryptedFilePath = filepath.Join(dirPath, "encrypted.sla")
	err = lock.Encrypt(ctx, filePath, encryptedFilePath, password)
	return
}
