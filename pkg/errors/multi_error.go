package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// multiError is a combination of multiple errors.
// They can be retrieved using RetrieveErrors function.
type multiError []error

// Append creates an error consisting of the passed errors.
// If one or both of them are multiErrors, it merges them and returns the
// result.
func Append(err1, err2 error) error {
	if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	}

	if err1Typed, ok := err1.(multiError); ok {
		if err2Typed, ok := err2.(multiError); ok {
			return append(err1Typed, err2Typed...)
		}

		return append(err1Typed, withStack(err2, 1))
	} else if err, ok := err2.(multiError); ok {
		return append(multiError{withStack(err1, 1)}, err...)
	}

	return multiError{withStack(err1, 1), withStack(err2, 1)}
}

// AppendSilent creates an error consisting of the passed errors.
// If one or both of them are multiErrors, it merges them and returns the
// result.
// All errors that are not of type multiError will be wrapped as a SilentError.
func AppendSilent(err1, err2 error) error {
	if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	}

	if err1Typed, ok := err1.(multiError); ok {
		if err2Typed, ok := err2.(multiError); ok {
			return append(err1Typed, err2Typed...)
		}

		serr2 := Silent(err2)
		serr2.stack = serr2.stack[1:]

		return append(err1Typed, serr2)
	} else if err, ok := err2.(multiError); ok {
		serr1 := Silent(err1)
		serr1.stack = serr1.stack[1:]

		return append(multiError{serr1}, err...)
	}

	serr1 := Silent(err1)
	if len(serr1.stack) > 1 {
		serr1.stack = serr1.stack[:len(serr1.stack)-1]
	}

	serr2 := Silent(err2)
	if len(serr2.stack) > 1 {
		serr2.stack = serr2.stack[:len(serr2.stack)-1]
	}

	return multiError{serr1, serr2}
}

// Combine combines the passed errors into a single error.
// If one of the passed errors is of type multiError, it will merged.
func Combine(errs ...error) error {
	if len(errs) == 0 {
		return nil
	} else if len(errs) == 1 {
		return errs[0]
	}

	var n int

	for _, err := range errs {
		if sub, ok := err.(multiError); ok {
			n += len(sub)
		} else {
			n++
		}
	}

	merr := make(multiError, 0, len(errs))

	for _, err := range errs {
		if sub, ok := err.(multiError); ok {
			merr = append(merr, sub...)
		} else {
			merr = append(merr, withStack(err, 1))
		}

	}

	return merr
}

// CombineSilent is the same as Combine, but wraps all errors that are not of
// type multiError as a SilentError.
func CombineSilent(errs ...error) error {
	for i, err := range errs {
		if _, ok := err.(multiError); !ok {
			silent := Silent(err)
			if len(silent.stack) > 1 {
				silent.stack = silent.stack[:len(silent.stack)-1]
			}

			errs[i] = silent
		}
	}

	return Combine(errs...)
}

// RetrieveErrors converts the passed errors to a single error.
// If the error is not of type multiError, it will return []error{err}.
func RetrieveErrors(err error) []error {
	merr, ok := err.(multiError)
	if ok {
		return merr
	}

	return []error{err}
}

func (merr multiError) Error() (s string) {
	s = merr[0].Error()

	for _, err := range merr[1:] {
		s += "; " + err.Error()
	}

	return
}

// Handle handles the passed errors.
// If one of them is an InternalError, i.e. if at least one error was not added
// silently, it will call handle on that.
// All other errors will be handled as a SilentError.
func (merr multiError) Handle(s *state.State, ctx *plugin.Context) (herr error) {
	internal := false

	for _, err := range merr {
		if !internal {
			if ierr, ok := err.(*InternalError); ok {
				herr = Append(herr, ierr.Handle(s, ctx))

				internal = true

				continue
			}
		}

		serr := Silent(err)

		herr = Append(herr, serr.Handle(s, ctx))
	}

	return herr
}
