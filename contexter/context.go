package contexter

import (
	"errors"
	"reflect"
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
func WithValue(err error, key, val any) error {
	if err == nil {
		panic("cannot create an error form nil error")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &contextualError{
		err: err,
		key: key,
		val: val,
	}
}

// Value returns the value associated with the error for key.
// When no value is associated with key, Value returns false.
func Value(err error, key any) (any, bool) {
	return value(err, key)
}

func value(err error, key any) (any, bool) {
	var ctxErr *contextualError
	if !errors.As(err, &ctxErr) {
		return nil, false
	}
	if ctxErr.key == key {
		return ctxErr.val, true
	}
	return value(ctxErr.err, key)
}
