package errors

import (
	"fmt"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/mavolin/logstract/pkg/logstract"

	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

// defaultInternalDescConfig is the localization.Config used by default as
// description for an InternalError.
var defaultInternalDescConfig = localization.Config{
	Term: termInternalDescription,
	Fallback: localization.Fallback{
		Other: "Oh no! Something went wrong and I couldn't finish executing your command. I've informed my team and " +
			"they'll get on fixing the bug asap.",
	},
}

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

	// description of the error, either is set
	descString string
	descConfig localization.Config
}

// WithStack enriches the passed error with a stack trace.
// If the error is nil or it already has a stack trace, WithStack returns
// the original error.
func WithStack(err error) error {
	return withStack(err, 1)
}

// withStack enriches the passed error with a stack trace.
// If the error is nil or it already has a stack trace, WithStack returns
// the original error.
// If however, there is no stack trace, withStack skips the passed amount of
// frames including withStack itself and saves the callers.
func withStack(err error, skip int) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(Handler); ok {
		return err
	}

	return &InternalError{
		cause:      err,
		stack:      stackTrace(err, 1+skip),
		descConfig: defaultInternalDescConfig,
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
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &InternalError{
		cause: &messageError{
			msg:   message,
			cause: err,
		},
		stack:      stackTrace(err, 1),
		descConfig: defaultInternalDescConfig,
	}
}

// Wrapf wraps the passed error using the formatted passed message, and
// enriches the new error with a stack trace.
// The returned error will print as
// '$fmt.Sprintf(format, args...): $err.Error()'.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &InternalError{
		cause: &messageError{
			msg:   fmt.Sprintf(format, args...),
			cause: err,
		},
		stack:      stackTrace(err, 1),
		descConfig: defaultInternalDescConfig,
	}
}

// WithDescription creates an internal error from the passed cause with the
// passed description.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.Description(localizer).
func WithDescription(cause error, desc string) error {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.descConfig = localization.Config{}
		ie.descString = desc
		return ie
	}

	return &InternalError{
		cause:      cause,
		stack:      stackTrace(cause, 1),
		descString: desc,
	}
}

// WithDescription creates an internal error from the passed cause using
// the formatted description.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.Description(localizer).
func WithDescriptionf(cause error, format string, args ...interface{}) error {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.descConfig = localization.Config{}
		ie.descString = fmt.Sprintf(format, args)
		return ie
	}

	return &InternalError{
		cause:      cause,
		stack:      stackTrace(cause, 1),
		descString: fmt.Sprintf(format, args),
	}
}

// WithDescriptionl creates an internal error from the passed cause using the
// localized description.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.Description(localizer).
func WithDescriptionl(cause error, c localization.Config) error {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.descConfig = c
		ie.descString = ""
		return ie
	}

	return &InternalError{
		cause:      cause,
		stack:      stackTrace(cause, 1),
		descConfig: c,
	}
}

// WithDescriptionlt creates an internal error from the passed cause using the
// message generated from the passed term.
// The description will be sent instead of a generic error message.
//
// When using a custom error handler, the description can be retrieved by
// calling internalError.Description(localizer).
func WithDescriptionlt(cause error, term string) error {
	if cause == nil {
		return nil
	}

	if ie, ok := cause.(*InternalError); ok {
		ie.descConfig = localization.Config{
			Term: term,
		}
		ie.descString = ""
		return ie
	}

	return &InternalError{
		cause: cause,
		stack: stackTrace(cause, 1),
		descConfig: localization.Config{
			Term: term,
		},
	}
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *InternalError) Description(l *localization.Localizer) (desc string) {
	var err error

	desc = e.descString
	if desc == "" {
		desc, err = l.Localize(e.descConfig)
		if err != nil { // use default on failure
			// we can safely discard this error, as it is impossible for
			// Localize to fail, ad defaultInternalDescConfig provides a
			// fallback.
			desc, _ = l.Localize(defaultInternalDescConfig)
		}
	}

	return desc
}

func (e *InternalError) Error() string         { return e.cause.Error() }
func (e *InternalError) Unwrap() error         { return e.cause }
func (e *InternalError) StackTrace() []uintptr { return e.stack }

// Handle logs the error, and sends it to sentry, if configured.
func (e *InternalError) Handle(_ *state.State, ctx *plugin.Context) error {
	logstract.
		WithFields(logstract.Fields{
			"cmd_ident": ctx.CommandIdentifier,
			"err":       e,
		}).
		Error("command returned with error")

	eventID := ctx.Hub.CaptureException(e)

	// We can ignore the error, as we have a fallback.
	title, _ := ctx.Localizer.Localize(errorTitleConfig)

	embed := discord.Embed{
		Title:       title,
		Description: e.Description(ctx.Localizer),
		Color:       constant.ErrorColor,
	}

	if eventID != nil {
		embed.Footer = &discord.EmbedFooter{
			Text: string(*eventID),
		}
	}

	_, err := ctx.ReplyEmbed(embed)

	return err
}
