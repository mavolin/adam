// Package errorutil provides utilities to handle errors.
package errorutil

// StackError is a wrapper error that provides a stacktrace for the error.
type StackError struct {
	cause error
	s     StackTrace
}

// WithStack wraps the passed error into a StackError and attaches the stack
// trace of the caller to it.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	return &StackError{
		cause: err,
		s:     GenerateStackTrace(1),
	}
}

func (e *StackError) StackTrace() StackTrace { return e.s }
func (e *StackError) Unwrap() error          { return e.cause }
func (e *StackError) Error() string          { return e.cause.Error() }
