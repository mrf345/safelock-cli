package myErrs

import "fmt"

// file path doesn't exist or invalid
type ErrInvalidFile struct {
	BaseError,
	Path string
}

func (e *ErrInvalidFile) Error() string {
	return fmt.Sprintf("invalid file path (%s)", e.Path)
}
