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
	if s, ok := err.(stackTracer); ok { //nolint:errorlint
		stack = s.StackTrace()
	} else {
		stack = errorutil.GenerateStackTrace(1 + skip)
	}

	return
}

// retrieveCause retrieves the cause of the passed error, and returns it
// alongside the stack trace.
func retrieveCause(err error) (cause error, stack []uintptr, ok bool) { //nolint:golint,stylecheck
	var eerr Error
	eerrOk := As(err, &eerr)
	if eerrOk {
		err = eerr
	}

	switch err := err.(type) { //nolint:errorlint
	case *SilentError:
		cause = err.cause
		stack = err.stack
	case *InternalError:
		cause = err.cause
		stack = err.stack
	default:
		if eerrOk { // another Error, ignore
			return eerr, nil, false
		}

		cause = err
		stack = stackTrace(err, 2)
	}

	return cause, stack, true
}
