package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

var DefaultRestrictionError = NewRestrictionErrorl(defaultRestrictionDesc)

type RestrictionError struct {
	// description of the error, either is set
	descString string
	descConfig localization.Config
}

// NewRestrictionError creates a new RestrictionError with the passed
// description.
func NewRestrictionError(desc string) *RestrictionError {
	return &RestrictionError{
		descString: desc,
	}
}

// NewRestrictionErrorl creates a new RestrictionError using the message
// generated from the passed localization.Config as description.
func NewRestrictionErrorl(description localization.Config) *RestrictionError {
	return &RestrictionError{
		descConfig: description,
	}
}

// NewUserInfolt creates a new RestrictionError using the message generated
// from the passed term as description.
func NewRestrictionErrorlt(description localization.Term) *RestrictionError {
	return NewRestrictionErrorl(localization.Config{
		Term: description,
	})
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *RestrictionError) Description(l *localization.Localizer) (string, error) {
	if e.descString != "" {
		return e.descString, nil
	}

	return l.Localize(e.descConfig)
}

func (e *RestrictionError) Error() string { return "user error" }

func (e *RestrictionError) Is(target error) bool {
	casted, ok := target.(*RestrictionError)
	if !ok {
		return false
	}

	return e.descString == casted.descString || e.descConfig == casted.descConfig
}

// Handle sends an error embed with the description of the UserError.
func (e *RestrictionError) Handle(_ *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
