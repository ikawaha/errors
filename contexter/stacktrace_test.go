package contexter_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ikawaha/errors/contexter"
)

func fn1() error {
	return contexter.WithStacktrace(errors.New("error caused by fn1"), 0)
}

func fn2() error {
	return fmt.Errorf("fn2 has error: %w", fn1())
}

func TestWithStackTrace(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err := fn2()
		if err == nil {
			t.Fatal("expected error, but nil")
		}
		want := "fn2 has error: error caused by fn1"
		if got := err.Error(); want != err.Error() {
			t.Errorf("want: %v, got: %v", want, got)
		}
		t.Run("frame trace", func(t *testing.T) {
			got := contexter.StackFrames(err)
			want := 1
			if len(got) != want {
				t.Errorf("want: len=%d, got: len=%d, %+v", want, len(got), got)
			}
		})
	})

	t.Run("wrapped wrapped error", func(t *testing.T) {
		err := contexter.WithStacktrace(fn2(), 0)
		if err == nil {
			t.Fatal("expected error, but nil")
		}
		want := "fn2 has error: error caused by fn1"
		if got := err.Error(); want != err.Error() {
			t.Errorf("want: %v, got: %v", want, got)
		}
		t.Run("frame trace", func(t *testing.T) {
			got := contexter.StackFrames(err)
			want := 2
			if len(got) != want {
				t.Errorf("want: len=%d, got: len=%d, %q", want, len(got), got)
			}
		})
	})
}
