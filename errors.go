package errors

import (
	stderrs "errors"
	"fmt"

	"github.com/ikawaha/errors/chainer"
)

var stacktraceCapture = true

// StacktraceCapture enables/disables stacktrace capture.
func StacktraceCapture(b bool) {
	stacktraceCapture = b
}

// Errorf is provided for compatibility with fmt.Errorf.
// It captures the stack trace.
func Errorf(format string, a ...any) error {
	ret := fmt.Errorf(format, a...)
	if stacktraceCapture {
		ret = WithStacktraceSkip(ret, 1)
	}
	return ret
}

// Wrap is provided for compatibility with github.com/pkg/errors.Wrap.
func Wrap(err error, text string) error {
	ret := fmt.Errorf("%s: %w", text, err)
	if stacktraceCapture {
		ret = WithStacktraceSkip(ret, 1)
	}
	return ret
}

// NewWithErrors creates new error of text and joins remaining errors, but the Error method of the joined error is
// the same as the head of the list. All information of the joined error is retained, so the information can be
// obtained with errors.Is, errors.As.
func NewWithErrors(text string, errs ...error) error {
	arg := make([]error, 0, len(errs)+1)
	arg = append(arg, stderrs.New(text))
	arg = append(arg, errs...)
	return chainer.Chain(arg...)
}
