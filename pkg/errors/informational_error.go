package errors

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	_ "unsafe" // for go:linkname

	"github.com/mavolin/adam/pkg/plugin"
)

// Abort is similar to a break in a for-loop.
// It stops the execution of a command silently, while producing neither a
// logged error nor a message to the calling user.
//
// It is intended to be used, if the user signals to cancel a command early
// and is therefore just a signaling error, rather than an actual exception.
var Abort error = NewInformationalError("abort")

func init() {
	// Hack to make plugin.ErrInsufficientSendPermissions of type
	// InformationalError without creating an import cycle.
	//
	// We safely do this, since the error, unreplaced, is of type
	// errors.errorString, a type that can't be type asserted.
	// This means no code can be written that works based on the old type.
	// As soon as a type check for the replaced error is performed, the init
	// function will have already run during imports and everything will work
	// fine as well.
	// And since direct/pointer comparisons are unaffected by this, there won't
	// be an issue using that either.
	// Hence, no faulty code can be written.
	//
	// I know this is not pretty, but still better than copying this type,
	// using that type instead, and having the end user handle the error
	// specially.
	plugin.ErrInsufficientSendPermissions = NewInformationalError("insufficient permissions to send message")
}

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

func (e *InformationalError) Error() string                              { return e.s }
func (e *InformationalError) Is(err error) bool                          { return e == err }
func (e *InformationalError) Handle(*state.State, *plugin.Context) error { return nil }
