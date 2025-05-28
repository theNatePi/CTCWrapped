package main

import "fmt"

func wrapError(err error, function string, context string) error {
	return fmt.Errorf("In %s: %s\n%s", function, err, context)
}
