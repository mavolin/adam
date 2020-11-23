package errors

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// multiError is a combination of multiple errors.
// They can be retrieved using RetrieveMultiError function.
type multiError []error

var _ Interface = new(multiError)

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

		return append(err1Typed, withStack(err2))
	} else if err2Typed, ok := err2.(multiError); ok {
		return append(multiError{withStack(err1)}, err2Typed...)
	}

	return multiError{withStack(err1), withStack(err2)}
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
		serr2.(*SilentError).stack = serr2.(*SilentError).stack[:len(serr2.(*SilentError).stack)-1]

		return append(err1Typed, serr2)
	} else if err2Typed, ok := err2.(multiError); ok {
		serr1 := Silent(err1)
		serr1.(*SilentError).stack = serr1.(*SilentError).stack[:len(serr1.(*SilentError).stack)-1]

		return append(multiError{serr1}, err2Typed...)
	}

	serr1 := Silent(err1)
	if len(serr1.(*SilentError).stack) > 1 {
		serr1.(*SilentError).stack = serr1.(*SilentError).stack[:len(serr1.(*SilentError).stack)-1]
	}

	serr2 := Silent(err2)
	if len(serr2.(*SilentError).stack) > 1 {
		serr2.(*SilentError).stack = serr2.(*SilentError).stack[:len(serr2.(*SilentError).stack)-1]
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
			merr = append(merr, withStack(err))
		}
	}

	return merr
}

// CombineSilent is the same as Combine, but wraps all errors that are not of
// type multiError as a SilentError.
func CombineSilent(errs ...error) error {
	if len(errs) == 0 {
		return nil
	} else if len(errs) == 1 {
		return Silent(errs[0])
	}

	var n int

	for i, err := range errs {
		if merr, ok := err.(multiError); ok {
			n += len(merr)
		} else {
			silent := Silent(err)
			if len(silent.(*SilentError).stack) > 1 {
				silent.(*SilentError).stack = silent.(*SilentError).stack[:len(silent.(*SilentError).stack)-1]
			}

			errs[i] = silent
			n++
		}
	}

	merr := make(multiError, 0, len(errs))

	for _, err := range errs {
		if sub, ok := err.(multiError); ok {
			merr = append(merr, sub...)
		} else {
			merr = append(merr, err)
		}
	}

	return merr
}

// RetrieveMultiError converts the passed errors to a single error.
// If the error is not of type multiError, it will return []error{err}.
func RetrieveMultiError(err error) []error {
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

// Handle handles the multi error.
// By default it iterates over the errors and calls the context's error handler
// on every one.
// If one of the errors is an InternalError, i.e. if at least one error was not
// added silently, it will be handled as such.
// All other errors will be handled as a SilentError.
func (merr multiError) Handle(s *state.State, ctx *plugin.Context) {
	HandleMultiError(merr, s, ctx)
}

var HandleMultiError = func(errs []error, s *state.State, ctx *plugin.Context) {
	internal := false

	for _, err := range errs {
		if !internal {
			if _, ok := err.(*InternalError); ok {
				ctx.HandleError(err)

				internal = true
				continue
			}
		}

		ctx.HandleErrorSilent(err)
	}
}
