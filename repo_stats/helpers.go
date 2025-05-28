package main

import "fmt"

func WrapError(err error, function string, context string) error {
	return fmt.Errorf("In %s: %s\n%s", function, err, context)
}
