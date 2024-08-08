package myErrs

import "fmt"

type ErrInvalidPassword struct {
	BaseError,
	Len int
}

func (e *ErrInvalidPassword) Error() string {
	return fmt.Sprintf("invalid password length (%d)", e.Len)
}
