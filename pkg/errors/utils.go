package errors

import (
	"github.com/mavolin/adam/internal/errorutil"
)

// stackTrace attempts to extract the stacktrace from the error.
// If there is none, it will generate a stack trace.
func stackTrace(err error, skip int) (stack []uintptr) {
	if s, ok := err.(stackTracer); ok {
		stack = s.StackTrace()
	} else {
		stack = errorutil.GenerateStackTrace(1 + skip)
	}

	return
}
