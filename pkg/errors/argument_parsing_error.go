package errors

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

// ArgumentParsingError is the error used if an argument or flag a user
// supplied is invalid.
type ArgumentParsingError struct {
	desc i18nutil.Text
}

var _ Interface = new(ArgumentParsingError)

// NewArgumentParsingError returns a new ArgumentParsingError with the passed
// description.
// The description mustn't be empty for this error to be handled properly.
func NewArgumentParsingError(description string) *ArgumentParsingError {
	return &ArgumentParsingError{
		desc: i18nutil.NewText(description),
	}
}

// NewArgumentParsingErrorl returns a new ArgumentParsingError using the passed
// i18n.Config to generate a description.
func NewArgumentParsingErrorl(description *i18n.Config) *ArgumentParsingError {
	return &ArgumentParsingError{
		desc: i18nutil.NewTextl(description),
	}
}

// NewArgumentParsingErrorlt returns a new ArgumentParsingError using the
// passed term to generate a description.
func NewArgumentParsingErrorlt(description i18n.Term) *ArgumentParsingError {
	return NewArgumentParsingErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *ArgumentParsingError) Description(l *i18n.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *ArgumentParsingError) Error() string { return "argument parsing error" }

// Handle send an error embed containing a description of which arg/flag was
// faulty and an optional reason for the error, in the channel the command
// was sent in.
func (e *ArgumentParsingError) Handle(_ *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := ErrorEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
