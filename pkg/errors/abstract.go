package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// Handler handles an error.
	Handler interface {
		error
		// Handle handles the error.
		// If the Handler itself encounters an error, it may return it.
		Handle(s *state.State, ctx *plugin.Context) error
	}

	// stackTracer is used internally to check if an error already has a
	// stacktrace.
	stackTracer interface {
		// StackTrace returns the program counters that led to the error.
		StackTrace() []uintptr
	}
)
