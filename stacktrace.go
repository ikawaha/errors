package errors

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/ikawaha/errors/contexter"
)

// Frame represents a stack frame.
type Frame []uintptr

// String returns a string representation of the frame.
func (s *Frame) String() string {
	var b bytes.Buffer
	frames := runtime.CallersFrames(*s)
	for fr, ok := frames.Next(); ok; fr, ok = frames.Next() {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("%s\n\t%s:%d", fr.Function, fr.File, fr.Line))
	}
	return b.String()
}

func defaultStackTracer(err error, skip int) Frame {
	if err == nil {
		return nil
	}
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(3+skip, pcs[:])
	return pcs[:n]
}

type traceKey struct{}

// WithStacktraceSkip adds a frame trace to the error.
// The skip specifies the number of frames to skip when getting the stack trace.
// Specify 0 when using it directly.
func WithStacktraceSkip(err error, skip int) error {
	if err == nil {
		return nil
	}
	fs := StackFrames(err)
	return contexter.WithValue(err, traceKey{}, append(fs, defaultStackTracer(err, skip)))
}

// WithStacktrace adds a frame trace to the error.
func WithStacktrace(err error) error {
	return WithStacktraceSkip(err, 1)
}

// StackFrames returns the frame trace added to the error.
func StackFrames(err error) []Frame {
	v, ok := contexter.Value(err, traceKey{})
	if !ok {
		return nil
	}
	if ret, ok := v.([]Frame); ok {
		return ret
	}
	return nil
}

// Stacktrace returns the frame trace added to the error.
func Stacktrace(err error) string {
	frames := StackFrames(err)
	if frames == nil {
		return ""
	}
	var b bytes.Buffer
	for i := range frames {
		if i != 0 {
			b.WriteString("\n")
		}
		b.WriteString(frames[i].String())
	}
	return b.String()
}
