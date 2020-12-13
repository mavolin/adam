package errors //nolint:dupl

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
	Embed *embedutil.Builder
}

var _ Error = new(UserInfo)

// NewCustomUserInfo creates a new UserInfo using the InfoEmbed as template.
func NewCustomUserInfo() *UserInfo {
	return &UserInfo{Embed: InfoEmbed.Clone()}
}

// NewUserInfoFromEmbed creates a new UserInfo from the passed
// embedutil.Builder.
func NewUserInfoFromEmbed(e *embedutil.Builder) *UserInfo {
	return &UserInfo{Embed: e}
}

// NewUserInfo creates a new UserInfo using the passed description.
// The description mustn't be empty for this error to be handled properly.
func NewUserInfo(description string) *UserInfo {
	return NewCustomUserInfo().
		WithDescription(description)
}

// NewUserInfol creates a new UserInfo using the message generated from the
// passed i18n.Config.
func NewUserInfol(description *i18n.Config) *UserInfo {
	return NewCustomUserInfo().
		WithDescriptionl(description)
}

// NewUserInfolt creates a new UserInfo using the message generated from the
// passed term.
func NewUserInfolt(description i18n.Term) *UserInfo {
	return NewUserInfol(description.AsConfig())
}

// WithSimpleTitle adds a plain title (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleTitle(title string) *UserInfo {
	i.Embed.WithSimpleTitle(title)
	return i
}

// WithSimpleTitlel adds a plain title (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleTitlel(title *i18n.Config) *UserInfo {
	i.Embed.WithSimpleTitlel(title)
	return i
}

// WithSimpleTitlelt adds a plain title (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleTitlelt(title i18n.Term) *UserInfo {
	return i.WithSimpleTitlel(title.AsConfig())
}

// WithTitle adds a title (max. 256 characters) with a link to the UserInfo.
func (i *UserInfo) WithTitle(title string, url discord.URL) *UserInfo {
	i.Embed.WithTitle(title, url)
	return i
}

// WithTitlel adds a title (max. 256 characters) with a link to the UserInfo.
func (i *UserInfo) WithTitlel(title *i18n.Config, url discord.URL) *UserInfo {
	i.Embed.WithTitlel(title, url)
	return i
}

// WithTitlelt adds a title (max. 256 characters) with a link to the UserInfo.
func (i *UserInfo) WithTitlelt(title i18n.Term, url discord.URL) *UserInfo {
	return i.WithTitlel(title.AsConfig(), url)
}

// WithDescription adds a description (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithDescription(description string) *UserInfo {
	i.Embed.WithDescription(description)
	return i
}

// WithDescriptionl adds a description (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithDescriptionl(description *i18n.Config) *UserInfo {
	i.Embed.WithDescriptionl(description)
	return i
}

// WithDescriptionlt adds a description (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithDescriptionlt(description i18n.Term) *UserInfo {
	return i.WithDescriptionl(description.AsConfig())
}

// WithTimestamp adds a discord.Timestamp to the UserInfo.
func (i *UserInfo) WithTimestamp(timestamp discord.Timestamp) *UserInfo {
	i.Embed.WithTimestamp(timestamp)
	return i
}

// WithTimestamp adds a timestamp of the current time to the UserInfo.
func (i *UserInfo) WithTimestampNow() *UserInfo {
	return i.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color of the Embed to the passed discord.Color.
func (i *UserInfo) WithColor(color discord.Color) *UserInfo {
	i.Embed.WithColor(color)
	return i
}

// WithSimpleFooter adds a plain footer (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithSimpleFooter(text string) *UserInfo {
	i.Embed.WithSimpleFooter(text)
	return i
}

// WithSimpleFooterl adds a plain footer (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithSimpleFooterl(text *i18n.Config) *UserInfo {
	i.Embed.WithSimpleFooterl(text)
	return i
}

// WithSimpleFooterlt adds a plain footer (max. 2048 characters) to the UserInfo.
func (i *UserInfo) WithSimpleFooterlt(text i18n.Term) *UserInfo {
	return i.WithSimpleFooterl(text.AsConfig())
}

// WithFooter adds a footer (max. 2048 character) with an icon to the UserInfo.
func (i *UserInfo) WithFooter(text string, icon discord.URL) *UserInfo {
	i.Embed.WithField(text, icon)
	return i
}

// WithFooterl adds a footer (max. 2048 character) with an icon to the UserInfo.
func (i *UserInfo) WithFooterl(text *i18n.Config, icon discord.URL) *UserInfo {
	i.Embed.WithFooterl(text, icon)
	return i
}

