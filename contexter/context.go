package contexter

import (
	"errors"
)

type contextualError struct {
	err error
	key any
	val any
}

// Error implements error interface.
func (e contextualError) Error() string {
	return e.err.Error()
}

// Unwrap implements Wrapper interface.
func (e contextualError) Unwrap() error {
	return errors.Unwrap(e.err)
}

// WithValue returns a wrapped error which has a value associated with key.
func WithValue(err error, key, val string) error {
	return &contextualError{
		err: err,
		key: key,
		val: val,
	}
}

// Value returns the value associated with the error for key.
// When no value is associated with key, Value returns false.
func Value(err error, key string) (string, bool) {
	return value(err, key)
}

func value(err error, key any) (string, bool) {
	var ctxErr *contextualError
	if !errors.As(err, &ctxErr) {
		return "", false
	}
	if ctxErr.key == key {
		v, ok := ctxErr.val.(string)
		return v, ok

	}
	return value(ctxErr.err, key)
}
