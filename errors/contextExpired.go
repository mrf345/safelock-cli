package myErrs

type ErrContextExpired struct {
	BaseError
}

func (e *ErrContextExpired) Error() string {
	return "context has expired"
}
