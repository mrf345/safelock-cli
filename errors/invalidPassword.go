package myErrs

import "fmt"

type ErrInvalidPassword struct {
	BaseError,
	Len int
	Need int
}

func (e *ErrInvalidPassword) Error() string {
	return fmt.Sprintf("invalid password length (%d) expected (%d)", e.Len, e.Need)
}
