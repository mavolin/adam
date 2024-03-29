package errors

import (
	"fmt"
	"io"
	"log"

	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
)

// InternalError represents a non-user triggered error.
//
// InternalErrors can be divided into two different categories:
//
// Silent errors are internal errors that are only handled internally by
// calling Log, and do not send an error message to the user.
// An InternalError is considered Silent if InternalError.Description returns
// an empty description.
//
// Non-Silent errors are internal errors that are handled internally by
// calling Log, and that send an error message to the user.
// By default, the description, i.e. the message sent to the user, does not
// contain any context or cause of the error, such as stack traces.
// However, a custom description can be added using WithDescription,
// WithDescriptionl or WithDescriptionlt.
type InternalError struct {
	// cause is the cause of the error.
	cause error
	// stack contains information about the callers.
	stackTrace StackTrace

	// desc is the description of the error.
	desc *i18n.Config
}

var (
	_ Error       = new(InternalError)
	_ StackTracer = new(InternalError)
)

// WithStack returns a new *InternalError using the caller's stack trace and
// the default description.
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the Error generated by As will be
// returned.
// If that error is of type *InternalError and has no description, the
// description will be set to the default description.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func WithStack(err error) Error {
	return withStack(err)
}

// withStack does the same as WithStack, but omits the calling method from the
// caller chain.
func withStack(err error) Error {
	if err == nil {
		return nil
	}

	var Err Error
	if As(err, &Err) {
		if ierr, ok := Err.(*InternalError); ok && ierr.desc == nil {
			cp := *ierr
			cp.desc = defaultInternalDesc
			return &cp
		}

		return Err
	}

	return &InternalError{
		cause:      err,
		stackTrace: stackTrace(err, 2),
		desc:       defaultInternalDesc,
	}
}

// Silent creates a new *InternalError with no description using the passed
// error as cause.
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfills As for Error, nil will be returned, unless the error
// generated by As is of type *InternalError.
// If the internal error has a description, it will be removed.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func Silent(err error) *InternalError {
	if err == nil {
		return nil
	}

	var Err Error
	if As(err, &Err) {
		if ierr, ok := Err.(*InternalError); ok {
			if ierr.desc == nil {
				return ierr
			}

			cp := *ierr
			cp.desc = nil
			return &cp
		}

		return nil
	}

	return &InternalError{
		cause:      err,
		stackTrace: stackTrace(err, 1),
	}
}

// MustInternal creates a new *InternalError from the passed error using the
// default description.
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfils As for *InternalError, the internal error generated by
// As will be returned.
// If the internal error has a description, it will be removed.
//
// In any other case, unlike WithStack, the error is wrapped in an
// *InternalError.
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, MustInternal will use that stack trace
// when wrapping, instead of creating one from the caller chain.
// MustInternal is therefore considered more forceful than WithStack, and cases
// that don't explicitly require this, should use WithStack.
func MustInternal(err error) *InternalError {
	if err == nil {
		return nil
	}

	var ierr *InternalError
	if As(err, &ierr) {
		if ierr.desc != nil {
			return ierr
		}

		cp := *ierr
		cp.desc = defaultInternalDesc
		return &cp
	}

	return &InternalError{
		cause:      err,
		stackTrace: stackTrace(err, 1),
		desc:       defaultInternalDesc,
	}
}

// MustSilent creates a new silent error using the passed error as cause.
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfils As for *InternalError, the internal error generated by
// As will be returned.
// If the internal error has no description, the description will be set to the
// default description.
//
// In any other case, unlike Silent, the error is wrapped in an *InternalError.
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, MustSilent will use that stack trace
// when wrapping, instead of creating one from the caller chain.
// MustSilent is therefore considered more forceful than Silent, and cases
// that don't explicitly require this should use Silent.
func MustSilent(err error) *InternalError {
	if err == nil {
		return nil
	}

	var ierr *InternalError
	if As(err, &ierr) {
		if ierr.desc == nil {
			return ierr
		}

		cp := *ierr
		cp.desc = nil
		return &cp
	}

	return &InternalError{
		cause:      err,
		stackTrace: stackTrace(err, 1),
	}
}

// Wrap wraps the passed error with the passed message and enriches it with a
// stack trace.
// The returned error will be an *InternalError using the default description,
// unless one of the below exceptions says otherwise.
//
// The returned error will print as:
// 	fmt.Sptrinf("%s: %s", message, err.Error())
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the Error generated by As will be
// returned.
// If that error is of type *InternalError, it's cause will be wrapped using
// the passed message.
// Furthermore, if the internal error has no description, it's description will
// be set to the default description.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func Wrap(err error, message string) Error {
	return wrap(err, message, false)
}

