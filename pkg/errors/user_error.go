package errors

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// UserError is an error on the user-side.
// The user will reported via a message containing a detailed description of
// the problem.
// The error won't be logged.
type UserError struct {
	embed *embedutil.Builder
}

// NewCustomUserError creates a new UserError using the ErrorEmbed as a
// template.
func NewCustomUserError() *UserError {
	return &UserError{
		embed: ErrorEmbed.Clone(),
	}
}

// NewUserError creates a new UserError with the passed description.
// The description mustn't be empty for this error to be handled properly.
func NewUserError(description string) *UserError {
	return NewCustomUserError().
		WithDescription(description)
}

// NewUserErrorl creates a new UserError using the message generated from the
// passed i18n.Config as description.
func NewUserErrorl(description i18n.Config) *UserError {
	return NewCustomUserError().
		WithDescriptionl(description)
}

// NewUserErrorlt creates a new UserError using the message generated from the
// passed term as description.
func NewUserErrorlt(description i18n.Term) *UserError {
	return NewUserErrorl(description.AsConfig())
}

// WithSimpleTitle adds a plain title (max. 256 characters) to the error.
func (e *UserError) WithSimpleTitle(title string) *UserError {
	e.embed.WithSimpleTitle(title)
	return e
}

// WithSimpleTitlel adds a plain title (max. 256 characters) to the error.
func (e *UserError) WithSimpleTitlel(title i18n.Config) *UserError {
	e.embed.WithSimpleTitlel(title)
	return e
}

// WithSimpleTitlelt adds a plain title (max. 256 characters) to the error.
func (e *UserError) WithSimpleTitlelt(title i18n.Term) *UserError {
	return e.WithSimpleTitlel(title.AsConfig())
}

// WithTitle adds a title (max. 256 characters) with a link to the error.
func (e *UserError) WithTitle(title string, url discord.URL) *UserError {
	e.embed.WithTitle(title, url)
	return e
}

// WithTitlel adds a title (max. 256 characters) with a link to the error.
func (e *UserError) WithTitlel(title i18n.Config, url discord.URL) *UserError {
	e.embed.WithTitlel(title, url)
	return e
}

// WithTitlelt adds a title (max. 256 characters) with a link to the error.
func (e *UserError) WithTitlelt(title i18n.Term, url discord.URL) *UserError {
	return e.WithTitlel(title.AsConfig(), url)
}

// WithDescription adds a description (max. 2048 characters) to the error.
func (e *UserError) WithDescription(description string) *UserError {
	e.embed.WithDescription(description)
	return e
}

// WithDescriptionl adds a description (max. 2048 characters) to the error.
func (e *UserError) WithDescriptionl(description i18n.Config) *UserError {
	e.embed.WithDescriptionl(description)
	return e
}

// WithDescriptionlt adds a description (max. 2048 characters) to the error.
func (e *UserError) WithDescriptionlt(description i18n.Term) *UserError {
	return e.WithDescriptionl(description.AsConfig())
}

// WithTimestamp adds a discord.Timestamp to the error.
func (e *UserError) WithTimestamp(timestamp discord.Timestamp) *UserError {
	e.embed.WithTimestamp(timestamp)
	return e
}

