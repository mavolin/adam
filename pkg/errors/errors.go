// Package errors is a drop-in replacement of the stdlib's errors package.
// It provides enhanced error types, that store caller stacks.
//
// Additionally, errors defines custom error types that are specially handled
// when returned by plugin.Command.Invoke.
package errors

import (
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
