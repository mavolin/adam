// Package errors is a drop-in replacement of the stdlib's errors package.
// It provides enhanced error types, that store caller stacks.
//
// Additionally, errors defines custom error types that are specially handled
// when returned by plugin.Command.Invoke.
package errors

import (
	"errors"
	"log"

	"github.com/mavolin/disstate/v3/pkg/state"

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

// Log is the logger used to log InternalErrors and SilentErrors.
var Log = func(err error, ctx *plugin.Context) {
	log.Printf("internal error in command %s: %+v\n", ctx.InvokedCommand.ID(), err)
}

// Handle handles the passed error.
// If the error does not implement interface Error, it will be wrapped using
// WithStack.
//
// Up to maxHandles errors returned by Error.Handle and subsequent calls will
// be handled.
// If maxHandles is negative, subsequent handles won't be limited.
func Handle(err error, s *state.State, ctx *plugin.Context, maxHandles int) {
	for maxHandles != 0 && err != nil {
		var Err Error
		if !errors.As(err, &Err) {
			Err = WithStack(err).(Error) //nolint:errorlint
		}

		err = Err.Handle(s, ctx)
		if maxHandles > 0 {
			maxHandles--
		}
	}
}