// WithTimestamp adds a timestamp of the current time to the error.
func (e *UserError) WithTimestampNow() *UserError {
	return e.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color of the embed to the passed discord.Color.
func (e *UserError) WithColor(color discord.Color) *UserError {
	e.embed.WithColor(color)
	return e
}

// WithSimpleFooter adds a plain footer (max. 2048 characters) to the error.
func (e *UserError) WithSimpleFooter(text string) *UserError {
	e.embed.WithSimpleFooter(text)
	return e
}

// WithSimpleFooterl adds a plain footer (max. 2048 characters) to the error.
func (e *UserError) WithSimpleFooterl(text i18n.Config) *UserError {
	e.embed.WithSimpleFooterl(text)
	return e
}

// WithSimpleFooterlt adds a plain footer (max. 2048 characters) to the error.
func (e *UserError) WithSimpleFooterlt(text i18n.Term) *UserError {
	return e.WithSimpleFooterl(text.AsConfig())
}

// WithFooter adds a footer (max. 2048 character) with an icon to the error.
func (e *UserError) WithFooter(text string, icon discord.URL) *UserError {
	e.embed.WithField(text, icon)
	return e
}

// WithFooterl adds a footer (max. 2048 character) with an icon to the error.
func (e *UserError) WithFooterl(text i18n.Config, icon discord.URL) *UserError {
	e.embed.WithFooterl(text, icon)
	return e
}

// WithFooterlt adds a footer (max. 2048 character) with an icon to the error.
func (e *UserError) WithFooterlt(text i18n.Term, icon discord.URL) *UserError {
	return e.WithFooterl(text.AsConfig(), icon)
}

// WithImage adds an image to the error.
func (e *UserError) WithImage(image discord.URL) *UserError {
	e.embed.WithImage(image)
	return e
}

// WithThumbnail adds a thumbnail to the error.
func (e *UserError) WithThumbnail(thumbnail discord.URL) *UserError {
	e.embed.WithThumbnail(thumbnail)
	return e
}

// WithSimpleAuthor adds a plain author (max. 256 characters) to the error.
func (e *UserError) WithSimpleAuthor(name string) *UserError {
	e.embed.WithSimpleAuthor(name)
	return e
}

// WithSimpleAuthorl adds a plain author (max. 256 characters) to the error.
func (e *UserError) WithSimpleAuthorl(name i18n.Config) *UserError {
	e.embed.WithSimpleAuthorl(name)
	return e
}

// WithSimpleAuthorlt adds a plain author (max. 256 characters) to the error.
func (e *UserError) WithSimpleAuthorlt(name i18n.Term) *UserError {
	return e.WithSimpleAuthorl(name.AsConfig())
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the embed.
func (e *UserError) WithSimpleAuthorWithURL(name string, url discord.URL) *UserError {
	e.embed.WithSimpleAuthorWithURL(name, url)
	return e
}

// WithSimpleAuthorWithURLl adds an author (max. 256 character) with a URL to
// the embed.
func (e *UserError) WithSimpleAuthorWithURLl(name i18n.Config, url discord.URL) *UserError {
	e.embed.WithSimpleAuthorWithURLl(name, url)
	return e
}

// WithSimpleAuthorWithURLlt adds an author (max. 256 character) with a URL to
// the embed.
func (e *UserError) WithSimpleAuthorWithURLlt(name i18n.Term, url discord.URL) *UserError {
	return e.WithSimpleAuthorWithURLl(name.AsConfig(), url)
}

// WithAuthor adds an author (max 256 characters) with an icon to the error.
func (e *UserError) WithAuthor(name string, icon discord.URL) *UserError {
	e.embed.WithAuthor(name, icon)
	return e
}

// WithAuthorl adds an author (max 256 characters) with an icon to the error.
func (e *UserError) WithAuthorl(name i18n.Config, icon discord.URL) *UserError {
	e.embed.WithAuthorl(name, icon)
	return e
}

// WithAuthorlt adds an author (max 256 characters) with an icon to the error.
func (e *UserError) WithAuthorlt(name i18n.Term, icon discord.URL) *UserError {
	return e.WithAuthorl(name.AsConfig(), icon)
}

// WithAuthorWithURL adds an author (max 256 characters) with an icon and a URL
// to the error.
func (e *UserError) WithAuthorWithURL(name string, icon, url discord.URL) *UserError {
	e.embed.WithAuthorWithURL(name, icon, url)
	return e
}

// WithAuthorWithURLl adds an author (max 256 characters) with an icon and a
// URL to the error.
func (e *UserError) WithAuthorWithURLl(name i18n.Config, icon, url discord.URL) *UserError {
	e.embed.WithAuthorWithURLl(name, icon, url)
	return e
}

// WithAuthorWithURLlt adds an author (max 256 characters) with an icon and a
// URL to the error.
func (e *UserError) WithAuthorWithURLlt(name i18n.Term, icon, url discord.URL) *UserError {
	return e.WithAuthorWithURLl(name.AsConfig(), icon, url)
}

// WithField adds the passed field to the error, and returns a pointer to the
// UserError to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (e *UserError) WithField(name, value string) *UserError {
	e.embed.WithField(name, value)
	return e
}

// WithFieldl adds the passed field to the error, and returns a pointer to the
// UserError to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (e *UserError) WithFieldl(name, value i18n.Config) *UserError {
	e.embed.WithFieldl(name, value)
	return e
}

// WithFieldlt adds the passed field to the error, and returns a pointer to the
// UserError to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (e *UserError) WithFieldlt(name, value i18n.Term) *UserError {
	return e.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithField adds the passed inlined field to the error, and returns a pointer
// to the UserError to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (e *UserError) WithInlinedField(name, value string) *UserError {
	e.embed.WithInlinedField(name, value)
	return e
}

// WithFieldl adds the passed inlined field to the error, and returns a pointer
// to the UserError to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (e *UserError) WithInlinedFieldl(name, value i18n.Config) *UserError {
	e.embed.WithInlinedFieldl(name, value)
	return e
}

// WithFieldlt adds the passed inlined field to the error, and returns a
// pointer to the UserError to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (e *UserError) WithInlinedFieldlt(name, value i18n.Term) *UserError {
	return e.WithFieldl(name.AsConfig(), value.AsConfig())
}

// Embed returns the embed of the UserError.
func (e *UserError) Embed(l *i18n.Localizer) (discord.Embed, error) {
	return e.embed.Build(l)
}

func (e *UserError) Error() string { return "user error" }

// Handle sends an error embed with the description of the UserError.
func (e *UserError) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	_, err = ctx.ReplyEmbedBuilder(e.embed)
	return
}
