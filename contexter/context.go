package contexter

import (
	"errors"
)

type contextualError struct {
	err error
	key interface{}
	val string
}

func (e contextualError) Error() string {
	return e.err.Error()
}

func (e contextualError) Unwrap() error {
	return e.err
}

var traceKey struct{}

func WithValue(err error, key, val string) error {
	return &contextualError{
		err: err,
		key: key,
		val: val,
	}
}

func Value(err error, key string) (string, bool) {
	return value(err, key)
}

func value(err error, key interface{}) (string, bool) {
	var ctxErr *contextualError
	if !errors.As(err, &ctxErr) {
		return "", false
	}
	if ctxErr.key == key {
		return ctxErr.val, true
	}
	return value(ctxErr.err, key)
}
