package linker

import (
	"errors"
)

type errLinker struct {
	err  error
	next error
}

func newErrLinker(err error) error {
	if err == nil {
		return nil
	}
	if errors.As(err, &errLinker{}) {
		return err
	}
	return &errLinker{err: err}
}

func (e errLinker) Error() string {
	return e.err.Error()
}

func (e errLinker) Is(err error) bool {
	return errors.Is(e.err, err)
}

func (e errLinker) Unwrap() error {
	return e.next
}

func Append(lhs error, rhs ...error) error {
	if len(rhs) == 0 || rhs[0] == nil {
		return newErrLinker(lhs)
	}
	if lhs == nil {
		return Append(rhs[0], rhs[1:]...)
	}
	if len(rhs) == 1 {
		return &errLinker{err: lhs, next: newErrLinker(rhs[0])}
	}
	return Append(lhs, Append(rhs[0], rhs[1:]...))
}

func Yield(err error) []error {
	var el *errLinker
	if !errors.As(err, &el) {
		return nil
	}
	return yield(err)
}

func yield(err error) []error {
	var el *errLinker
	if !errors.As(err, &el) {
		return []error{err}
	}
	if el.next == nil {
		return yield(el.err)
	}
	return append(yield(el.err), yield(el.next)...)
}
