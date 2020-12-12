package errors

import (
	"fmt"

	"github.com/mavolin/disstate/v2/pkg/state"
	log "github.com/mavolin/logstract/pkg/logstract"

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

var _ Interface = new(SilentError)

// Silent creates a new silent error using the passed error as cause.
// If the error is already a SilentError, it will be returned as is.
// Furthermore, if the error is of type InternalError, the cause of the error
// will be extracted first, before creating a SilentError.
func Silent(err error) error {
	if err == nil {
		return nil
	}

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		return typedErr
	case *InternalError:
		err = typedErr.Unwrap()
	}

	return &SilentError{
		cause: err,
		stack: stackTrace(err, 1),
	}
}

// WrapSilent wraps the passed error with passed message, enriches the
// error with a stack trace, and marks the error as log-only.
// If the passed error is an InternalError or a SilentError WrapSilent will
// unwrap it first.
// The returned error will print as '$message: $err.Error()'.
func WrapSilent(err error, message string) error {
	if err == nil {
		return nil
	}

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case *InternalError:
		err = typedErr.Unwrap()
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
func WrapSilentf(err error, format string, args ...interface{}) error {
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

// Handle handles the SilentError.
// By default it logs the error.
//
// It will never return an error.
func (e *SilentError) Handle(s *state.State, ctx *plugin.Context) {
	HandleSilentError(e, s, ctx)
}

var HandleSilentError = func(serr *SilentError, s *state.State, ctx *plugin.Context) {
	log.
		WithFields(log.Fields{
			"cmd_ident": ctx.InvokedCommand.Identifier,
			"err":       serr.Unwrap().Error(),
		}).
		Error("command returned with error")
}
