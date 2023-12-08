package errors

import (
	stderrs "errors"

	"github.com/ikawaha/errors/contexter"
)

// As is provided for compatibility with errors.As.
func As(err error, target any) bool {
	return stderrs.As(err, target)
}

// Is is provided for compatibility with errors.Is.
func Is(err, target error) bool {
	return stderrs.Is(err, target)
}

// Join is provided for compatibility with errors.Join.
func Join(errs ...error) error {
	return stderrs.Join(errs...)
}

// New is provided for compatibility with errors.New.
// It captures the stack trace.
func New(text string) error {
	ret := stderrs.New(text)
	if stacktraceCapture {
		ret = contexter.WithStacktrace(ret, 1)
	}
	return ret
}

// Unwrap is provided for compatibility with errors.Unwrap.
// In particular Unwrap does not unwrap errors returned by Join or Chain.
func Unwrap(err error) error {
	return stderrs.Unwrap(err)
}
