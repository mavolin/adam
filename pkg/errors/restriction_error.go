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
func NewRestrictionErrorl(desc localization.Config) *RestrictionError {
	return &RestrictionError{
		descConfig: desc,
	}
}

// NewUserInfolt creates a new RestrictionError using the message generated
// from the passed term as description.
func NewRestrictionErrorlt(descTerm string) *RestrictionError {
	return NewRestrictionErrorl(localization.Config{
		Term: descTerm,
	})
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *RestrictionError) Description(l *localization.Localizer) (desc string) {
	if e.descString != "" {
		return e.descString
	}

	var err error
	if desc, err = l.Localize(e.descConfig); err != nil {
		// we can ignore the error, as there is a fallback
		desc, _ = l.Localize(defaultInternalDesc)
	}

	return desc
}

func (e *RestrictionError) Error() string { return "user error" }

// Handle sends an error embed with the description of the UserError.
func (e *RestrictionError) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(e.Description(ctx.Localizer))

	_, err = ctx.ReplyEmbedBuilder(embed)
	return
}
