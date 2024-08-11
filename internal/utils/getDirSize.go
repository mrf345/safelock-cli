package utils

import (
	"io/fs"
	"path/filepath"
)

func GetDirSize(path string) (size int64, err error) {
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return
}
