package contexter

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"sync"
)

var traceKey struct{}

type stackTraceLogger struct {
	mux sync.Mutex
	log func(error) string
}

var defLogger = stackTraceLogger{
	log: DefaultLogger,
}

// SetLogger sets a stack trace logger instead of the DefaultLogger.
func SetLogger(loggerFn func(error) string) {
	if loggerFn == nil {
		return
	}
	defLogger.mux.Lock()
	defer defLogger.mux.Unlock()
	defLogger.log = loggerFn
}

// DefaultLogger is the default logger for stackTraceLogger that logs
// the stack trace at this function called.
func DefaultLogger(err error) string {
	if err == nil {
		return ""
	}
	b := bytes.NewBufferString(err.Error())
	b.WriteString(":\n")
	for skip := 2; skip < 7; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc).Name()
		fmt.Fprintf(b, "    %s:%d %s\n", file, line, fn)
	}
	return b.String()
}

// WithStackTrace adds a stack trace to the error.
func WithStackTrace(err error) error {
	return &contextualError{
		err: err,
		key: traceKey,
		val: defLogger.log(err),
	}
}

// StackTrace returns the stack trace added to the error.
func StackTrace(err error) (string, bool) {
	var b bytes.Buffer
	find := false
	for {
		var ctxErr *contextualError
		if !errors.As(err, &ctxErr) {
			break
		}
		if ctxErr.key == traceKey {
			find = true
			b.WriteString(ctxErr.val)
		}
		err = ctxErr.err
	}
	return b.String(), find
}
