package contexter

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"slices"
)

type Frame []uintptr

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

// WithStacktrace adds a frame trace to the error.
// The skip specifies the number of frames to skip when getting the stack trace.
// Specify 0 when using it directly.
func WithStacktrace(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &contextualError{
		err: err,
		key: traceKey{},
		val: defaultStackTracer(err, skip),
	}
}

func StackFrames(err error) []Frame {
	find := false
	var frames []Frame
	for {
		var ctxErr *contextualError
		if !errors.As(err, &ctxErr) {
			break
		}
		if ctxErr.key == any(traceKey{}) {
			v, ok := ctxErr.val.(Frame)
			if !ok {
				continue
			}
			find = true
			frames = append(frames, v)
		}
		err = ctxErr.err
	}
	if !find {
		return nil
	}
	return frames
}

// StackTrace returns the Frame trace added to the error.
func StackTrace(err error) string {
	frames := StackFrames(err)
	if frames == nil {
		return ""
	}
	var b bytes.Buffer
	slices.Reverse(frames)
	for i := range frames {
		if i != 0 {
			b.WriteString("\n")
		}
		b.WriteString(frames[i].String())
	}
	return b.String()
}
