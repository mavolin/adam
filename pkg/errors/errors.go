// Package errors is a drop-in replacement of the stdlib's errors package.
// It provides enhanced error types, that store caller stacks.
//
// Additionally, errors defines custom error types that are specially handled
// when returned by plugin.Command.Invoke.
package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Interface is an abstraction of a handleable error.
// It extends the built-in error.
type Interface interface {
	error
	// Handle handles the error.
	// If the Interface itself encounters an error, it may return it.
	Handle(s *state.State, ctx *plugin.Context) error
}
