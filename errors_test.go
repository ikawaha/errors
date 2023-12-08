package errors_test

import (
	"testing"

	"github.com/ikawaha/errors"
)

func TestChain(t *testing.T) {
	tests := []struct {
		name    string
		head    error
		tails   []error
		want    bool
		wantErr string
	}{
		{name: "empty", head: nil, tails: nil, want: false},
		{name: "head", head: errors.New("head"), tails: []error{errors.New("t0"), errors.New("t1"), errors.New("t2")}, want: true, wantErr: "head"},
		{name: "head is nil, tail has some error", head: nil, tails: []error{nil, errors.New("t0"), errors.New("t1"), errors.New("t2")}, want: true, wantErr: "t0"},
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

func TestWithStacktrace(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := errors.WithStacktrace(nil)
		if err != nil {
			t.Errorf("want: nil, got: %v", err)
		}
	})
	t.Run("with stack trace", func(t *testing.T) {
		err := errors.WithStacktrace(errors.New("error"))
		if err == nil {
			t.Errorf("want: not nil, got: nil")
		}
		got := errors.Stacktrace(err)
		if len(got) == 0 {
			t.Errorf("want not empty, but empty stack trace")
		}
	})
}
