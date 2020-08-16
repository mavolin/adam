package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// noHandlingError is an error that won't be handled.
type noHandlingError struct {
	s string
}

func (e *noHandlingError) Error() string                              { return e.s }
func (e *noHandlingError) Is(err error) bool                          { return e == err }
func (e *noHandlingError) Handle(*state.State, *plugin.Context) error { return nil }

// Abort is similar to a break in a for-loop.
// It stops the execution of a command silently, while producing neither a
// logged error nor a message to the calling user.
//
// It is intended to be used, if the user signals to cancel a command early
// and is therefore just a signaling error, rather than an actual exception.
var Abort error = &noHandlingError{
	s: "abort",
}
