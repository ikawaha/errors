package chainer

import (
	"errors"
)

type errChain struct {
	err  error
	next error
}

func newErrChain(err error) error {
	if err == nil {
		return nil
	}
	if errors.As(err, &errChain{}) {
		return err
	}
	return &errChain{err: err}
}

// Error implements error interface.
func (e errChain) Error() string {
	return e.err.Error()
}

// Is reports whether any error in err's chain matches target.
func (e errChain) Is(err error) bool {
	return errors.Is(e.err, err)
}

// As reports whether any error in err's chain matches target type.
func (e errChain) As(target interface{}) bool {
	return errors.As(e.err, target)
}

// Unwrap implements Wrapper interface.
func (e errChain) Unwrap() error {
	return e.next
}

// Append appends errors and obtains chained errors.
func Append(lhs error, rhs ...error) error {
	if len(rhs) == 0 || rhs[0] == nil {
		return newErrChain(lhs)
	}
	if lhs == nil {
		return Append(rhs[0], rhs[1:]...)
	}
	if len(rhs) == 1 {
		return &errChain{err: lhs, next: newErrChain(rhs[0])}
	}
	return Append(lhs, Append(rhs[0], rhs[1:]...))
}

// Yield yields a list of errors from the chained error.
func Yield(err error) []error {
	var el *errChain
	if !errors.As(err, &el) {
		return nil
	}
	return yield(err)
}

func yield(err error) []error {
	var el *errChain
	if !errors.As(err, &el) {
		return []error{err}
	}
	if el.next == nil {
		return yield(el.err)
	}
	return append(yield(el.err), yield(el.next)...)
}
