// explicitly defined errors you should handle
package myErrs

type BaseError struct {
	M string
}

func (e *BaseError) Error() string {
	return e.M
}
