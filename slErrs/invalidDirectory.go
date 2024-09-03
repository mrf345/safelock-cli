package slErrs

import "fmt"

// encountered invalid output path
type ErrInvalidOutputPath struct {
	BaseError,
	Path string
	Err error
}

func (e *ErrInvalidOutputPath) Error() string {
	return fmt.Sprintf("invalid output path (%s) > %s", e.Path, e.Err.Error())
}

func (e *ErrInvalidOutputPath) Is(t error) bool {
	_, ok := t.(*ErrInvalidOutputPath)
	return ok
}
