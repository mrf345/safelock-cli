package myErrs

import "fmt"

// wrong password or corrupted encryption
type ErrFailedToAuthenticate struct {
	BaseError,
	Msg string
}

func (e *ErrFailedToAuthenticate) Error() string {
	return fmt.Sprintf("invalid password or corrupted encryption > %s", e.Msg)
}
