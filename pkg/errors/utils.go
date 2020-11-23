package errors

import (
	"github.com/mavolin/adam/internal/errorutil"
)

// stackTracer is used internally to check if an error already has a
// stacktrace.
type stackTracer interface {
	// StackTrace returns the program counters that led to the error.
	StackTrace() []uintptr
}

// stackTrace attempts to extract the stacktrace from the error.
// If there is none, it will generate a stack trace.
func stackTrace(err error, skip int) (stack []uintptr) {
	if s, ok := err.(stackTracer); ok { //nolint: errorlint
		stack = s.StackTrace()
	} else {
		stack = errorutil.GenerateStackTrace(1 + skip)
	}

	return
}
