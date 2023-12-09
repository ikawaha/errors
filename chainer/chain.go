package chainer

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

// Unwrap implements Wrapper interface.
func (e *errChain) Unwrap() []error {
	return e.errs
}
