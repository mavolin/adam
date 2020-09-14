package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

// DefaultRestrictionError is a restriction error with a default, generic
// description.
var DefaultFatalRestrictionError = NewFatalRestrictionErrorl(defaultRestrictionDesc)

// FatalRestrictionError is like a RestrictionError, but it indicates that a
// command should not be displayed in the help message.
// The error messages are the same
type FatalRestrictionError struct {
	// description of the error, either is set
	descString string
	descConfig localization.Config
}

// NewFatalRestrictionError creates a new FatalRestrictionError with the passed
// description.
func NewFatalRestrictionError(desc string) *FatalRestrictionError {
	return &FatalRestrictionError{
		descString: desc,
	}
}

// NewFatalRestrictionErrorl creates a new FatalRestrictionError using the message
// generated from the passed localization.Config as description.
func NewFatalRestrictionErrorl(description localization.Config) *FatalRestrictionError {
	return &FatalRestrictionError{
		descConfig: description,
	}
}

// NewFatalRestrictionErrorlt creates a new FatalRestrictionError using the message generated
// from the passed term as description.
func NewFatalRestrictionErrorlt(description localization.Term) *FatalRestrictionError {
	return NewFatalRestrictionErrorl(localization.Config{
		Term: description,
	})
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *FatalRestrictionError) Description(l *localization.Localizer) (string, error) {
	if e.descString != "" {
		return e.descString, nil
	}

	return l.Localize(e.descConfig)
}

func (e *FatalRestrictionError) Error() string { return "restriction error" }

// Handle sends an error embed with the description of the UserError.
func (e *FatalRestrictionError) Handle(_ *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := ErrorEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
