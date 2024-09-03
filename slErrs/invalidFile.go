package slErrs

import "fmt"

// encountered invalid input path
type ErrInvalidInputPath struct {
	BaseError,
	Path string
	Err error
}

func (e *ErrInvalidInputPath) Error() string {
	return fmt.Sprintf("invalid input path (%s) > %s", e.Path, e.Err.Error())
}

func (e *ErrInvalidInputPath) Is(t error) bool {
	_, ok := t.(*ErrInvalidInputPath)
	return ok
}
