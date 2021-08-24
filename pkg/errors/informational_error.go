package errors

import (
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Abort stops the execution of a command silently, while producing neither a
// logged error nor a message to the calling user.
//
// It is intended to be used if the user signals to cancel a command early.
var Abort error = NewInformationalError("abort")

// InformationalError is an error that won't be handled.
// It is used to communicate information, similar to io.EOF.
// See Abort for an example.
type InformationalError struct {
	s string
}

var _ Error = new(InformationalError)

// NewInformationalError creates a new InformationalError with the passed
// error message.
func NewInformationalError(s string) *InformationalError {
	return &InformationalError{s: s}
}

func (e *InformationalError) Error() string                              { return e.s }
func (e *InformationalError) Handle(*state.State, *plugin.Context) error { return nil }
