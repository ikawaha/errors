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
	t.Run("wrapped error", func(t *testing.T) {
		err := fn2()
		if err == nil {
			t.Fatal("expected error, but nil")
		}
		want := "fn2 has error: error caused by fn1"
		if got := err.Error(); want != err.Error() {
			t.Errorf("want: %v, got: %v", want, got)
		}
		t.Run("Frame trace", func(t *testing.T) {
			got := contexter.StackTrace(err)
			want := "github.com/ikawaha/errors/contexter_test.fn1\n" +
				"\t/Users/ikawaha/go/src/github.com/ikawaha/errors/contexter/stacktrace_test.go:12\n" +
				"github.com/ikawaha/errors/contexter_test.fn2\n" +
				"\t/Users/ikawaha/go/src/github.com/ikawaha/errors/contexter/stacktrace_test.go:16\n" +
				"github.com/ikawaha/errors/contexter_test.TestWithStackTrace.func1\n" +
				"\t/Users/ikawaha/go/src/github.com/ikawaha/errors/contexter/stacktrace_test.go:21\n" +
				"testing.tRunner\n\t/opt/homebrew/opt/go/libexec/src/testing/testing.go:1595"
			if got != want {
				t.Errorf("want: %q, got: %q", want, got)
			}
		})
	})

	t.Run("wrapped wrapped error", func(t *testing.T) {
		err := contexter.WithStackTrace(fn2())
		if err == nil {
			t.Fatal("expected error, but nil")
		}
		want := "fn2 has error: error caused by fn1"
		if got := err.Error(); want != err.Error() {
			t.Errorf("want: %v, got: %v", want, got)
		}
		t.Run("Frame trace", func(t *testing.T) {
			got := contexter.StackTrace(err)
			// Frame trace is reverse order
			want := "error caused by fn1" +
				"\n" + "fn2 has error: error caused by fn1"
			if got != want {
				t.Errorf("want: %q, got: %q", want, got)
				fmt.Println(got)
			}
		})
	})

}
