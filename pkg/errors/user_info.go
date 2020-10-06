package errors

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// UserInfo is less sever error on the user-side.
// The error will reported to the user via a message containing a detailed
// description of the problem.
// It won't be logged
type UserInfo struct {
	embed *embedutil.Builder
}

// NewUserInfo creates a new UserInfo using the passed description.
// The description mustn't be empty for this error to be handled properly.
func NewUserInfo(description string) *UserInfo {
	return &UserInfo{
		embed: InfoEmbed.Clone().
			WithDescription(description),
	}
}

// NewUserInfol creates a new UserInfo using the message generated from the
// passed i18n.Config.
func NewUserInfol(description i18n.Config) *UserInfo {
	return &UserInfo{
		embed: InfoEmbed.Clone().
			WithDescriptionl(description),
	}
}

// NewUserInfolt creates a new UserInfo using the message generated from the
// passed term.
func NewUserInfolt(description i18n.Term) *UserInfo {
	return NewUserInfol(description.AsConfig())
}

// WithSimpleTitle adds a plain title (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleTitle(title string) *UserInfo {
	i.embed.WithSimpleTitle(title)
	return i
}

// WithSimpleTitlel adds a plain title (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleTitlel(title i18n.Config) *UserInfo {
	i.embed.WithSimpleTitlel(title)
	return i
}

// WithSimpleTitlelt adds a plain title (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleTitlelt(title i18n.Term) *UserInfo {
	return i.WithSimpleTitlel(title.AsConfig())
}

// WithTitle adds a title (max. 256 characters) with a link to the UserInfo.
func (i *UserInfo) WithTitle(title string, url discord.URL) *UserInfo {
	i.embed.WithTitle(title, url)
	return i
}

// WithTitlel adds a title (max. 256 characters) with a link to the UserInfo.
func (i *UserInfo) WithTitlel(title i18n.Config, url discord.URL) *UserInfo {
	i.embed.WithTitlel(title, url)
	return i
}

// WithTitlelt adds a title (max. 256 characters) with a link to the UserInfo.
func (i *UserInfo) WithTitlelt(title i18n.Term, url discord.URL) *UserInfo {
	return i.WithTitlel(title.AsConfig(), url)
}

// WithDescription adds a description (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithDescription(description string) *UserInfo {
	i.embed.WithDescription(description)
	return i
}

// WithDescriptionl adds a description (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithDescriptionl(description i18n.Config) *UserInfo {
	i.embed.WithDescriptionl(description)
	return i
}

// WithDescriptionlt adds a description (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithDescriptionlt(description i18n.Term) *UserInfo {
	return i.WithDescriptionl(description.AsConfig())
}

// WithTimestamp adds a discord.Timestamp to the UserInfo.
func (i *UserInfo) WithTimestamp(timestamp discord.Timestamp) *UserInfo {
	i.embed.WithTimestamp(timestamp)
	return i
}

// WithTimestamp adds a timestamp of the current time to the UserInfo.
func (i *UserInfo) WithTimestampNow() *UserInfo {
	return i.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color of the embed to the passed discord.Color.
func (i *UserInfo) WithColor(color discord.Color) *UserInfo {
	i.embed.WithColor(color)
	return i
}

// WithSimpleFooter adds a plain footer (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithSimpleFooter(text string) *UserInfo {
	i.embed.WithSimpleFooter(text)
	return i
}

// WithSimpleFooterl adds a plain footer (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithSimpleFooterl(text i18n.Config) *UserInfo {
	i.embed.WithSimpleFooterl(text)
	return i
}

// WithSimpleFooterlt adds a plain footer (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithSimpleFooterlt(text i18n.Term) *UserInfo {
	return i.WithSimpleFooterl(text.AsConfig())
}

// WithFooter adds a footer (max. 2048 character) with an icon to the UserInfo.
func (i *UserInfo) WithFooter(text string, icon discord.URL) *UserInfo {
	i.embed.WithField(text, icon)
	return i
}

// WithFooterl adds a footer (max. 2048 character) with an icon to the UserInfo.
func (i *UserInfo) WithFooterl(text i18n.Config, icon discord.URL) *UserInfo {
	i.embed.WithFooterl(text, icon)
	return i
}

// WithFooterlt adds a footer (max. 2048 character) with an icon to the UserInfo.
func (i *UserInfo) WithFooterlt(text i18n.Term, icon discord.URL) *UserInfo {
	return i.WithFooterl(text.AsConfig(), icon)
}

// WithImage adds an image to the UserInfo.
func (i *UserInfo) WithImage(image discord.URL) *UserInfo {
	i.embed.WithImage(image)
	return i
}

