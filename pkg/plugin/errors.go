package plugin

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/internal/errorutil"
)

// noHandlingError is an error that won't be handled.
// This is a copy of errors.noHandlingError, to prevent import cycles
type noHandlingError struct {
	s string
}

func (e *noHandlingError) Error() string                       { return e.s }
func (e *noHandlingError) Is(err error) bool                   { return e == err }
func (e *noHandlingError) Handle(*state.State, *Context) error { return nil }

// stackError is a copy of errors.InternalError to prevent an import cycle.
type stackError struct {
	cause error
	s     errorutil.Stack
}

// errWithStack wraps the passed error in an error type containing stack
// information.
// If err is nil, errWithStack returns nil.
func errWithStack(err error) error {
	if err == nil {
		return nil
	}

	return &stackError{
		cause: err,
		s:     errorutil.GenerateStackTrace(1),
	}
}

func (s *stackError) Error() string { return s.cause.Error() }
func (s *stackError) Unwrap() error { return s.cause }
