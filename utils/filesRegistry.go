package utils

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

var mutex sync.Mutex

// keeps track of files created and helps to cleanup on sudden termination
type FilesRegistry struct {
	Disable    bool
	PrefixText string
	Paths      map[string]struct{}
}

// removable extended [os.File] type for [utils.FilesRegistry]
type RegFile struct {
	*os.File
	Remove func()
}

// get a prefix for temporary files
func (fs *FilesRegistry) Prefix(name string) string {
	return fmt.Sprintf("%s_%s", fs.PrefixText, name)
}

// create a new temporary file
func (fs *FilesRegistry) NewFile(dir, name string) (tempFile *RegFile, err error) {
	var file *os.File

	if file, err = os.CreateTemp(dir, fs.Prefix(name)); err != nil {
		return
	}

	fs.Paths[file.Name()] = struct{}{}
	tempFile = &RegFile{
		File: file,
		Remove: func() {
			_ = fs.Remove(file.Name())
		},
	}

	return
}

// remove previously created/registered file
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

// remove all registered files
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

// register a file
func (fs *FilesRegistry) Register(file *os.File) func() {
	mutex.Lock()
	defer mutex.Unlock()
	fs.Paths[file.Name()] = struct{}{}

	return func() {
		delete(fs.Paths, file.Name())
	}
}
