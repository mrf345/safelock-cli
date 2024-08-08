package myErrs

import "fmt"

type ErrInvalidFile struct {
	BaseError,
	Path string
}

func (e *ErrInvalidFile) Error() string {
	return fmt.Sprintf("invalid file path (%s)", e.Path)
}
