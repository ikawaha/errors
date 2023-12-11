package chainer

import "bytes"

type errChain struct {
	errs []error
}

// Chain chains non-nil errors and obtains chained errors.
func Chain(errs ...error) error {
	n := 0
	for _, v := range errs {
		if v != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	es := make([]error, 0, n)
	for _, v := range errs {
		if v != nil {
			es = append(es, v)
		}
	}
	return &errChain{errs: es}
}

// Error returns the head of the error chain.
func (e *errChain) Error() string {
	return e.errs[0].Error()
}

// Cause returns the messages of the remaining errors except the head error, joined by a newline.
func (e *errChain) Cause() string {
	var b bytes.Buffer
	for i, v := range e.errs[1:] {
		if i != 0 {
			b.WriteString("\n")
		}
		b.WriteString(v.Error())
	}
	return b.String()
}

// Unwrap implements Wrapper interface.
func (e *errChain) Unwrap() []error {
	return e.errs
}
