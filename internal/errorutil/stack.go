package errorutil

import "runtime"

// stackDepth is the maximum depth of a Stack.
const stackDepth = 32

// Stack is a collection of program counters that show the stack trace of an
// error.
type Stack []uintptr

// GenerateStackTrace generates a Stack.
func GenerateStackTrace(skip int) Stack {
	pcs := make([]uintptr, stackDepth)

	n := runtime.Callers(2+skip, pcs)

	return pcs[0:n]
}

type StackError struct {
	cause error
	s     Stack
}

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
