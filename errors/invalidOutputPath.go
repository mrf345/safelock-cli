package myErrs

import "fmt"

// path already exists or invalid
type ErrInvalidOutputPath struct {
	BaseError,
	Path string
}

func (e *ErrInvalidOutputPath) Error() string {
	return fmt.Sprintf("output path already exists (%s)", e.Path)
}
