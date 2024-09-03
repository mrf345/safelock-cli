package safelock

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"

	"github.com/GianlucaGuarini/go-observable"
	"github.com/klauspost/compress/zstd"
	"github.com/mholt/archiver/v4"
)

// encryption/decryption configuration settings
type EncryptionConfig struct {
	IterationCount int
	// encryption key length (default: 48)
	KeyLength int
	// nonce length used to generate the encryption key (default: 12)
	NonceLength int
	// minimum password length allowed (default: 8)
	MinPasswordLength int
	// hashing method used to generate the encryption key (default: sha384)
	Hash func() hash.Hash
	// ratio to create file header size based on
	HeaderRatio int
	// lazy loaded GCM seals buffer size (default: 500)
	GcmBufferSize int
}

func (ec *EncryptionConfig) getMinPasswordLength() int {
	return ec.MinPasswordLength / 2
}

func (ec *EncryptionConfig) getHeaderSize(fileSize int) int {
	if ec.HeaderRatio > fileSize/100 {
		return ec.HeaderRatio
	}

	return int(fileSize) / ec.HeaderRatio
}

// archiving and compression configuration settings
type ArchiverConfig struct {
	// files compression (default: zstd.SpeedFastest)
	Compression archiver.Compression
	// files archiving (default: tar)
	Archival archiver.Archival
}

func (ac *ArchiverConfig) archive(ctx context.Context, output io.Writer, files []archiver.File) error {
	return archiver.CompressedArchive{Compression: ac.Compression, Archival: ac.Archival}.
		Archive(ctx, output, files)
}

// the main object used to configure safelock
type Safelock struct {
	EncryptionConfig
	ArchiverConfig

	// disable all output and logs (default: false)
	Quiet bool
	// observable instance that allows us to stream the status to multiple listeners
	StatusObs *observable.Observable
}

// creates a new [safelock.Safelock] instance with the default recommended options
func New() *Safelock {
	return &Safelock{
		StatusObs: observable.New(),
		ArchiverConfig: ArchiverConfig{
			Archival: archiver.Tar{},
			Compression: archiver.Zstd{
				EncoderOptions: []zstd.EOption{
					zstd.WithEncoderLevel(zstd.SpeedFastest),
				},
			},
		},
		EncryptionConfig: EncryptionConfig{
			IterationCount:    32,
			KeyLength:         48,
			NonceLength:       12,
			MinPasswordLength: 8,
			HeaderRatio:       7400,
			GcmBufferSize:     500,
			Hash:              sha512.New384,
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

// creates a new [safelock.Safelock] instance with sha512 hashing, could be slightly faster
func NewSha512() *Safelock {
	options := New()
	options.Hash = sha512.New
	options.KeyLength = 64
	return options
}
