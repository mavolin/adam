package errors

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// UserError is an error on the user side.
	// The error will be reported via a message containing a detailed
	// description of the problem.
	// It won't be logged.
	UserError struct {
		title  *i18n.Config
		desc   *i18n.Config
		fields []userErrorField
	}

	userErrorField struct {
		name *i18n.Config
		val  *i18n.Config
	}
)

// NewUserError creates a new *UserInfo with the passed description.
func NewUserError(description string) *UserError {
	return &UserError{desc: i18n.NewStaticConfig(description)}
}

// NewUserErrorf returns the result of calling NewUserError with
// fmt.Sprinf(description, a...).
func NewUserErrorf(description string, a ...interface{}) *UserError {
	return NewUserError(fmt.Sprintf(description, a...))
}

// NewUserErrorl creates a new *UserInfo using the message generated from the
// passed *i18n.Config as description.
func NewUserErrorl(description *i18n.Config) *UserError {
	return &UserError{desc: description}
}

// WithTitle overwrites the default with the passed title
// (max. 256 characters).
func (e *UserError) WithTitle(title string) *UserError {
	return e.WithTitlel(i18n.NewStaticConfig(title))
}

// WithTitlel overwrites the default with the passed title
// (max. 256 characters).
func (e *UserError) WithTitlel(title *i18n.Config) *UserError {
	e.title = title
	return e
}

// WithField adds a field (name: max. 256 characters, value: max. 1024
// characters) to the embed.
func (e *UserError) WithField(name, value string) *UserError {
	return e.WithFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithFieldl adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (e *UserError) WithFieldl(name, value *i18n.Config) *UserError {
	e.fields = append(e.fields, userErrorField{name: name, val: value})
	return e
}

// Title returns the custom title of the UserInfo.
// If there is no custom title, ("", nil) is returned.
func (e *UserError) Title(l *i18n.Localizer) (string, error) {
	if e.title == nil {
		return "", nil
	}

	return l.Localize(e.title)
}

// Description returns the description of the error.
func (e *UserError) Description(l *i18n.Localizer) (string, error) {
	return l.Localize(e.desc)
}

func (e *UserError) Fields(l *i18n.Localizer) ([]discord.EmbedField, error) {
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

func (e *UserError) Error() string { return "user error" }

// Handle handles the UserInfo.
//
// By default, it creates a NewErrorEmbed and then fills it with the data from
// the UserInfo.
func (e *UserError) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleUserError(s, ctx, e)
}

var HandleUserError = func(s *state.State, ctx *plugin.Context, uerr *UserError) error {
	e := NewErrorEmbed(ctx.Localizer)

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
