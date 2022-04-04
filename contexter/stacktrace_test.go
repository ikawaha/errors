package contexter_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ikawaha/errors/contexter"
)

func fn1() error {
	return contexter.WithStackTrace(errors.New("error caused by fn1"))
}

func fn2() error {
	return fmt.Errorf("fn2 has error: %w", fn1())
}

func TestWithStackTrace(t *testing.T) {
	contexter.SetLogger(func(err error) string {
		return err.Error()
	})
	defer contexter.SetLogger(contexter.DefaultLogger)

	err := fn2()
	if err == nil {
		t.Fatal("expected error, but nil")
	}
	t.Run("wrapped error", func(t *testing.T) {
		want := "fn2 has error: error caused by fn1"
		if got := err.Error(); want != err.Error() {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})
	t.Run("stack trace", func(t *testing.T) {
		got, ok := contexter.StackTrace(err)
		if !ok {
			t.Errorf("expected true, but false")
		}
		want := "error caused by fn1"
		if got != want {
			t.Errorf("want: %q, got: %q", got, want)
		}
	})
}
