package errors

import (
	"fmt"

	"github.com/ikawaha/errors/chainer"
	"github.com/ikawaha/errors/contexter"
)

var captureStackTrace = true

func CaptureStackTrace(b bool) {
	captureStackTrace = b
}

func Errorf(format string, a ...any) error {
	ret := fmt.Errorf(format, a...)
	if captureStackTrace {
		ret = contexter.WithStackTrace(ret)
	}
	return ret
}

func Wrap(err error, text string) error {
	ret := Errorf("%s: %w", text, err)
	if captureStackTrace {
		ret = contexter.WithStackTrace(ret)
	}
	return ret
}

func Chain(errs ...error) error {
	return chainer.Chain(errs...)
}

func WithStacktrace(err error) error {
	return contexter.WithStackTrace(err)
}

func StackTrace(err error) string {
	ret := contexter.StackTrace(err)
	return ret
}
