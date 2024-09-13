package safelock

import (
	"context"
	"io"
	"runtime"

	"github.com/GianlucaGuarini/go-observable"
	"github.com/klauspost/compress/zstd"
	"github.com/mholt/archiver/v4"
)

// encryption/decryption configuration settings
type EncryptionConfig struct {
	// encryption key length (default: 32)
	KeyLength uint32
	// encryption salt length (default: 12)
	SaltLength int
	// number of argon2 hashing iterations (default: 3)
	IterationCount uint32
	// memory allocated for generating argon2 key (default: 64 * 1024)
	MemSize uint32
	// number of threads used to generate argon2 key (default: runtime.NumCPU())
	Threads uint8
	// minimum password length allowed (default: 8)
	MinPasswordLength int
	// ratio to create file header size based on (default: 1024 * 4)
	HeaderRatio int
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
		ArchiverConfig: ArchiverConfig{
			Archival: archiver.Tar{},
			Compression: archiver.Zstd{
				EncoderOptions: []zstd.EOption{
					zstd.WithEncoderLevel(zstd.SpeedFastest),
				},
			},
		},
		EncryptionConfig: EncryptionConfig{
			IterationCount:    3,
			KeyLength:         32,
			SaltLength:        12,
			MinPasswordLength: 8,
			HeaderRatio:       1024 * 4,
			MemSize:           64 * 1024,
			Threads:           uint8(runtime.NumCPU()),
		},
		StatusObs: observable.New(),
	}
}
