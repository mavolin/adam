package errors

import (
	"errors"

	"github.com/mavolin/adam/internal/errorutil"
)

// errorStringWithStack is a trivial implementation of error with a stacktrace.
type errorStringWithStack struct {
	s     string
	stack errorutil.Stack
}

func (e *errorStringWithStack) Error() string         { return e.s }
func (e *errorStringWithStack) StackTrace() []uintptr { return e.stack }

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is
// identical.
func New(text string) error { return errors.New(text) }

// New returns an error that formats as the given text and stores the stack
// trace of the caller chain.
// Each call to New returns a distinct error value even if the text is
// identical.
func NewWithStack(text string) error {
	return &errorStringWithStack{
		s:     text,
		stack: errorutil.GenerateStackTrace(),
	}
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error { return errors.Unwrap(err) }

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained
// by repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error.
// For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == os.ErrExist }
//
// then Is(MyError{}, os.ErrExist) returns true.
// See syscall.Errno.Is for an example in the standard library.
func Is(err error, target error) bool { return errors.Is(err, target) }

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
// Otherwise, it returns false.
//
// The chain consists of err itself followed by the sequence of errors obtained
// by repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the
// value pointed to by target, or if the error has a method As(interface{})
// bool such that As(target) returns true.
// In the latter case, the As method is responsible for setting target.
//
// An error type might provide an As method so it can be treated as if it were
// a different error type.
//
// As panics if target is not a non-nil pointer to either a type that
// implements error, or to any interface type.
func As(err error, target interface{}) bool { return errors.As(err, target) }