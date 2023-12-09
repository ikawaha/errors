errors
===

The errors package is an error handling library compatible with the standard Go errors package.It enables recording of stack traces during error handling and conceals error details that need not be shown to the user, without losing any error information.

```go
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
```
OUTPUT:
```
RAW ERROR: error causes in f2: error causes in f1
FOR USERS: error message for users
---stacktrace---
main.f1
	/Users/ikawaha/go/src/github.com/ikawaha/errors/_sample/main.go:10
main.f2
	/Users/ikawaha/go/src/github.com/ikawaha/errors/_sample/main.go:14
main.main
	/Users/ikawaha/go/src/github.com/ikawaha/errors/_sample/main.go:21
runtime.main
	/Users/ikawaha/sdk/go1.21.5/src/runtime/proc.go:267
main.f2
	/Users/ikawaha/go/src/github.com/ikawaha/errors/_sample/main.go:15
main.main
	/Users/ikawaha/go/src/github.com/ikawaha/errors/_sample/main.go:21
runtime.main
	/Users/ikawaha/sdk/go1.21.5/src/runtime/proc.go:267
```
---
MIT