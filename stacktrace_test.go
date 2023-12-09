package errors_test

import (
	stderrs "errors"
	"fmt"
	"testing"

	"github.com/ikawaha/errors"
)

func TestWithStackTrace(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		fn1 := func() error {
			return errors.WithStacktrace(stderrs.New("error caused by fn1"))
		}
		fn2 := func() error {
			return fmt.Errorf("fn2 has error: %w", fn1())
		}
		fn3 := func() error {
			return errors.WithStacktrace(fn2())
		}
		tests := []struct {
			name string
			exec func() error
			want int // # of stack frames
		}{
			{name: "f1() cases error with stack trace", exec: fn1, want: 1},
			{name: "f2() cases error with stack trace", exec: fn2, want: 1},
			{name: "f3() cases error with 2 stack trace", exec: fn3, want: 2},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.exec()
				if err == nil {
					t.Fatal("expected error, but nil")
				}
				t.Run("frame trace", func(t *testing.T) {
					got := errors.StackFrames(err)
					if len(got) != tt.want {
						t.Errorf("want: len=%d, got: len=%d, %+v", tt.want, len(got), got)
					}
				})
				t.Run("print stack trace (for debug)", func(t *testing.T) {
					t.Skip(tt.name)
					fmt.Println("BEGIN:\n" + errors.Stacktrace(err) + "\nEND")
				})
			})
		}
	})
	t.Run("joined error", func(t *testing.T) {
		err0 := stderrs.New("err0")
		err1 := errors.New("err1")
		err2 := errors.Wrap(err1, "err2")

		tests := []struct {
			name string
			err  error
			want int // # of stack frames
		}{
			{name: "join err0, err1", err: stderrs.Join(err0, err1), want: 1},
			{name: "join err0, err1, err2", err: stderrs.Join(err0, err1, err2), want: 1},
			{name: "join err0, err2, err1", err: stderrs.Join(err0, err2, err1), want: 2},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Run("frame trace", func(t *testing.T) {
					got := errors.StackFrames(tt.err)
					if len(got) != tt.want {
						t.Errorf("want: len=%d, got: len=%d, %+v", tt.want, len(got), got)
					}
				})
				t.Run("print stack trace (for debug)", func(t *testing.T) {
					t.Skip(tt.name)
					fmt.Println("BEGIN:\n" + errors.Stacktrace(tt.err) + "\nEND")
				})
			})
		}
	})
}
