package safelock

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"os"

	"github.com/GianlucaGuarini/go-observable"
	"github.com/mholt/archiver/v4"
	"github.com/mrf345/safelock-cli/utils"
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
	// nonce length used to generate the encryption key (default: 12)
	NonceLength int
	// encrypted/decrypted files buffer size (default: 64 * 1024)
	BufferSize int
	// encryption/decryption channels buffer size increasing/decreasing it might improve performance (default: 30)
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
	Registry *utils.FilesRegistry
}

// creates a new [safelock.Safelock] instance with the default recommended options
func New() *Safelock {
	return &Safelock{
		Compression:       archiver.Zstd{},
		Archival:          archiver.Tar{},
		IterationCount:    32,
		KeyLength:         64,
		BufferSize:        64 * 1024,
		NonceLength:       12,
		MinPasswordLength: 8,
		ChannelSize:       30,
		Hash:              sha512.New,
		StatusObs:         observable.New(),
		Registry: &utils.FilesRegistry{
			PrefixText: "SF_temp",
			Paths:      make(map[string]struct{}),
			TempDir:    os.TempDir(),
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
	return sl.BufferSize - (sl.NonceLength + 4)
}

func (sl *Safelock) decryptionBufferSize() int {
	return sl.BufferSize + sl.NonceLength
}