// WrapSilent wraps the passed error with passed message, enriches the error with a stack trace. The returned error will
// be an *InternalError with no description, unless one of the below exceptions says applies.
//
// The returned error will print as:
// 	fmt.Sprintf("%s: %s, fmt.Sprintf(format, a...), err.Error())
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the nil will be returned, unless the
// error generated by As is of type *InternalError.
// Furthermore, if the description of the internal error will be set to nil.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func WrapSilent(err error, message string) Error {
	return wrap(err, message, true)
}

// Wrapf wraps the passed error using the formatted passed message, and
// enriches the new error with a stack trace.
// The returned error will be an *InternalError using the default description,
// unless one of the below exceptions says otherwise.
//
// The returned error will print as:
// 	fmt.Sprintf("%s: %s, fmt.Sprintf(format, a...), err.Error())
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the Error generated by As will be
// returned.
// If that error is of type *InternalError, it's cause will be wrapped using
// the passed message.
// Furthermore, if the internal error has no description, it's description will
// be set to the default description.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func Wrapf(err error, format string, a ...interface{}) Error {
	return wrap(err, fmt.Sprintf(format, a...), false)
}

// WrapSilentf wraps the passed error using the formatted passed message,
// enriches the new error with a stack trace.
// The returned error will be an *InternalError using the default description,
// unless one of the below exceptions says otherwise.
//
// The returned error will print as:
// 	fmt.Sprintf("%s: %s, fmt.Sprintf(format, a...), err.Error())
//
// Exceptions
//
// If the error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the nil will be returned, unless the
// error generated by As is of type *InternalError.
// Furthermore, if the description of the internal error will be set to nil.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func WrapSilentf(err error, format string, a ...interface{}) Error {
	return wrap(err, fmt.Sprintf(format, a...), true)
}

// messageError is a simple error used for wrapped errors.
type messageError struct {
	msg   string
	cause error
}

func (e *messageError) Error() string { return fmt.Sprintf("%s: %s", e.msg, e.cause.Error()) }
func (e *messageError) Unwrap() error { return e.cause }

// wrap is the same as Wrap, but it omits the calling function from the stack
// trace.
// Additionally, it adds the silent parameter that defines whether or not the
// returned error shall be silenced.
func wrap(err error, message string, silent bool) Error {
	if err == nil {
		return nil
	}

	var Err Error
	if As(err, &Err) {
		if ierr, ok := Err.(*InternalError); ok {
			cp := *ierr
			ierr = &cp

			if silent {
				ierr.desc = nil
			} else if ierr.desc == nil {
				ierr.desc = defaultInternalDesc
			}

			ierr.cause = &messageError{
				msg:   message,
				cause: ierr.cause,
			}

			return ierr
		}

		if silent {
			return nil
		}

		return Err
	}

	ierr := &InternalError{
		cause: &messageError{
			msg:   message,
			cause: err,
		},
		stackTrace: stackTrace(err, 2),
	}

	if !silent {
		ierr.desc = defaultInternalDesc
	}

	return ierr
}

// WithDescription creates a new non-silent *InternalError from the passed
// error using the passed description.
//
// Exceptions
//
// If the passed error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the Error generated by As will be
// returned.
// If that error is of type *InternalError, it's description will be set to the
// given description.
func WithDescription(err error, description string) Error {
	if description == "" {
		return withDescription(err, nil)
	}

	return withDescription(err, i18n.NewStaticConfig(description))
}

// WithDescriptionf creates a new *InternalError from the passed error using
// the formatted description.
//
// Exceptions
//
// If the passed error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the Error generated by As will be
// returned.
// If that error is of type *InternalError, it's description will be set to the
// given description.
func WithDescriptionf(err error, format string, a ...interface{}) Error {
	descString := fmt.Sprintf(format, a...)
	if len(descString) == 0 {
		return withDescription(err, nil)
	}

	return withDescription(err, i18n.NewStaticConfig(descString))
}

// WithDescriptionl creates a new *InternalError from the passed cause using
// the localized description.
//
// Exceptions
//
// If the passed error is nil, it will be returned as is.
//
// If the error fulfills As for Error, the Error generated by As will be
// returned.
// If that error is of type *InternalError, it's description will be set to the
// given description.
func WithDescriptionl(err error, description *i18n.Config) Error {
	return withDescription(err, description)
}

