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

func (s *StackError) Error() string { return s.cause.Error() }
func (s *StackError) Unwrap() error { return s.cause }
