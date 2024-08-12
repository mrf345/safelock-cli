// Fast files encryption (AES-GSM) package ⚡
//
// # Example
//
//	package main
//
//	import "github.com/mrf345/safelock-cli/safelock"
//
//	func main() {
//	  lock := safelock.New()
//	  inputPath := "/home/testing/important"
//	  outputPath := "/home/testing/encrypted.sla"
//	  extractTo := "/home/testing"
//	  password := "testing123456"
//
//	  // Encrypts `inputPath` with the default settings
//	  if err := lock.Encrypt(nil, inputPath, outputPath, password); err != nil {
//	    panic(err)
//	  }
//
//	  // Decrypts `outputPath` with the default settings
//	  if err := lock.Decrypt(nil, outputPath, extractTo, password); err != nil {
//	    panic(err)
//	  }
//	}
package safelock
