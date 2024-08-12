package myErrs

// passed context has ended or expired
type ErrContextExpired struct {
	BaseError
}

func (e *ErrContextExpired) Error() string {
	return "context has expired"
}
