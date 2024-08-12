package utils

import (
	"io/fs"
	"os"

	myErrs "github.com/mrf345/safelock-cli/errors"
)

// check if the path is a valid directory
func IsValidDir(path string) (valid bool, err error) {
	var info fs.FileInfo

	if info, err = os.Stat(path); err != nil {
		return
	}

	valid = info.IsDir()

	if !valid {
		err = &myErrs.ErrInvalidDirectory{Path: path}
	}

	return
}
