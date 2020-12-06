package plugin

import (
	"github.com/mavolin/disstate/v2/pkg/state"
)

// informationalError is an error that won't be handled.
// This is a copy of errors.InformationalError, to prevent import cycles.
type informationalError struct {
	s string
}

func (e *informationalError) Error() string                 { return e.s }
func (e *informationalError) Handle(*state.State, *Context) {}
