package utils

import (
	"errors"
	"io/fs"
	"os"

	myErrs "github.com/mrf345/safelock-cli/errors"
)

// check if path is a valid file path
func IsValidFile(path string) (valid bool, err error) {
	var info fs.FileInfo

	if info, err = os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = &myErrs.ErrInvalidFile{Path: path}
		}
		return
	}

	valid = !info.IsDir()

	if !valid {
		err = &myErrs.ErrInvalidFile{Path: path}
	}

	return
}
