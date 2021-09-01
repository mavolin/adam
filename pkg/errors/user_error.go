package errors //nolint:dupl

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/internal/embedbuilder"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// UserError is an error on the user-side.
// The error will be reported via a message containing a detailed description
// of the problem.
// It won't be logged.
type UserError struct {
	Embed *embedbuilder.Builder
}

var _ Error = new(UserError)

// NewCustomUserError creates a new *UserError using a NewErrorEmbed as a
// template.
func NewCustomUserError() *UserError {
	return &UserError{Embed: NewErrorEmbed()}
}

// NewUserErrorFromEmbed creates a new *UserError from the passed
// *msgbuilder.EmbedBuilder.
func NewUserErrorFromEmbed(e *embedbuilder.Builder) *UserError {
	return &UserError{Embed: e}
}

// NewUserError creates a new *UserError with the passed description using a
// NewErrorEmbed as template.
// The description mustn't be empty for this error to be handled properly.
func NewUserError(description string) *UserError {
	return NewCustomUserError().
		WithDescription(description)
}

// NewUserErrorl creates a new *UserError using the message generated from the
// passed *i18n.Config as description.
func NewUserErrorl(description *i18n.Config) *UserError {
	return NewCustomUserError().
		WithDescriptionl(description)
}

// NewUserErrorlt creates a new *UserError using the message generated from the
// passed term as description.
func NewUserErrorlt(description i18n.Term) *UserError {
	return NewUserErrorl(description.AsConfig())
}

// WithTitle sets the title (max. 256 characters) to the passed title.
func (e *UserError) WithTitle(title string) *UserError {
	return e.WithTitlel(i18n.NewStaticConfig(title))
}

// WithTitlelt sets the title (max. 256 characters) to the passed title.
func (e *UserError) WithTitlelt(title i18n.Term) *UserError {
	return e.WithTitlel(title.AsConfig())
}

// WithTitlel sets the title (max. 256 characters) to the passed title.
func (e *UserError) WithTitlel(title *i18n.Config) *UserError {
	e.Embed.WithTitlel(title)
	return e
}

// WithTitleURL assigns a discord.URL to the title.
func (e *UserError) WithTitleURL(url discord.URL) *UserError {
	e.Embed.WithTitleURL(url)
	return e
}

// WithDescription sets the description (max. 2048 characters) to the passed
// description.
func (e *UserError) WithDescription(description string) *UserError {
	return e.WithDescriptionl(i18n.NewStaticConfig(description))
}

// WithDescriptionlt sets the description (max. 2048 characters) to the passed
// description.
func (e *UserError) WithDescriptionlt(description i18n.Term) *UserError {
	return e.WithDescriptionl(description.AsConfig())
}

// WithDescriptionl sets the description (max. 2048 characters) to the passed
// description.
func (e *UserError) WithDescriptionl(description *i18n.Config) *UserError {
	e.Embed.WithDescriptionl(description)
	return e
}

// WithTimestamp sets the timestamp to the passed discord.Timestamp.
func (e *UserError) WithTimestamp(timestamp discord.Timestamp) *UserError {
	e.Embed.WithTimestamp(timestamp)
	return e
}

// WithTimestampNow sets the timestamp to a timestamp of the current time.
func (e *UserError) WithTimestampNow() *UserError {
	return e.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color to the passed discord.Color.
func (e *UserError) WithColor(color discord.Color) *UserError {
	e.Embed.WithColor(color)
	return e
}

// WithFooter sets the text of the footer (max. 2048 characters) to the passed
// text.
func (e *UserError) WithFooter(text string) *UserError {
	return e.WithFooterl(i18n.NewStaticConfig(text))
}

// WithFooterlt sets the text of the footer (max. 2048 characters) to the
// passed text.
func (e *UserError) WithFooterlt(text i18n.Term) *UserError {
	return e.WithFooterl(text.AsConfig())
}

// WithFooterl sets the text of the footer (max. 2048 characters) to the passed
// text.
func (e *UserError) WithFooterl(text *i18n.Config) *UserError {
	e.Embed.WithFooterl(text)
	return e
}

// WithFooterIcon sets the icon of the footer to the passed icon url.
func (e *UserError) WithFooterIcon(icon discord.URL) *UserError {
	e.Embed.WithFooterIcon(icon)
	return e
}

// WithImage sets the image to the passed image url.
func (e *UserError) WithImage(image discord.URL) *UserError {
	e.Embed.WithImage(image)
	return e
}

// WithThumbnail adds a thumbnail to the error.
func (e *UserError) WithThumbnail(thumbnail discord.URL) *UserError {
	e.Embed.WithThumbnail(thumbnail)
	return e
}

// WithAuthor sets the author's name (max. 256 characters) to the passed
// name.
func (e *UserError) WithAuthor(name string) *UserError {
	return e.WithAuthorl(i18n.NewStaticConfig(name))
}

// WithAuthorlt sets the author's name (max. 256 characters) to the passed
// name.
func (e *UserError) WithAuthorlt(name i18n.Term) *UserError {
	return e.WithAuthorl(name.AsConfig())
}

// WithAuthorl sets the author's name (max. 256 characters) to the passed
// name.
func (e *UserError) WithAuthorl(name *i18n.Config) *UserError {
	e.Embed.WithAuthorl(name)
	return e
}

// WithAuthorURL assigns the author the passed discord.URL.
func (e *UserError) WithAuthorURL(url discord.URL) *UserError {
	e.Embed.WithAuthorURL(url)
	return e
}

// WithAuthorIcon sets the icon of the author to the passed icon url.
func (e *UserError) WithAuthorIcon(icon discord.URL) *UserError {
	e.Embed.WithAuthorIcon(icon)
	return e
}

// WithField adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (e *UserError) WithField(name, value string) *UserError {
	return e.WithFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithFieldlt adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (e *UserError) WithFieldlt(name, value i18n.Term) *UserError {
	return e.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithFieldl adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (e *UserError) WithFieldl(name, value *i18n.Config) *UserError {
	e.Embed.WithFieldl(name, value)
	return e
}

// WithInlinedField adds an inlined field (name: max. 256 characters, value:
// max 1024 characters) to the embed.
func (e *UserError) WithInlinedField(name, value string) *UserError {
	return e.WithInlinedFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithInlinedFieldlt adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (e *UserError) WithInlinedFieldlt(name, value i18n.Term) *UserError {
	return e.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithInlinedFieldl adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (e *UserError) WithInlinedFieldl(name, value *i18n.Config) *UserError {
	e.Embed.WithInlinedFieldl(name, value)
	return e
}

func (e *UserError) Error() string { return "user error" }

// Handle handles the UserError.
// By default it sends the error Embed.
func (e *UserError) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleUserError(s, ctx, e)
}

var HandleUserError = func(s *state.State, ctx *plugin.Context, uerr *UserError) error {
	_, err := ctx.ReplyEmbedBuilders(uerr.Embed)
	return err
}