// WithFooterlt adds a footer (max. 2048 character) with an icon to the UserInfo.
func (i *UserInfo) WithFooterlt(text i18n.Term, icon discord.URL) *UserInfo {
	return i.WithFooterl(text.AsConfig(), icon)
}

// WithImage adds an image to the UserInfo.
func (i *UserInfo) WithImage(image discord.URL) *UserInfo {
	i.Embed.WithImage(image)
	return i
}

// WithThumbnail adds a thumbnail to the UserInfo.
func (i *UserInfo) WithThumbnail(thumbnail discord.URL) *UserInfo {
	i.Embed.WithThumbnail(thumbnail)
	return i
}

// WithSimpleAuthor adds a plain author (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleAuthor(name string) *UserInfo {
	i.Embed.WithSimpleAuthor(name)
	return i
}

// WithSimpleAuthorl adds a plain author (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleAuthorl(name *i18n.Config) *UserInfo {
	i.Embed.WithSimpleAuthorl(name)
	return i
}

// WithSimpleAuthorlt adds a plain author (max. 256 characters) to the UserInfo.
func (i *UserInfo) WithSimpleAuthorlt(name i18n.Term) *UserInfo {
	return i.WithSimpleAuthorl(name.AsConfig())
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the Embed.
func (i *UserInfo) WithSimpleAuthorWithURL(name string, url discord.URL) *UserInfo {
	i.Embed.WithSimpleAuthorWithURL(name, url)
	return i
}

// WithSimpleAuthorWithURLl adds an author (max. 256 character) with a URL to
// the Embed.
func (i *UserInfo) WithSimpleAuthorWithURLl(name *i18n.Config, url discord.URL) *UserInfo {
	i.Embed.WithSimpleAuthorWithURLl(name, url)
	return i
}

// WithSimpleAuthorWithURLlt adds an author (max. 256 character) with a URL to
// the Embed.
func (i *UserInfo) WithSimpleAuthorWithURLlt(name i18n.Term, url discord.URL) *UserInfo {
	return i.WithSimpleAuthorWithURLl(name.AsConfig(), url)
}

// WithAuthor adds an author (max 256 characters) with an icon to the UserInfo.
func (i *UserInfo) WithAuthor(name string, icon discord.URL) *UserInfo {
	i.Embed.WithAuthor(name, icon)
	return i
}

// WithAuthorl adds an author (max 256 characters) with an icon to the UserInfo.
func (i *UserInfo) WithAuthorl(name *i18n.Config, icon discord.URL) *UserInfo {
	i.Embed.WithAuthorl(name, icon)
	return i
}

// WithAuthorlt adds an author (max 256 characters) with an icon to the UserInfo.
func (i *UserInfo) WithAuthorlt(name i18n.Term, icon discord.URL) *UserInfo {
	return i.WithAuthorl(name.AsConfig(), icon)
}

// WithAuthorWithURL adds an author (max 256 characters) with an icon and a URL
// to the UserInfo.
func (i *UserInfo) WithAuthorWithURL(name string, icon, url discord.URL) *UserInfo {
	i.Embed.WithAuthorWithURL(name, icon, url)
	return i
}

// WithAuthorWithURLl adds an author (max 256 characters) with an icon and a
// URL to the UserInfo.
func (i *UserInfo) WithAuthorWithURLl(name *i18n.Config, icon, url discord.URL) *UserInfo {
	i.Embed.WithAuthorWithURLl(name, icon, url)
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
	i.Embed.WithField(name, value)
	return i
}

// WithFieldl adds the passed field to the UserInfo, and returns a pointer to
// the
// UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithFieldl(name, value *i18n.Config) *UserInfo {
	i.Embed.WithFieldl(name, value)
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
	i.Embed.WithInlinedField(name, value)
	return i
}

// WithFieldl adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedFieldl(name, value *i18n.Config) *UserInfo {
	i.Embed.WithInlinedFieldl(name, value)
	return i
}

// WithFieldlt adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedFieldlt(name, value i18n.Term) *UserInfo {
	return i.WithFieldl(name.AsConfig(), value.AsConfig())
}

func (i *UserInfo) Error() string { return "user info" }

// Handle sends the info Embed.
func (i *UserInfo) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleUserInfo(i, s, ctx)
}

var HandleUserInfo = func(info *UserInfo, s *state.State, ctx *plugin.Context) error {
	_, err := ctx.ReplyEmbedBuilder(info.Embed)
	return err
}
