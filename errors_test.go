package errors_test

import (
	"fmt"
	"testing"

	"github.com/ikawaha/errors"
)

func TestNewWithErrors(t *testing.T) {
	tests := []struct {
		name    string
		head    string
		tails   []error
		want    bool
		wantErr string
	}{
		{name: "head", head: "head", tails: []error{errors.New("t0"), errors.New("t1"), errors.New("t2")}, want: true, wantErr: "head"},
		{name: "tail has some error including nil", head: "head", tails: []error{nil, errors.New("t0"), errors.New("t1"), errors.New("t2")}, want: true, wantErr: "head"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.NewWithErrors(tt.head, tt.tails...)
			if (err != nil) != tt.want {
				t.Fatalf("NewWithErrors() error = %v, want not nil", err)
			}
			if tt.want {
				if got := err.Error(); tt.wantErr != got {
					t.Errorf("want: %q, got: %q", tt.wantErr, got)
				}
			}
		})
	}
}

func TestWithErrors(t *testing.T) {
	tests := []struct {
		name    string
		head    error
		tails   []error
		want    bool
		wantErr string
	}{
		{name: "head", head: fmt.Errorf("head"), tails: []error{errors.New("t0"), errors.New("t1"), errors.New("t2")}, want: true, wantErr: "head"},
		{name: "tail has some error including nil", head: fmt.Errorf("head"), tails: []error{nil, errors.New("t0"), errors.New("t1"), errors.New("t2")}, want: true, wantErr: "head"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.WithErrors(tt.head, tt.tails...)
			if (err != nil) != tt.want {
				t.Fatalf("NewWithErrors() error = %v, want not nil", err)
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

func TestErrorf(t *testing.T) {
	t.Run("verb: %w", func(t *testing.T) {
		err0 := errors.New("err0")
		err := errors.Errorf("wrapped: %w", err0)
		if !errors.Is(err, err0) {
			t.Errorf("%%w should be able to be detected by errors.Is, but it was not.")
		}
	})
	t.Run("verb: %v", func(t *testing.T) {
		err0 := errors.New("err0")
		err := errors.Errorf("does not wrapped: %v", err0)
		if errors.Is(err, err0) {
			t.Errorf("%%v should not be able to be detected by errors.Is, but it was.")
		}
	})
}
