package encryption

import (
	"context"
	"crypto/sha512"

	"github.com/mholt/archiver/v4"
	myTypes "github.com/mrf345/safelock-cli/types"
)

// TODO: add docs
// https://security.stackexchange.com/questions/110084/parameters-for-pbkdf2-for-password-hashing
func GetDefaultEncryptionOptions() myTypes.EncryptionOptions {
	return myTypes.EncryptionOptions{
		Compression:     archiver.Zstd{},
		Archival:        archiver.Tar{},
		IterationCount:  32,
		KeyLength:       64,
		BufferSize:      4096,
		NonceLength:     12,
		Hash:            sha512.New,
		Context:         context.Background(),
		Status:          make(chan string, 2),
		ProgressPercent: make(chan string, 2),
	}
}
