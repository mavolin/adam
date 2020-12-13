package errors

import (
	"fmt"

	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/mavolin/logstract/pkg/logstract"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

// InternalError represents a non-user triggered error.
// By default, an InternalError does not explicitly state any information about
// the cause or context of the error, and instead sends a generalised message.
// However, a custom description can be added using WithDescription,
// WithDescriptionl or WithDescriptionlt.
type InternalError struct {
	// cause is the cause of the error.
	cause error
	// stack contains information about the callers.
	stack errorutil.Stack

	// description of the error
	desc *i18nutil.Text
}

var _ Error = new(InternalError)

// WithStack returns a new InternalError using the callers stack trace.
//
// If the error is a *SilentError, WithStack will unwrap it first.
//
// In case the error is nil or another Error type, WithStack will return the
// error as is.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func WithStack(err error) error {
	return withStack(err)
}

// withStack enriches the passed error with a stack trace.
//
// If the error is nil or it is another Error expect *SilentError, withStack
// will return the error as is.
//
// If the passed error already provides a stack trace via a
// err.StackTrace() []uintptr method, WithStack will use that stack trace when
// wrapping, instead of creating one from the caller chain.
func withStack(err error) error {
	if err == nil {
		return nil
	}

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case Error:
		return err
	}

	return &InternalError{
		cause: err,
		stack: stackTrace(err, 2),
		desc:  i18nutil.NewTextl(defaultInternalDesc),
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
//
// If the passed error is an *InternalError or a *SilentError, Wrap will unwrap
// it first.
//
// The returned error will print as '$message: $err.Error()'.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	stack := stackTrace(err, 1)

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case *InternalError:
		err = typedErr.Unwrap()
	}

	return &InternalError{
		cause: &messageError{
			msg:   message,
			cause: err,
		},
		stack: stack,
		desc:  i18nutil.NewTextl(defaultInternalDesc),
	}
}

// Wrapf wraps the passed error using the formatted passed message, and
// enriches the new error with a stack trace.
//
// If the passed error is an *InternalError or a *SilentError, Wrapf will unwrap
// it first.
// The returned error will print as
// '$fmt.Sprintf(format, args...): $err.Error()'.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	stack := stackTrace(err, 1)

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case *InternalError:
		err = typedErr.Unwrap()
	}

	return &InternalError{
		cause: &messageError{
			msg:   fmt.Sprintf(format, args...),
			cause: err,
		},
		stack: stack,
		desc:  i18nutil.NewTextl(defaultInternalDesc),
	}
}

// WithDescription creates a new InternalError from the passed error using the
// passed description.
//
// If the passed description is empty, the error will use the default
// description.
//
// If the passed error is a *SilentError or an *InternalError, WithDescription
// will unwrap it first.
func WithDescription(err error, description string) error {
	if err == nil {
		return nil
	}

	stack := stackTrace(err, 1)

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case *InternalError:
		typedErr.desc = i18nutil.NewText(description)
		return typedErr
	}

	return &InternalError{
		cause: err,
		stack: stack,
		desc:  i18nutil.NewText(description),
	}
}

// WithDescription creates a new *InternalError from the passed error using the
// formatted description.
//
// If the passed error is a *SilentError or an *InternalError, WithDescriptionf
// will unwrap it first.
func WithDescriptionf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	stack := stackTrace(err, 1)

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case *InternalError:
		typedErr.desc = i18nutil.NewText(fmt.Sprintf(format, args...))
		return typedErr
	}

	return &InternalError{
		cause: err,
		stack: stack,
		desc:  i18nutil.NewText(fmt.Sprintf(format, args...)),
	}
}

// WithDescriptionl creates a new *InternalError from the passed cause using
// the localized description.
//
// If the passed error is a *SilentError or an *InternalError, WithDescriptionl
// will unwrap it first.
func WithDescriptionl(err error, description *i18n.Config) error {
	if err == nil {
		return nil
	}

	stack := stackTrace(err, 1)

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case *InternalError:
		typedErr.desc = i18nutil.NewTextl(description)
		return typedErr
	}

	return &InternalError{
		cause: err,
		stack: stack,
		desc:  i18nutil.NewTextl(description),
	}
}

// WithDescriptionlt creates an internal error from the passed cause using the
// message generated from the passed term.
//
// If the passed error is a *SilentError or an *InternalError,
// WithDescriptionlt will unwrap it first.
func WithDescriptionlt(err error, description i18n.Term) error {
	if err == nil {
		return nil
	}

	stack := stackTrace(err, 1)

	switch typedErr := err.(type) { //nolint:errorlint
	case *SilentError:
		err = typedErr.Unwrap()
	case *InternalError:
		typedErr.desc = i18nutil.NewTextlt(description)
		return typedErr
	}

	return &InternalError{
		cause: err,
		stack: stack,
		desc:  i18nutil.NewTextlt(description),
	}
}

// Description returns the description of the error and localizes it, if
// possible.
// If there is no custom description or the custom description is empty,
// Description will fall back on the default description.
func (e *InternalError) Description(l *i18n.Localizer) string {
	if e.desc != nil {
		desc, err := e.desc.Get(l)
		if err == nil && len(desc) > 0 {
			return desc
		}
	}

	desc, _ := l.Localize(defaultInternalDesc)
	return desc
}

func (e *InternalError) Error() string         { return e.cause.Error() }
func (e *InternalError) Unwrap() error         { return e.cause }
func (e *InternalError) StackTrace() []uintptr { return e.stack }

// Handle handles the InternalError.
// By default it logs the error and sends out an internal error Embed.
func (e *InternalError) Handle(s *state.State, ctx *plugin.Context) error {
	// prevent infinite error cycle, by not allowing error returns
	HandleInternalError(e, s, ctx)

	return nil
}

var HandleInternalError = func(ierr *InternalError, s *state.State, ctx *plugin.Context) {
	if derr := discorderr.As(ierr.Unwrap()); derr != nil {
		switch {
		case discorderr.Is(derr, discorderr.InsufficientPermissions):
			// prevent cyclic error handling, in case this error was cause by
			// the same permission needed to handle the
			// InsufficientPermissionError
			_ = DefaultInsufficientPermissionsError.Handle(s, ctx)

			return
		case discorderr.Is(derr, discorderr.TemporarilyDisabled):
			ierr.desc = i18nutil.NewTextl(discordErrorFeatureTemporarilyDisabled)
		case derr.Status >= 500:
			ierr.desc = i18nutil.NewTextl(discordErrorServerError)
		}
	}

	logstract.
		WithFields(logstract.Fields{
			"cmd_ident": ctx.InvokedCommand.Identifier,
			"err":       ierr.cause,
		}).
		Error("internal error")

	embed := ErrorEmbed.Clone().
		WithSimpleTitlel(internalErrorTitle).
		WithDescription(ierr.Description(ctx.Localizer))

	_, _ = ctx.ReplyEmbedBuilder(embed)
}
