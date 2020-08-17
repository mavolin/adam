package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

// UserError is an error on the user-side.
// The user will reported via a message containing a detailed description of
// the problem.
// The error won't be logged or captured by sentry.
type UserError struct {
	// description of the error, either is set
	descString string
	descConfig localization.Config
}

// NewUserError creates a new UserError with the passed description.
func NewUserError(desc string) *UserError {
	return &UserError{
		descString: desc,
	}
}

// NewUserErrorl creates a new UserError using the message generated from the
// passed localization.Config as description.
func NewUserErrorl(description localization.Config) *UserError {
	return &UserError{
		descConfig: description,
	}
}

// NewUserErrorlt creates a new UserError using the message generated from the
// passed term as description.
func NewUserErrorlt(description localization.Term) *UserError {
	return NewUserErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *UserError) Description(l *localization.Localizer) (string, error) {
	if e.descString != "" {
		return e.descString, nil
	}

	return l.Localize(e.descConfig)
}

func (e *UserError) Error() string { return "user error" }

// Handle sends an error embed with the description of the UserError.
func (e *UserError) Handle(_ *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
