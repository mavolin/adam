package errors

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// UserInfo is less sever error on the user-side.
// The error will reported to the user via a message containing a detailed
// description of the problem.
// It won't be logged or captured by sentry.
type UserInfo struct {
	// description of the info, either is set
	descString string
	descConfig localization.Config

	fields *embedutil.Builder
}

// NewUserInfo creates a new UserInfo using the passed description.
// The description mustn't be empty for this error to be handled properly.
func NewUserInfo(desc string) *UserInfo {
	return &UserInfo{
		descString: desc,
		fields:     embedutil.NewBuilder(),
	}
}

// NewUserInfol creates a new UserInfo using the message generated from the
// passed localization.Config.
func NewUserInfol(description localization.Config) *UserInfo {
	return &UserInfo{
		descConfig: description,
		fields:     embedutil.NewBuilder(),
	}
}

// NewUserInfolt creates a new UserInfo using the message generated from the
// passed term.
func NewUserInfolt(description localization.Term) *UserInfo {
	return NewUserInfol(description.AsConfig())
}

// WithField adds the passed field to the UserInfo, and returns a pointer to
// the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithField(name, value string) *UserInfo {
	i.fields.WithField(name, value)
	return i
}

// WithFieldl adds the passed field to the UserInfo, and returns a pointer to
// the
// UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithFieldl(name, value localization.Config) *UserInfo {
	i.fields.WithFieldl(name, value)
	return i
}

// WithFieldlt adds the passed field to the UserInfo, and returns a pointer to
// the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithFieldlt(name, value localization.Term) *UserInfo {
	return i.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithField adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedField(name, value string) *UserInfo {
	i.fields.WithInlinedField(name, value)
	return i
}

// WithFieldl adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedFieldl(name, value localization.Config) *UserInfo {
	i.fields.WithInlinedFieldl(name, value)
	return i
}

// WithFieldlt adds the passed inlined field to the UserInfo, and returns a
// pointer to the UserInfo to allow chaining.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (i *UserInfo) WithInlinedFieldlt(name, value localization.Term) *UserInfo {
	return i.WithFieldl(name.AsConfig(), value.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (i *UserInfo) Description(l *localization.Localizer) (string, error) {
	if i.descString != "" {
		return i.descString, nil
	}

	return l.Localize(i.descConfig)
}

// Fields returns the discord.EmbedFields of the UserInfo.
// This can be safely ignored, if only used for UserInfos generated by adam, as
// these will never have fields.
func (i *UserInfo) Fields(l *localization.Localizer) ([]discord.EmbedField, error) {
	embed, err := i.fields.Build(l)
	if err != nil {
		return nil, err
	}

	return embed.Fields, nil
}

func (i *UserInfo) Error() string { return "user info" }

// Handle sends an info embed with the description of the UserInfo.
func (i *UserInfo) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	desc, err := i.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed, err := InfoEmbed.Clone().
		WithDescription(desc).
		Build(ctx.Localizer)
	if err != nil {
		return err
	}

	fields, err := i.Fields(ctx.Localizer)
	if err != nil {
		return err
	}

	embed.Fields = fields

	_, err = ctx.ReplyEmbed(embed)
	return err
}
