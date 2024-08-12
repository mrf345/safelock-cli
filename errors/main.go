// explicitly defined errors you should handle
package myErrs

// base error struct to extend
type BaseError struct {
	M string
}

func (e *BaseError) Error() string {
	return e.M
}
