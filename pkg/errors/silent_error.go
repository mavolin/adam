package errors

import (
	"fmt"

	"github.com/mavolin/disstate/pkg/state"
	"github.com/mavolin/logstract/pkg/logstract"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/plugin"
)

// SilentError is an error that meant to be logged only, hence silent.
// A SilentError contains a complete stacktrace of when the error was first
// created.
type SilentError struct {
	// cause is the cause of the SilentError.
	cause error
	// stack contains the caller stack of the error.
	stack errorutil.Stack
}

// Silent creates a new silent error using the passed error as cause.
func Silent(err error) *SilentError {
	if err == nil {
		return nil
	}

	return &SilentError{
		cause: err,
		stack: stackTrace(err, 1),
	}
}

// WrapSilent wraps the passed error with passed message, enriches the
// error with a stack trace, and marks the error as log-only.
// The returned error will print as '$message: $err.Error()'.
func WrapSilent(err error, message string) *SilentError {
	if err == nil {
		return nil
	}

	return &SilentError{
		cause: &messageError{
			msg:   message,
			cause: err,
		},
		stack: stackTrace(err, 1),
	}
}

// WrapSilent wraps the passed error using the formatted passed message,
// enriches the error with a stack trace, and marks the error as log-only.
// The returned error will print as
// '$fmt.Sprintf(format, args...): $err.Error()'.
func WrapSilentf(err error, format string, args ...interface{}) *SilentError {
	if err == nil {
		return nil
	}

	return &SilentError{
		cause: &messageError{
			msg:   fmt.Sprintf(format, args...),
			cause: err,
		},
		stack: stackTrace(err, 1),
	}
}

func (e *SilentError) Error() string         { return e.cause.Error() }
func (e *SilentError) Unwrap() error         { return e.cause }
func (e *SilentError) StackTrace() []uintptr { return e.stack }

// Handle logs the error and sends it to sentry, if configured.
func (e *SilentError) Handle(_ *state.State, ctx *plugin.Context) error {
	logstract.
		WithFields(logstract.Fields{
			"cmd_ident": ctx.CommandIdentifier,
			"err":       e,
		}).
		Error("command returned with error")

	ctx.Hub.CaptureException(e)

	return nil
}
