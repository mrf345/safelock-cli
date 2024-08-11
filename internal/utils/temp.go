package utils

import (
	"errors"
	"fmt"
	"os"
)

type TempStore struct {
	Cleanup    bool
	PrefixText string
	paths      []string
}

type TempFile struct {
	*os.File
	Remove        func() error
	RemoveQuietly func()
}

type TempDir struct {
	Path          string
	Remove        func() error
	RemoveQuietly func()
}

func (ts *TempStore) Prefix(name string) string {
	return fmt.Sprintf("%s_%s", ts.PrefixText, name)
}

func (ts *TempStore) NewFile(dir, name string) (tempFile *TempFile, err error) {
	var file *os.File

	if file, err = os.CreateTemp(dir, ts.Prefix(name)); err != nil {
		return
	}

	ts.paths = append(ts.paths, file.Name())

	tempFile = &TempFile{
		File: file,
		RemoveQuietly: func() {
			ts.Remove(file.Name())
		},
		Remove: func() error {
			return ts.Remove(file.Name())
		},
	}

	return
}

func (ts *TempStore) NewDir(dir, name string) (tempDir *TempDir, err error) {
	var path string

	if path, err = os.MkdirTemp(dir, ts.Prefix(name)); err != nil {
		return
	}

	ts.paths = append(ts.paths, path)

	tempDir = &TempDir{
		Path: path,
		RemoveQuietly: func() {
			ts.Remove(path)
		},
		Remove: func() error {
			return ts.Remove(path)
		},
	}

	return
}

func (ts *TempStore) Remove(path string) (err error) {
	for idx := len(ts.paths) - 1; idx >= 0; idx-- {
		if path == ts.paths[idx] {

			if err = os.RemoveAll(path); err != nil && !errors.Is(err, os.ErrNotExist) {
				return
			}

			err = nil
			ts.paths = append(ts.paths[:idx], ts.paths[idx+1:]...)
			return
		}
	}
	return
}

func (ts *TempStore) RemoveAll() {
	for idx := len(ts.paths) - 1; idx >= 0; idx-- {
		if err := os.RemoveAll(ts.paths[idx]); err != nil && !errors.Is(err, os.ErrNotExist) {
			ts.paths = append(ts.paths[:idx], ts.paths[idx+1:]...)
		}
	}
}
