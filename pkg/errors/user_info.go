package errors //nolint:dupl

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// UserInfo is less severe error on the user-side.
// The error will reported to the user via a message containing a detailed
// description of the problem.
// It won't be logged.
type UserInfo struct {
	Embed *embedutil.Builder
}

var _ Error = new(UserInfo)

// NewCustomUserInfo creates a new *UserInfo using a NewInfoEmbed as template.
func NewCustomUserInfo() *UserInfo {
	return &UserInfo{Embed: NewInfoEmbed()}
}

// NewUserInfoFromEmbed creates a new *UserInfo from the passed
// *embedutil.Builder.
func NewUserInfoFromEmbed(e *embedutil.Builder) *UserInfo {
	return &UserInfo{Embed: e}
}

// NewUserInfo creates a new *UserInfo using the passed description and the
// NewInfoEmbed template.
// The description mustn't be empty for this error to be handled properly.
func NewUserInfo(description string) *UserInfo {
	return NewCustomUserInfo().
		WithDescription(description)
}

// NewUserInfol creates a new *UserInfo using the message generated from the
// passed *i18n.Config.
func NewUserInfol(description *i18n.Config) *UserInfo {
	return NewCustomUserInfo().
		WithDescriptionl(description)
}

// NewUserInfolt creates a new *UserInfo using the message generated from the
// passed term.
func NewUserInfolt(description i18n.Term) *UserInfo {
	return NewUserInfol(description.AsConfig())
}

// WithTitle sets the title (max. 256 characters) to the passed title.
func (i *UserInfo) WithTitle(title string) *UserInfo {
	return i.WithTitlel(i18n.NewStaticConfig(title))
}

// WithTitlelt sets the title (max. 256 characters) to the passed title.
func (i *UserInfo) WithTitlelt(title i18n.Term) *UserInfo {
	return i.WithTitlel(title.AsConfig())
}

// WithTitlel sets the title (max. 256 characters) to the passed title.
func (i *UserInfo) WithTitlel(title *i18n.Config) *UserInfo {
	i.Embed.WithTitlel(title)
	return i
}

// WithTitleURL assigns a discord.URL to the title.
func (i *UserInfo) WithTitleURL(url discord.URL) *UserInfo {
	i.Embed.WithTitleURL(url)
	return i
}

// WithDescription sets the description (max. 2048 characters) to the passed
// description.
func (i *UserInfo) WithDescription(description string) *UserInfo {
	return i.WithDescriptionl(i18n.NewStaticConfig(description))
}

// WithDescriptionlt sets the description (max. 2048 characters) to the passed
// description.
func (i *UserInfo) WithDescriptionlt(description i18n.Term) *UserInfo {
	return i.WithDescriptionl(description.AsConfig())
}

// WithDescriptionl sets the description (max. 2048 characters) to the passed
// description.
func (i *UserInfo) WithDescriptionl(description *i18n.Config) *UserInfo {
	i.Embed.WithDescriptionl(description)
	return i
}

// WithTimestamp sets the timestamp to the passed discord.Timestamp.
func (i *UserInfo) WithTimestamp(timestamp discord.Timestamp) *UserInfo {
	i.Embed.WithTimestamp(timestamp)
	return i
}

// WithTimestampNow sets the timestamp to a timestamp of the current time.
func (i *UserInfo) WithTimestampNow() *UserInfo {
	return i.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color to the passed discord.Color.
func (i *UserInfo) WithColor(color discord.Color) *UserInfo {
	i.Embed.WithColor(color)
	return i
}

// WithFooter sets the text of the footer (max. 2048 characters) to the passed
// text.
func (i *UserInfo) WithFooter(text string) *UserInfo {
	return i.WithFooterl(i18n.NewStaticConfig(text))
}

// WithFooterlt sets the text of the footer (max. 2048 characters) to the
// passed text.
func (i *UserInfo) WithFooterlt(text i18n.Term) *UserInfo {
	return i.WithFooterl(text.AsConfig())
}

// WithFooterl sets the text of the footer (max. 2048 characters) to the passed
// text.
func (i *UserInfo) WithFooterl(text *i18n.Config) *UserInfo {
	i.Embed.WithFooterl(text)
	return i
}

// WithFooterIcon sets the icon of the footer to the passed icon url.
func (i *UserInfo) WithFooterIcon(icon discord.URL) *UserInfo {
	i.Embed.WithFooterIcon(icon)
	return i
}

// WithImage sets the image to the passed image url.
func (i *UserInfo) WithImage(image discord.URL) *UserInfo {
	i.Embed.WithImage(image)
	return i
}

// WithThumbnail adds a thumbnail to the error.
func (i *UserInfo) WithThumbnail(thumbnail discord.URL) *UserInfo {
	i.Embed.WithThumbnail(thumbnail)
	return i
}

// WithAuthor sets the author's name (max. 256 characters) to the passed
// name.
func (i *UserInfo) WithAuthor(name string) *UserInfo {
	return i.WithAuthorl(i18n.NewStaticConfig(name))
}

// WithAuthorlt sets the author's name (max. 256 characters) to the passed
// name.
func (i *UserInfo) WithAuthorlt(name i18n.Term) *UserInfo {
	return i.WithAuthorl(name.AsConfig())
}

// WithAuthorl sets the author's name (max. 256 characters) to the passed
// name.
func (i *UserInfo) WithAuthorl(name *i18n.Config) *UserInfo {
	i.Embed.WithAuthorl(name)
	return i
}

// WithAuthorURL assigns the author the passed discord.URL.
func (i *UserInfo) WithAuthorURL(url discord.URL) *UserInfo {
	i.Embed.WithAuthorURL(url)
	return i
}

// WithAuthorIcon sets the icon of the author to the passed icon url.
func (i *UserInfo) WithAuthorIcon(icon discord.URL) *UserInfo {
	i.Embed.WithAuthorIcon(icon)
	return i
}

// WithField adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (i *UserInfo) WithField(name, value string) *UserInfo {
	return i.WithFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithFieldlt adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (i *UserInfo) WithFieldlt(name, value i18n.Term) *UserInfo {
	return i.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithFieldl adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (i *UserInfo) WithFieldl(name, value *i18n.Config) *UserInfo {
	i.Embed.WithFieldl(name, value)
	return i
}

// WithInlinedField adds an inlined field (name: max. 256 characters, value:
// max 1024 characters) to the embed.
func (i *UserInfo) WithInlinedField(name, value string) *UserInfo {
	return i.WithInlinedFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithInlinedFieldlt adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (i *UserInfo) WithInlinedFieldlt(name, value i18n.Term) *UserInfo {
	return i.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithInlinedFieldl adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (i *UserInfo) WithInlinedFieldl(name, value *i18n.Config) *UserInfo {
	i.Embed.WithInlinedFieldl(name, value)
	return i
}

func (i *UserInfo) Error() string { return "user info" }

// Handle sends the info Embed.
func (i *UserInfo) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleUserInfo(i, s, ctx)
}

var HandleUserInfo = func(info *UserInfo, s *state.State, ctx *plugin.Context) error {
	_, err := ctx.ReplyEmbedBuilders(info.Embed)
	return err
}
