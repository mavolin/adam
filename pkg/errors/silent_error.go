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

var _ Error = new(SilentError)

// Silent creates a new silent error using the passed error as cause.
//
// If the error is already a *SilentError, it will be returned as is, and if
// the error is of type *InternalError, the cause of the error will be
// extracted first, before creating a SilentError.
// Furthermore, if the passed error fulfills As for Error but neither of the
// above conditions, nil will be returned, to prevent logging of errors that
// normally wouldn't be logged either.
func Silent(err error) error {
	if err == nil {
		return nil
	}

	err, stack, ok := retrieveCause(err)
	if !ok {
		return nil
	}

	return &SilentError{
		cause: err,
		stack: stack,
	}
}

// Silent creates a new silent error using the passed error as cause.
//
// If the error is already a *SilentError, it will be returned as is, and if
// the error is of type *InternalError, the cause of the error
// will be extracted first, before creating a SilentError.
// In any other case, unlike Silent, the error will wrapped in a silent error.
func MustSilent(err error) error {
	if err == nil {
		return nil
	}

	err, stack, _ := retrieveCause(err)

	return &SilentError{
		cause: err,
		stack: stack,
	}
}

// WrapSilent wraps the passed error with passed message, enriches the
// error with a stack trace, and marks the error as log-only.
//
// If the error is already a *SilentError, it will be returned as is, and if
// the error is of type *InternalError, the cause of the error will be
// extracted first, before creating a SilentError.
// Furthermore, if the passed error fulfills As for Error but neither of the
// above conditions, nil will be returned, to prevent logging of errors that
// normally wouldn't be logged either.
//
// The returned error will print as '$message: $err.Error()'.
func WrapSilent(err error, message string) error {
	if err == nil {
		return nil
	}

	err, stack, ok := retrieveCause(err)
	if !ok {
		return nil
	}

	return &SilentError{
		cause: &messageError{
			msg:   message,
			cause: err,
		},
		stack: stack,
	}
}

// WrapSilent wraps the passed error using the formatted passed message,
// enriches the error with a stack trace, and marks the error as log-only.
//
// If the error is already a *SilentError, it will be returned as is, and if
// the error is of type *InternalError, the cause of the error will be
// extracted first, before creating a SilentError.
// Furthermore, if the passed error fulfills As for Error but neither of the
// above conditions, nil will be returned, to prevent logging of errors that
// normally wouldn't be logged either.
//
// The returned error will print as
// '$fmt.Sprintf(format, args...): $err.Error()'.
func WrapSilentf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	err, stack, ok := retrieveCause(err)
	if !ok {
		return nil
	}

	return &SilentError{
		cause: &messageError{
			msg:   fmt.Sprintf(format, args...),
			cause: err,
		},
		stack: stack,
	}
}

func (e *SilentError) Error() string         { return e.cause.Error() }
func (e *SilentError) Unwrap() error         { return e.cause }
func (e *SilentError) StackTrace() []uintptr { return e.stack }

// Handle handles the SilentError.
// By default it logs the error.
//
// It will never return an error.
func (e *SilentError) Handle(s *state.State, ctx *plugin.Context) error {
	// prevent infinite error cycle, by not allowing error returns
	HandleSilentError(e, s, ctx)

	return nil
}

var HandleSilentError = func(serr *SilentError, s *state.State, ctx *plugin.Context) {
	log.
		WithFields(log.Fields{
			"cmd_ident": ctx.InvokedCommand.Identifier,
			"err":       serr.Unwrap().Error(),
		}).
		Error("command returned with error")
}
