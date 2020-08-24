package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

// ArgumentParsingError is the error used if the arguments or flags a user
// supplied are invalid.
// It consists of two separate parts:
//
// The description is mandatory, and contains information about which argument,
// flag is affected, or similar information, such as signaling an error with
// the argument prefix.
//
// The reason is optional and is usually filled during parsing.
// It contains information about why this error occurred, and what can be done
// to fix it.
type ArgumentParsingError struct {
	descString string
	descConfig localization.Config

	reasonString string
	reasonConfig localization.Config
}

// NewArgumentParsingError returns a new ArgumentParsingError with the passed
// description.
// The description mustn't be empty for this error to be handled properly.
func NewArgumentParsingError(description string) *ArgumentParsingError {
	return &ArgumentParsingError{
		descString: description,
	}
}

// NewArgumentParsingErrorl returns a new ArgumentParsingError using the passed
// localization.Config to generate a description.
func NewArgumentParsingErrorl(description localization.Config) *ArgumentParsingError {
	return &ArgumentParsingError{
		descConfig: description,
	}
}

// NewArgumentParsingErrorlt returns a new ArgumentParsingError using the
// passed term to generate a description.
func NewArgumentParsingErrorlt(description localization.Term) *ArgumentParsingError {
	return NewArgumentParsingErrorl(description.AsConfig())
}

// WithReason creates a copy of the error and adds the passed reason to it.
func (e ArgumentParsingError) WithReason(reason string) *ArgumentParsingError {
	e.reasonString = reason
	return &e
}

// WithReasonl creates a copy of the error and adds the passed reason to it.
func (e ArgumentParsingError) WithReasonl(reason localization.Config) *ArgumentParsingError {
	e.reasonConfig = reason
	return &e
}

// WithReasonlt creates a copy of the error and adds the passed reason to it.
func (e ArgumentParsingError) WithReasonlt(reason localization.Term) *ArgumentParsingError {
	return e.WithReasonl(reason.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *ArgumentParsingError) Description(l *localization.Localizer) (string, error) {
	if e.descString != "" {
		return e.descString, nil
	}

	return l.Localize(e.descConfig)
}

// Reason returns the reason of the error and to localizes it, if
// possible.
// If there is no description, an empty string will be returned.
func (e *ArgumentParsingError) Reason(l *localization.Localizer) string {
	if e.reasonString != "" {
		return e.reasonString
	}

	reason, err := l.Localize(e.reasonConfig)
	if err != nil { // we have no reason
		return ""
	}

	return reason
}

func (e *ArgumentParsingError) Error() string { return "argument parsing error" }

func (e *ArgumentParsingError) Is(target error) bool {
	casted, ok := target.(*ArgumentParsingError)
	if !ok {
		return false
	}

	return e.descString == casted.descString || e.descConfig == casted.descConfig
}

// Handle send an error embed containing a description of which arg/flag was
// faulty and an optional reason for the error, in the channel the command
// was sent in.
func (e *ArgumentParsingError) Handle(_ *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(desc)

	if reasonVal := e.Reason(ctx.Localizer); reasonVal != "" {
		// we can ignore the error, as we have a fallback
		reasonName, _ := ctx.Localize(argumentParsingReasonFieldName)

		embed.WithField(reasonName, reasonVal)
	}

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
