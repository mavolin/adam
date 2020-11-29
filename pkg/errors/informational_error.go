package errors

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Abort is similar to a break in a for-loop.
// It stops the execution of a command silently, while producing neither a
// logged error nor a message to the calling user.
//
// It is intended to be used, if the user signals to cancel a command early
// and is therefore just a signaling error, rather than an actual exception.
var Abort error = &InformationalError{s: "abort"}

// InformationalError is an error that won't be handled.
// It is used to communicate preliminary stop of execution without signaling
// an actual error.
// See Abort for an example.
type InformationalError struct {
	s string
}

var _ Interface = new(InformationalError)

// NewInformationalError creates a new InformationalError with the passed
// error message.
func NewInformationalError(s string) *InformationalError {
	return &InformationalError{s: s}
}

func (e *InformationalError) Error() string                        { return e.s }
func (e *InformationalError) Is(err error) bool                    { return e == err }
func (e *InformationalError) Handle(*state.State, *plugin.Context) {}
