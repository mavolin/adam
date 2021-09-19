// Package errors is a drop-in replacement of the stdlib's errors package.
// It provides enhanced error types, that store caller stacks.
//
// Additionally, errors defines custom error types that are specially handled
// when returned by plugin.Command.Invoke.
package errors

import (
	"errors"

	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/plugin"
)

// Error is an abstraction of a handleable error.
// It extends the built-in error.
type Error interface {
	error
	// Handle handles the error.
	//
	// If an error occurs during handling, it should be returned.
	// However, Handlers must make sure, that they don't infinitely return
	// errors, i.e. the handler returns the same error it is supposed to handle
	// either directly or through other Errors.
	//
	// To prevent this from happening, errors that deal with internal errors
	// should never return errors, or it must be made sure that only a finite
	// error chain will arise.
	Handle(s *state.State, ctx *plugin.Context) error
}

// Handle handles the passed error.
// If the error does not make itself available as an Error via As, it will be
// wrapped using WithStack.
//
// Up to maxHandles errors returned by Error.Handle and subsequent calls will
// be handled.
// If maxHandles is negative, subsequent handles won't be limited.
func Handle(s *state.State, ctx *plugin.Context, err error, maxHandles int) {
	for maxHandles != 0 && err != nil {
		var Err Error
		if !errors.As(err, &Err) {
			Err = WithStack(err)
		}

		err = Err.Handle(s, ctx)
		if maxHandles > 0 {
			maxHandles--
		}
	}
}

type StackTrace = errorutil.StackTrace

// GenerateStackTrace generates a StackTrace.
// It skips the given amount of callers.
func GenerateStackTrace(skip int) StackTrace {
	return errorutil.GenerateStackTrace(skip + 1)
}

// StackTracer is the interface implemented by all types providing stack
// traces.
type StackTracer interface {
	StackTrace() StackTrace
}
