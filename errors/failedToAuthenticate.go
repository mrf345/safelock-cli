package myErrs

import "fmt"

type ErrFailedToAuthenticate struct {
	BaseError,
	Msg string
}

func (e *ErrFailedToAuthenticate) Error() string {
	return fmt.Sprintf("invalid password or corrupted encryption > %s", e.Msg)
}
