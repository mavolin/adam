package errors

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

var (
	// DefaultRestrictionError is a restriction error with a default, generic
	// description.
	DefaultRestrictionError = NewRestrictionErrorl(defaultRestrictionDesc)
	// DefaultFatalRestrictionError is a restriction error with a default,
	// generic description and Fatal set to true.
	DefaultFatalRestrictionError = NewFatalRestrictionErrorl(defaultRestrictionDesc)
)

// RestrictionError is the error returned if a restriction fails.
// It contains a description stating the conditions that need to be fulfilled
// for a command to execute.
//
// Besides restrictions, this will also be returned, if a user invokes the
// command in a channel, that is not specified in the plugin.Meta's
// ChannelTypes.
//
// Note that the description might contain mentions, which are intended not
// to ping anyone, e.g. "You need @role to use this command.".
// This means you should use allowed mentions if you are custom handling this
// error and not using an embed, which suppresses mentions by default.
type RestrictionError struct {
	// description of the error
	desc i18nutil.Text

	// Fatal defines if the RestrictionError is fatal.
	// Fatal errors won't be shown in the help message.
	Fatal bool
}

var _ Interface = new(RestrictionError)

// NewRestrictionError creates a new RestrictionError with the passed
// description.
func NewRestrictionError(description string) *RestrictionError {
	return &RestrictionError{
		desc: i18nutil.NewText(description),
	}
}

// NewRestrictionErrorl creates a new RestrictionError using the message
// generated from the passed i18n.Config as description.
func NewRestrictionErrorl(description *i18n.Config) *RestrictionError {
	return &RestrictionError{
		desc: i18nutil.NewTextl(description),
	}
}

// NewRestrictionErrorlt creates a new RestrictionError using the message
// generated from the passed term as description.
func NewRestrictionErrorlt(description i18n.Term) *RestrictionError {
	return NewRestrictionErrorl(description.AsConfig())
}

// NewFatalRestrictionError creates a new fatal RestrictionError with the
// passed description.
func NewFatalRestrictionError(description string) *RestrictionError {
	return &RestrictionError{
		desc:  i18nutil.NewText(description),
		Fatal: true,
	}
}

// NewFatalRestrictionErrorl creates a new fatal RestrictionError using the
// message generated from the passed i18n.Config as description.
func NewFatalRestrictionErrorl(description *i18n.Config) *RestrictionError {
	return &RestrictionError{
		desc:  i18nutil.NewTextl(description),
		Fatal: true,
	}
}

// NewFatalRestrictionErrorlt creates a new fatal RestrictionError using the
// message generated from the passed term as description.
func NewFatalRestrictionErrorlt(description i18n.Term) *RestrictionError {
	return NewFatalRestrictionErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *RestrictionError) Description(l *i18n.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *RestrictionError) Error() string { return "restriction error" }

// Handle sends an error embed with the description of the ThrottlingError.
func (e *RestrictionError) Handle(_ *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := ErrorEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
