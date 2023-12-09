package chainer

import (
	stderrs "errors"
)

type errChain struct {
	err  error
	next *errChain
}

// Chain chains non-nil errors and obtains chained errors.
func Chain(errs ...error) error {
	top := -1
	for i, v := range errs {
		if v != nil {
			top = i
			break
		}
	}
	if top < 0 {
		return nil
	}
	head := &errChain{
		err: errs[top],
	}
	p := head
	for i := top + 1; i < len(errs); i++ {
		if errs[i] == nil {
			continue
		}
		p.next = &errChain{
			err: errs[i],
		}
		p = p.next
	}
	return head
}

// Error returns the head of the error chain.
func (e *errChain) Error() string {
	return e.err.Error()
}

// Unwrap implements Wrapper interface.
func (e *errChain) Unwrap() error {
	if e.next == nil {
		return nil // untyped error
	}
	return e.next
}

func (e *errChain) Is(target error) bool {
	return stderrs.Is(e.err, target)
}

func (e *errChain) As(target any) bool {
	return stderrs.As(e.err, target)
}
