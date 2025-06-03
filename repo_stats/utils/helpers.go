package utils

import "fmt"

// WrapError
// Wraps errors with more context
//
// Parameters:
//   - err: original error
//   - function: the function from which the error will be returned
//   - context: additional text context to send with error
func WrapError(err error, function string, context string) error {
	return fmt.Errorf("In %s: %s\n%s", function, err, context)
}
