package main

import (
	"fmt"

	"github.com/ikawaha/errors"
)

func f1() error {
	return errors.New("error causes in f1")
}

func f2() error {
	if err := f1(); err != nil {
		return errors.Wrap(err, "error causes in f2")
	}
	return nil
}

func main() {
	if err := f2(); err != nil {
		// Print error details
		fmt.Printf("RAW ERROR: %v\n", err)

		// Creating a new error will create a new error message, but the information about the previous error will not be lost.
		err = errors.NewWithErrors("error message for users", err)
		fmt.Printf("FOR USERS: %v\n", err)

		// Print stack trace
		fmt.Println("---stacktrace---")
		fmt.Println(errors.Stacktrace(err))
	}
}
