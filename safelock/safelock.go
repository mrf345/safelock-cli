// ⚡ Fast files encryption (AES-GSM) package.
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

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/GianlucaGuarini/go-observable"
	"github.com/mholt/archiver/v4"
	"github.com/mrf345/safelock-cli/internal/utils"
)

// the main struct used to configure advanced encryption/decryption options
type Safelock struct {
	// files compression format (default: zstd)
	Compression archiver.Compression
	// files archiving tool (default: tar)
	Archival archiver.Archival
	// number of iterations performed to generate the encryption key (default: 32)
	IterationCount int
	// encryption key length (default: 64)
	KeyLength int
	// salt length used to generate the encryption key (default: 12)
	SaltLength int
	// encrypted/decrypted files buffer size (default: 4096)
	BufferSize int
	// encryption/decryption channels buffer size increasing/decreasing it might improve performance (default: 5)
	ChannelSize int
	// minimum password length allowed (default: 8)
	MinPasswordLength int
	// hashing method used to generate the encryption key (default: sha512)
	Hash func() hash.Hash
	// observable instance that allows us to stream the status to multiple listeners
	StatusObs *observable.Observable
	// disable all output and logs (default: false)
	Quiet bool
	// configures how we handle temporary files
	TempStore utils.TempStore
}

// creates a new [safelock.Safelock] instance with the default recommended options
func New() *Safelock {
	return &Safelock{
		Compression:       archiver.Zstd{},
		Archival:          archiver.Tar{},
		IterationCount:    32,
		KeyLength:         64,
		BufferSize:        4096,
		SaltLength:        12,
		MinPasswordLength: 8,
		ChannelSize:       5,
		Hash:              sha512.New,
		StatusObs:         observable.New(),
		TempStore: utils.TempStore{
			Cleanup:    true,
			PrefixText: "SF_temp",
		},
	}
}

// creates a new [safelock.Safelock] instance with sha256 hashing, it's faster (x.3) but less secure
func NewSha256() *Safelock {
	options := New()
	options.Hash = sha256.New
	options.KeyLength = 32
	return options
}

func (sl *Safelock) encryptionBufferSize() int {
	return sl.BufferSize - (sl.SaltLength + 4)
}

func (sl *Safelock) decryptionBufferSize() int {
	return sl.BufferSize + sl.SaltLength
}
