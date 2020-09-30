package errors

import (
	"fmt"

	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/mavolin/logstract/pkg/logstract"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/locutil"
)

// InternalError represents a non-user triggered error, that is reported to
// the user.
// By default, an InternalError does not explicitly state any information about
// the cause or context of the error, but sends a generalised message.
// However, using WithDescription, WithDetailsl or WithDetailslt, the default
// description can be replaced.
type InternalError struct {
	// cause is the cause of the error.
	cause error
	// stack contains information about the callers.
	stack errorutil.Stack

	// description of the error
	desc locutil.Text
}

// WithStack enriches the passed error with a stack trace.
// If the error is nil or it is another Interface, WithStack will return the
// error as is.
func WithStack(err error) error {
	return withStack(err)
}

// withStack enriches the passed error with a stack trace.
// If the error is nil or it is another Interface, withStack will return the
// error as is.
// If however, there is no stack trace, withStack skips the passed amount of
// frames including withStack itself and saves the callers.
func withStack(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(Interface); ok {
		return err
	}

	return &InternalError{
		cause: err,
		stack: stackTrace(err, 2),
		desc:  locutil.NewLocalizedText(defaultInternalDesc),
	}
}

// messageError is a simple error used for wrapped errors.
type messageError struct {
	msg   string
	cause error
}

func (e *messageError) Error() string { return fmt.Sprintf("%s: %s", e.msg, e.cause.Error()) }
func (e *messageError) Unwrap() error { return e.cause }

// Wrap wraps the passed error with the passed message and enriches it with a
// stack trace.
// The returned error will print as '$message: $err.Error()'.
func Wrap(err error, message string) *InternalError {
	if err == nil {
		return nil
	}

	return &InternalError{
		cause: &messageError{
			msg:   message,
			cause: err,
		},
		stack: stackTrace(err, 1),
		desc:  locutil.NewLocalizedText(defaultInternalDesc),
	}
}

// Wrapf wraps the passed error using the formatted passed message, and
// enriches the new error with a stack trace.
// The returned error will print as
// '$fmt.Sprintf(format, args...): $err.Error()'.
func Wrapf(err error, format string, args ...interface{}) *InternalError {
	if err == nil {
		return nil
	}

	return &InternalError{
		cause: &messageError{
			msg:   fmt.Sprintf(format, args...),
			cause: err,
		},
		stack: stackTrace(err, 1),
		desc:  locutil.NewLocalizedText(defaultInternalDesc),
	}
}

// WithDescription creates an internal error from the passed cause with the
// passed description.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.WithDescription(localizer).
func WithDescription(cause error, description string) *InternalError {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.desc = locutil.NewStaticText(description)
		return ie
	}

	return &InternalError{
		cause: cause,
		stack: stackTrace(cause, 1),
		desc:  locutil.NewStaticText(description),
	}
}

// WithDescription creates an internal error from the passed cause using
// the formatted description.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.WithDescription(localizer).
func WithDescriptionf(cause error, format string, args ...interface{}) *InternalError {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.desc = locutil.NewStaticText(fmt.Sprintf(format, args...))
		return ie
	}

	return &InternalError{
		cause: cause,
		stack: stackTrace(cause, 1),
		desc:  locutil.NewStaticText(fmt.Sprintf(format, args...)),
	}
}

// WithDescriptionl creates an internal error from the passed cause using the
// localized description.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.WithDescription(localizer).
func WithDescriptionl(cause error, description localization.Config) *InternalError {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.desc = locutil.NewLocalizedText(description)
		return ie
	}

	return &InternalError{
		cause: cause,
		stack: stackTrace(cause, 1),
		desc:  locutil.NewLocalizedText(description),
	}
}

// WithDescriptionlt creates an internal error from the passed cause using the
// message generated from the passed term.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.WithDescription(localizer).
func WithDescriptionlt(cause error, description localization.Term) *InternalError {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.desc = locutil.NewLocalizedText(description.AsConfig())
		return ie
	}

	return &InternalError{
		cause: cause,
		stack: stackTrace(cause, 1),
		desc:  locutil.NewLocalizedText(description.AsConfig()),
	}
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *InternalError) Description(l *localization.Localizer) string {
	if !e.desc.IsEmpty() {
		desc, err := e.desc.Get(l)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(defaultInternalDesc)
	return desc
}

func (e *InternalError) Error() string         { return e.cause.Error() }
func (e *InternalError) Unwrap() error         { return e.cause }
func (e *InternalError) StackTrace() []uintptr { return e.stack }

// Handle logs the error and sends out an internal error embed.
func (e *InternalError) Handle(_ *state.State, ctx *plugin.Context) error {
	logstract.
		WithFields(logstract.Fields{
			"cmd_ident": ctx.Command.Identifier,
			"err":       e,
		}).
		Error("command returned with an error")

	embed := ErrorEmbed.Clone().
		WithSimpleTitlel(internalErrorTitle).
		WithDescription(e.Description(ctx.Localizer))

	_, err := ctx.ReplyEmbedBuilder(embed)

	return err
}
