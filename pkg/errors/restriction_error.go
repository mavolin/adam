package errors

import (
	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

// DefaultRestrictionError is a restriction error with a default, generic
// description.
var DefaultRestrictionError = NewRestrictionErrorl(defaultRestrictionDesc)

// RestrictionError is the error returned if a restriction fails.
// It contains a description stating the conditions that need to be fulfilled
// for a command to execute.
//
// Besides restrictions, this will also be returned, if a user invokes the
// command in a channel, that is not specified in the plugin.Meta's
// ChannelTypes.
//
// Note that all mentions except the mention of the message author are
// suppressed.
// This allows you to more easily communicate errors with users etc. without
// any unintended pings.
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
func (e *RestrictionError) Handle(s *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	builder := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(desc)

	embed, err := builder.Build(ctx.Localizer)
	if err != nil {
		return err
	}

	_, err = s.SendMessageComplex(ctx.ChannelID, api.SendMessageData{
		Embed: &embed,
		AllowedMentions: &api.AllowedMentions{
			Users: []discord.UserID{ctx.Author.ID},
		},
	})

	return err
}
