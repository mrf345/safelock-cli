package utils

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

var mutex sync.Mutex

type FilesRegistry struct {
	Disable    bool
	PrefixText string
	Paths      map[string]struct{}
}

type RegFile struct {
	*os.File
	Remove func()
}

func (fs *FilesRegistry) Prefix(name string) string {
	return fmt.Sprintf("%s_%s", fs.PrefixText, name)
}

func (fs *FilesRegistry) NewFile(dir, name string) (tempFile *RegFile, err error) {
	var file *os.File

	if file, err = os.CreateTemp(dir, fs.Prefix(name)); err != nil {
		return
	}

	fs.Paths[file.Name()] = struct{}{}
	tempFile = &RegFile{
		File: file,
		Remove: func() {
			fs.Remove(file.Name())
		},
	}

	return
}

func (fs *FilesRegistry) Remove(path string) (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if fs.Disable {
		return
	}

	if _, ok := fs.Paths[path]; ok {
		if err = os.RemoveAll(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return
		}
		delete(fs.Paths, path)
	}

	return
}

func (fs *FilesRegistry) RemoveAll() (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if fs.Disable {
		return
	}

	for path := range fs.Paths {
		if err = os.RemoveAll(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return
		}

		delete(fs.Paths, path)
	}

	return
}

func (fs *FilesRegistry) Register(file *os.File) func() {
	mutex.Lock()
	defer mutex.Unlock()
	fs.Paths[file.Name()] = struct{}{}

	return func() {
		delete(fs.Paths, file.Name())
	}
}