// WithThumbnail adds a thumbnail to the UserInfo.
func (i *UserInfo) WithThumbnail(thumbnail discord.URL) *UserInfo {
	i.embed.WithThumbnail(thumbnail)
	return i
}

// WithSimpleAuthor adds a plain author (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleAuthor(name string) *UserInfo {
	i.embed.WithSimpleAuthor(name)
	return i
}

// WithSimpleAuthorl adds a plain author (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleAuthorl(name i18n.Config) *UserInfo {
	i.embed.WithSimpleAuthorl(name)
	return i
}

// WithSimpleAuthorlt adds a plain author (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleAuthorlt(name i18n.Term) *UserInfo {
	return i.WithSimpleAuthorl(name.AsConfig())
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the embed.
func (i *UserInfo) WithSimpleAuthorWithURL(name string, url discord.URL) *UserInfo {
	i.embed.WithSimpleAuthorWithURL(name, url)
	return i
}

// WithSimpleAuthorWithURLl adds an author (max. 256 character) with a URL to
// the embed.
func (i *UserInfo) WithSimpleAuthorWithURLl(name i18n.Config, url discord.URL) *UserInfo {
	i.embed.WithSimpleAuthorWithURLl(name, url)
	return i
}

// WithSimpleAuthorWithURLlt adds an author (max. 256 character) with a URL to
// the embed.
func (i *UserInfo) WithSimpleAuthorWithURLlt(name i18n.Term, url discord.URL) *UserInfo {
	return i.WithSimpleAuthorWithURLl(name.AsConfig(), url)
}

// WithAuthor adds an author (max 256 characters) with an icon to the UserInfo.
func (i *UserInfo) WithAuthor(name string, icon discord.URL) *UserInfo {
	i.embed.WithAuthor(name, icon)
	return i
}

// WithAuthorl adds an author (max 256 characters) with an icon to the UserInfo.
func (i *UserInfo) WithAuthorl(name i18n.Config, icon discord.URL) *UserInfo {
	i.embed.WithAuthorl(name, icon)
	return i
}

// WithAuthorlt adds an author (max 256 characters) with an icon to the UserInfo.
func (i *UserInfo) WithAuthorlt(name i18n.Term, icon discord.URL) *UserInfo {
	return i.WithAuthorl(name.AsConfig(), icon)
}

// WithAuthorWithURL adds an author (max 256 characters) with an icon and a URL
// to the UserInfo.
func (i *UserInfo) WithAuthorWithURL(name string, icon, url discord.URL) *UserInfo {
	i.embed.WithAuthorWithURL(name, icon, url)
	return i
}

// WithAuthorWithURLl adds an author (max 256 characters) with an icon and a
// URL to the UserInfo.
func (i *UserInfo) WithAuthorWithURLl(name i18n.Config, icon, url discord.URL) *UserInfo {
	i.embed.WithAuthorWithURLl(name, icon, url)
	return i
}

// WithAuthorWithURLlt adds an author (max 256 characters) with an icon and a
// URL to the UserInfo.
func (i *UserInfo) WithAuthorWithURLlt(name i18n.Term, icon, url discord.URL) *UserInfo {
	return i.WithAuthorWithURLl(name.AsConfig(), icon, url)
}

// WithField adds the passed field to the UserInfo, and returns a pointer to
// the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithField(name, value string) *UserInfo {
	i.embed.WithField(name, value)
	return i
}

// WithFieldl adds the passed field to the UserInfo, and returns a pointer to
// the
// UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithFieldl(name, value i18n.Config) *UserInfo {
	i.embed.WithFieldl(name, value)
	return i
}

// WithFieldlt adds the passed field to the UserInfo, and returns a pointer to
// the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithFieldlt(name, value i18n.Term) *UserInfo {
	return i.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithField adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedField(name, value string) *UserInfo {
	i.embed.WithInlinedField(name, value)
	return i
}

// WithFieldl adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedFieldl(name, value i18n.Config) *UserInfo {
	i.embed.WithInlinedFieldl(name, value)
	return i
}

// WithFieldlt adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedFieldlt(name, value i18n.Term) *UserInfo {
	return i.WithFieldl(name.AsConfig(), value.AsConfig())
}

// Embed returns the embed of the UserInfo.
func (i *UserInfo) Embed(l *i18n.Localizer) (discord.Embed, error) {
	return i.embed.Build(l)
}

func (i *UserInfo) Error() string { return "user info" }

// Handle sends an info embed with the description of the UserInfo.
func (i *UserInfo) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	_, err = ctx.ReplyEmbedBuilder(i.embed)
	return
}
