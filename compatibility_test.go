package errors_test

import (
	stderrs "errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/ikawaha/errors"
)

type asTestError struct {
	msg string
}

func (e *asTestError) Error() string {
	return e.msg
}

func TestAs(t *testing.T) {
	err0 := &asTestError{msg: "aloha"}
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "new: equal", err: err0, want: true},
		{name: "new: not found", err: errors.New("error"), want: false},
		{name: "found (embedded by errorf)", err: errors.Errorf("caused by %w", err0), want: true},
		{name: "found (embedded by wrapper)", err: errors.Wrap(err0, "caused by"), want: true},
		{name: "found (chained: top)", err: errors.NewWithErrors("head", err0, errors.New("tails")), want: true},
		{name: "found (chained: tail)", err: errors.NewWithErrors("head", errors.New("tail1"), err0), want: true},
		{name: "found (chained: middle)", err: errors.NewWithErrors("head", errors.New("tail1"), err0, errors.New("tail2")), want: true},
		{name: "found (chained: [[err, err, err0], MyErr])", err: errors.NewWithErrors("head0", errors.NewWithErrors("head1", errors.New("tail0"), errors.New("tail1"), err0), &asTestError{msg: "tail2"}), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// errors.Is
			var got *asTestError
			if errors.As(tt.err, &got) != tt.want {
				t.Errorf("want: %t, got: %t", tt.want, !tt.want)
			}
			if !tt.want {
				return
			}
			// std errors.Is
			var want *asTestError
			if stderrs.As(tt.err, &want) != tt.want {
				t.Errorf("want: %t, got: %t", tt.want, !tt.want)
			}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want: %v, got: %v", want, got)
			}
			if err0.Error() != got.Error() {
				t.Errorf("want: %v, got: %v", err0, got)
			}
		})
	}
}

func TestIs(t *testing.T) {
	err0 := errors.New("err0")
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{name: "nil equal", err: nil, target: nil, want: true},
		{name: "new: equal", err: err0, target: err0, want: true},
		{name: "new: not equal", err: errors.New("error"), target: errors.New("error"), want: false},
		{name: "equal (embedded by errorf)", err: errors.Errorf("caused by %w", err0), target: err0, want: true},
		{name: "equal (embedded by wrapper)", err: errors.Wrap(err0, "caused by"), target: err0, want: true},
		{name: "equal (chained: head)", err: errors.NewWithErrors("top", err0, errors.New("tails")), target: err0, want: true},
		{name: "equal (chained: tail)", err: errors.NewWithErrors("top", errors.New("head"), errors.New("tail1"), err0), target: err0, want: true},
		{name: "equal (chained: middle)", err: errors.NewWithErrors("top", errors.New("head"), errors.New("tail1"), err0, errors.New("tail2")), target: err0, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// errors.Is
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("want: %t, got: %t", tt.want, got)
			}
			// std errors.Is
			if got := stderrs.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("want: %t, got: %t", tt.want, got)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")
	tests := []struct {
		name string
		errs []error
		want error
	}{
		{name: "nil", errs: nil, want: nil},
		{name: "empty", errs: []error{}, want: nil},
		{name: "single", errs: []error{err1}, want: stderrs.Join(err1)},
		{name: "multiple", errs: []error{err1, err2, err3}, want: stderrs.Join(err1, err2, err3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Join(tt.errs...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %q, got: %q", tt.want, got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	err := errors.New("new_error")
	stack := errors.Stacktrace(err)
	if len(stack) == 0 {
		t.Errorf("want not empty, but empty stack trace")
	}
}

func TestUnwrap(t *testing.T) {
	err0 := errors.New("err0")
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "nil equal", err: nil, want: nil},
		{name: "single by new", err: err0, want: nil},
		{name: "single by std new", err: stderrs.New("err0"), want: nil},
		{name: "wrapped by wrap", err: errors.Wrap(err0, "wrapped"), want: err0},
		{name: "wrapped by errorf", err: errors.Errorf("caused by %w", err0), want: err0},
		{name: "wrapped by fmt errorf", err: fmt.Errorf("caused by %w", err0), want: err0},
		{name: "chain", err: errors.NewWithErrors("head", err0), want: err0},
		{name: "join", err: errors.Join(err0), want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Unwrap(tt.err)
			want := stderrs.Unwrap(tt.err)
			if !errors.Is(want, got) {
				t.Errorf("std errror.Unwrap(%v): %v, errors.Unwrap(%v): %v", tt.err, want, tt.err, got)
			}
			if !errors.Is(got, tt.want) {
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
			fmt.Println(got, want)
		})
	}
}
