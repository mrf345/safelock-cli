package safelock_test

import (
	"github.com/mholt/archiver/v4"
	"github.com/mrf345/safelock-cli/safelock"
)

func GetQuietSafelock() *safelock.Safelock {
	sl := safelock.New()
	sl.Quiet = true
	return sl
}

func GetQuietGzipSafelock() *safelock.Safelock {
	sl := safelock.New()
	sl.Compression = archiver.Gz{}
	sl.Quiet = true
	return sl
}
