package plugin

import "github.com/mavolin/disstate/pkg/state"

// noHandlingError is an error that won't be handled.
// This is a copy of errors.noHandlingError, to prevent import cycles
type noHandlingError struct {
	s string
}

func (e *noHandlingError) Error() string                       { return e.s }
func (e *noHandlingError) Is(err error) bool                   { return e == err }
func (e *noHandlingError) Handle(*state.State, *Context) error { return nil }
