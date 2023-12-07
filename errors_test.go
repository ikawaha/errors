package errors_test

import (
	"fmt"
	"testing"

	"github.com/ikawaha/errors"
)

// error with the stack trace
func errFn1() error {
	return errors.Errorf("error caused by fn1")
}

// fn1 -> fn2
func errFn2() error {
	return fmt.Errorf("fn2 has error: %w", errFn1())
}

func TestNew(t *testing.T) {
	err := errors.New("new_error")
	stack := errors.StackTrace(err)
	if len(stack) == 0 {
		t.Errorf("want not empty, but empty stack trace")
	}
}

func Test_Is(t *testing.T) {
	err0 := errors.New("error0")
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{name: "nil equal", err: nil, target: nil, want: true},
		{name: "equal", err: err0, target: err0, want: true},
		{name: "not equal", err: errors.New("error"), target: errors.New("error"), want: false},
		{name: "equal (embedded by std wrapper)", err: fmt.Errorf("caused by %w", err0), target: err0, want: true},
		{name: "equal (embedded by wrapper)", err: errors.Wrap(err0, "caused by"), target: err0, want: true},
		{name: "equal (chained: head)", err: errors.Chain(err0, errors.New("tails")), target: err0, want: true},
		{name: "equal (chained: tail)", err: errors.Chain(errors.New("head"), errors.New("tail1"), err0), target: err0, want: true},
		{name: "equal (chained: middle)", err: errors.Chain(errors.New("head"), errors.New("tail1"), err0, errors.New("tail2")), target: err0, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("want: %t, got: %t", tt.want, got)
			}
		})
	}
}

func TestChain(t *testing.T) {
	type args struct {
		head  error
		tails []error
	}
	tests := []struct {
		name    string
		head    error
		tails   []error
		want    bool
		wantErr string
	}{
		{
			name:  "empty",
			head:  nil,
			tails: nil,
			want:  false,
		},
		{
			name:    "head",
			head:    errors.New("head"),
			tails:   []error{errors.New("t0"), errors.New("t1"), errors.New("t2")},
			want:    true,
			wantErr: "head",
		},
		{
			name:    "head is nil, tail has some error",
			head:    nil,
			tails:   []error{nil, errors.New("t0"), errors.New("t1"), errors.New("t2")},
			want:    true,
			wantErr: "t0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Chain(append([]error{tt.head}, tt.tails...)...)
			if (err != nil) != tt.want {
				t.Fatalf("Chain() error = %v, want not nil", err)
			}
			if tt.want {
				if got := err.Error(); tt.wantErr != got {
					t.Errorf("want: %q, got: %q", tt.wantErr, got)
				}
			}
		})
	}
}
