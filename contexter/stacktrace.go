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

func defaultStackTracer(err error) Frame {
	if err == nil {
		return nil
	}
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n]
}

type traceKey struct{}

// WithStackTrace adds a Frame trace to the error.
func WithStackTrace(err error) error {
	if err == nil {
		return nil
	}
	return &contextualError{
		err: err,
		key: traceKey{},
		val: defaultStackTracer(err),
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
