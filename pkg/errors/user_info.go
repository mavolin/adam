package errors

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// UserInfo is an error on the user side.
	// It is less severe than an UserError, hence Info.
	//
	// The error will be reported via a message containing a detailed
	// description of the problem.
	// It won't be logged.
	UserInfo struct {
		title  *i18n.Config
		desc   *i18n.Config
		fields []userInfoField
	}

	userInfoField struct {
		name *i18n.Config
		val  *i18n.Config
	}
)

// NewUserInfo creates a new *UserInfo with the passed description.
func NewUserInfo(description string) *UserInfo {
	return &UserInfo{desc: i18n.NewStaticConfig(description)}
}

// NewUserInfof returns the result of calling NewUserInfo with
// fmt.Sprinf(description, a...).
func NewUserInfof(description string, a ...interface{}) *UserInfo {
	return NewUserInfo(fmt.Sprintf(description, a...))
}

// NewUserInfol creates a new *UserInfo using the message generated from the
// passed *i18n.Config as description.
func NewUserInfol(description *i18n.Config) *UserInfo {
	return &UserInfo{desc: description}
}

// WithTitle overwrites the default with the passed title
// (max. 256 characters).
func (e *UserInfo) WithTitle(title string) *UserInfo {
	return e.WithTitlel(i18n.NewStaticConfig(title))
}

// WithTitlel overwrites the default with the passed title
// (max. 256 characters).
func (e *UserInfo) WithTitlel(title *i18n.Config) *UserInfo {
	e.title = title
	return e
}

// WithField adds a field (name: max. 256 characters, value: max. 1024
// characters) to the embed.
func (e *UserInfo) WithField(name, value string) *UserInfo {
	return e.WithFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithFieldl adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (e *UserInfo) WithFieldl(name, value *i18n.Config) *UserInfo {
	e.fields = append(e.fields, userInfoField{name: name, val: value})
	return e
}

// Title returns the custom title of the UserInfo.
// If there is no custom title, ("", nil) is returned.
func (e *UserInfo) Title(l *i18n.Localizer) (string, error) {
	if e.title == nil {
		return "", nil
	}

	return l.Localize(e.title)
}

// Description returns the description of the error.
func (e *UserInfo) Description(l *i18n.Localizer) (string, error) {
	return l.Localize(e.desc)
}

func (e *UserInfo) Fields(l *i18n.Localizer) ([]discord.EmbedField, error) {
	lfields := make([]discord.EmbedField, len(e.fields))

	for i, field := range e.fields {
		name, err := l.Localize(field.name)
		if err != nil {
			return nil, err
		}

		val, err := l.Localize(field.val)
		if err != nil {
			return nil, err
		}

		lfields[i] = discord.EmbedField{Name: name, Value: val}
	}

	return lfields, nil
}

func (e *UserInfo) Error() string { return "user info" }

// Handle handles the UserInfo.
//
// By default, it creates a NewInfoEmbed and then fills it with the data from
// the UserInfo.
func (e *UserInfo) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleUserInfo(s, ctx, e)
}

var HandleUserInfo = func(s *state.State, ctx *plugin.Context, uerr *UserInfo) error {
	e := NewInfoEmbed(ctx.Localizer)

	title, err := uerr.Title(ctx.Localizer)
	if err != nil {
		// we have a title to fall back on, so don't return with the error
		ctx.HandleErrorSilently(err)
	} else if title != "" {
		e.Title = title
	}

	e.Description, err = uerr.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	e.Fields, err = uerr.Fields(ctx.Localizer)
	if err != nil {
		return err
	}

	_, err = ctx.ReplyEmbeds(e)
	return err
}
