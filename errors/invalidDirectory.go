package myErrs

import "fmt"

type ErrInvalidDirectory struct {
	BaseError,
	Path string
}

func (e *ErrInvalidDirectory) Error() string {
	return fmt.Sprintf("invalid directory path (%s)", e.Path)
}
