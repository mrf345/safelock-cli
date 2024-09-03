// explicitly defined errors you should handle
package slErrs

import "errors"

// base error struct to extend
type BaseError struct {
	M string
}

func (e *BaseError) Error() string {
	return e.M
}

// check if error or unwrapped error  matches target
func Is[T error](err error) bool {
	target, ok := err.(T)
	return ok && errors.Is(err, target) || errors.Is(errors.Unwrap(err), target)
}
