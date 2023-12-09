package chainer_test

import (
	"errors"
	"testing"

	"github.com/ikawaha/errors/chainer"
)

func TestChain(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if err := chainer.Chain(); err != nil {
			t.Errorf("chainer.Chain() = %v, want nil", err)
		}
	})
	t.Run("nils", func(t *testing.T) {
		if err := chainer.Chain(nil); err != nil {
			t.Errorf("chainer.Chain(nil) = %v, want nil", err)
		}
		if err := chainer.Chain(nil, nil); err != nil {
			t.Errorf("chainer.Chain(nil, nil) = %v, want nil", err)
		}
	})
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	tests := []struct {
		name   string
		input  []error
		error  error
		unwrap error
	}{
		{name: "[err1]", input: []error{err1}, error: err1, unwrap: nil},
		{name: "[err1, err2]", input: []error{err1, err2}, error: err1, unwrap: err2},
		{name: "[err1, nil, err2]", input: []error{err1, nil, err2}, error: err1, unwrap: err2},
		{name: "[nil, err1, nil, err2, nil]", input: []error{nil, err1, nil, err2, nil}, error: err1, unwrap: err2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chainer.Chain(tt.input...)
			if !errors.Is(got, tt.error) {
				t.Errorf("want: %q, got: %q", tt.error, got)
			}
			if unwrapped := errors.Unwrap(got); !errors.Is(unwrapped, tt.unwrap) {
				t.Errorf("Unwrap(): want: %+v, got: %+v", tt.unwrap, got)
			}
		})
	}
}

func Test_Is(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")
	tests := []struct {
		name   string
		errs   []error
		target error
		want   bool
	}{
		{name: "Is([err1, err2], err1)", errs: []error{err1, err2}, target: err1, want: true},
		{name: "Is([err1, err2], err2)", errs: []error{err1, err2}, target: err2, want: true},
		{name: "Is([err1, err2], err3)", errs: []error{err1, err2}, target: err3, want: false},
		{name: "Is([err1, [err2, err3]], err3)", errs: []error{err1, chainer.Chain(err2, err3)}, target: err3, want: true},
		{name: "Is([err1, err2], err3], err3)", errs: []error{chainer.Chain(err1, err2), err3}, target: err3, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(chainer.Chain(tt.errs...), tt.target); got != tt.want {
				t.Errorf("expected %t, but %t", tt.want, got)
			}
		})
	}
}
