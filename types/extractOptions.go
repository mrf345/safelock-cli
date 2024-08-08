package myTypes

import (
	"context"
	"hash"

	"github.com/mholt/archiver/v4"
)

// TODO: add docs
type EncryptionOptions struct {
	Compression     archiver.Compression
	Archival        archiver.Archival
	IterationCount  int
	KeyLength       int
	NonceLength     int
	BufferSize      int
	Hash            func() hash.Hash
	Context         context.Context
	Status          chan string
	ProgressPercent chan string
	Quiet           bool
}