// withDescription is the same as WithDescriptionl, but omits the calling
// function from the stack trace.
func withDescription(err error, description *i18n.Config) Error {
	if err == nil {
		return nil
	}

	var Err Error
	if As(err, &Err) {
		if ierr, ok := Err.(*InternalError); ok {
			cp := *ierr
			cp.desc = description
			return &cp
		}

		return Err
	}

	return &InternalError{
		cause:      err,
		stackTrace: stackTrace(err, 2),
		desc:       description,
	}
}

// IsSilent reports whether the error is silent, meaning it should not send any
// messages to the user upon being handled.
func (e *InternalError) IsSilent() bool {
	return e.desc == nil
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *InternalError) Description(l *i18n.Localizer) string {
	if e.desc == nil {
		return ""
	}

	desc, err := l.Localize(e.desc)
	if err == nil {
		return desc
	}

	desc, err = l.Localize(defaultInternalDesc)
	if err != nil {
		return ""
	}

	return desc
}

func (e *InternalError) Error() string          { return e.cause.Error() }
func (e *InternalError) Unwrap() error          { return e.cause }
func (e *InternalError) StackTrace() StackTrace { return e.stackTrace }

//goland:noinspection GoUnhandledErrorResult
func (e *InternalError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s\n", e.cause.Error())
			e.stackTrace.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error()) //nolint:errcheck
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// Handle handles the InternalError.
// If the error fulfills As for *httputil.HTTPError, it calls
// HandleDiscordError.
// If HandleDiscordError returns an Error of a type other than *InternalError,
// it handles it and returns.
//
// Otherwise, the error is logged, and handled using HandleInternalError.
// This does also apply if the error did not fulfil As for *httputil.HTTPError.
//
// Handle will never return a non-nil error.
func (e *InternalError) Handle(s *state.State, ctx *plugin.Context) error {
	if e == nil {
		return nil
	}

	var err Error = e

	if derr := discorderr.As(e); derr != nil {
		err = HandleDiscordError(e, derr)
		if err == nil {
			return nil
		}
	}

	ierr, ok := err.(*InternalError)
	if !ok {
		// prevent an infinite error cycle
		_ = err.Handle(s, ctx)
		return nil
	}

	Log(ctx, ierr)

	if !e.IsSilent() {
		newErr := HandleInternalError(s, ctx, e)
		if newIErr := Silent(newErr); newIErr != nil {
			_ = newIErr.Handle(s, ctx) // this will never return an error
		}
	}

	return nil
}

// HandleDiscordError takes in an *InternalError that fulfills As for
// *httputil.HTTPError and modifies it if appropriate.
// It returns the modified error.
var HandleDiscordError = func(ierr *InternalError, derr *httputil.HTTPError) Error {
	switch {
	case discorderr.Is(derr, discorderr.InsufficientPermissions):
		return plugin.DefaultBotPermissionsError
	case discorderr.Is(derr, discorderr.TemporarilyDisabled):
		return NewUserErrorl(discordErrorFeatureTemporarilyDisabled)
	case derr.Status >= 500:
		return NewUserErrorl(discordErrorServerError)
	default:
		return ierr
	}
}

// Log logs an InternalError.
//
// As InternalErrors can arise during any middleware, not all fields of the
// Context might be set.
// Hence, nil-checks should be performed on every nillable field not set by the
// router directly.
var Log = func(ctx *plugin.Context, err *InternalError) {
	if ctx.InvokedCommand == nil {
		log.Printf("internal error: %+v\n%+v\n", err, err.StackTrace())
	} else {
		log.Printf("internal error in command %s: %+v\n%+v", ctx.InvokedCommand.ID(), err, err.StackTrace())
	}
}

// HandleInternalError is called to handle a non-silent InternalError.
//
// By default, it sends an error embed containing the description of the error.
//
// As InternalErrors can arise during the execution of a default middleware,
// not all context fields might be set.
// Hence, nil-checks should be performed on every nillable field, that is not
// set by the router directly.
//
// If HandleInternalError returns an error, it will be given to Silent and then
// handled.
var HandleInternalError = func(s *state.State, ctx *plugin.Context, ierr *InternalError) (err error) {
	e := NewErrorEmbed(ctx.Localizer)

	e.Title, err = ctx.Localize(internalErrorTitle)
	if err != nil {
		return err
	}

	e.Description = ierr.Description(ctx.Localizer)

	_, err = ctx.ReplyEmbeds(e)
	return err
}

// stackTrace attempts to extract the stacktrace from the error.
// If there is none, it will generate a stack trace.
func stackTrace(err error, skip int) StackTrace {
	var tracer StackTracer
	if As(err, &tracer) {
		return tracer.StackTrace()
	}

	return GenerateStackTrace(skip + 1)
}
