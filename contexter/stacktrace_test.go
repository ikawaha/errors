package contexter_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ikawaha/errors/contexter"
)

func TestWithStackTrace(t *testing.T) {
	fn1 := func() error {
		return contexter.WithStacktrace(errors.New("error caused by fn1"), 0)
	}
	fn2 := func() error {
		return fmt.Errorf("fn2 has error: %w", fn1())
	}
	fn3 := func() error {
		return contexter.WithStacktrace(fn2(), 0)
	}
	tests := []struct {
		name string
		exec func() error
		want int
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
				got := contexter.StackFrames(err)
				if len(got) != tt.want {
					t.Errorf("want: len=%d, got: len=%d, %+v", tt.want, len(got), got)
				}
			})
			t.Run("print stack trace", func(t *testing.T) {
				t.Skip(tt.name)
				fmt.Println("BEGIN:\n" + contexter.StackTrace(err) + "\nEND")
			})
		})
	}
}
