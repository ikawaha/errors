package errors

import (
	stderrs "errors"

	"github.com/ikawaha/errors/contexter"
)

func As(err error, target any) bool {
	return stderrs.As(err, target)
}

func Is(err, target error) bool {
	return stderrs.Is(err, target)
}

func Join(errs ...error) error {
	return stderrs.Join(errs...)
}

func New(text string) error {
	ret := stderrs.New(text)
	if captureStackTrace {
		ret = contexter.WithStackTrace(ret)
	}
	return ret
}

func Unwrap(err error) error {
	return stderrs.Unwrap(err)
}
