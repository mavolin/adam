package errors

import (
	"fmt"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestWithStack(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithStack(nil)
		assert.Nil(t, err)
	})

	t.Run("handler", func(t *testing.T) {
		cause := NewWithStack("abc")

		err := WithStack(cause)

		assert.True(t, cause == err)
	})

	t.Run("not handler", func(t *testing.T) {
		cause := New("abc")

		err := WithStack(cause)
		unwrapper := err.(interface {
			Unwrap() error
		})

		assert.Equal(t, cause, unwrapper.Unwrap())
	})
}

func TestWrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, "")
		assert.Nil(t, err)
	})

	t.Run("not nil", func(t *testing.T) {
		var (
			cause   = New("abc")
			message = "def"
		)

		err := Wrap(cause, message)

		assert.Equal(t, fmt.Sprintf("%s: %s", message, cause.Error()), err.Error())
	})
}

func TestWrapf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrapf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("not nil", func(t *testing.T) {
		var (
			cause   = New("abc")
			message = "def ghi"
		)

		err := Wrapf(cause, "def %s", "ghi")

		assert.Equal(t, fmt.Sprintf("%s: %s", message, cause.Error()), err.Error())
	})
}

func TestWithDescription(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescription(nil, "")
		assert.Nil(t, err)
	})

	t.Run("internal error", func(t *testing.T) {
		var (
			cause = new(InternalError)
			desc  = "abc"
		)

		err := WithDescription(cause, desc)
		assert.True(t, err == cause)
		assert.Equal(t, desc, cause.descString)
	})

	t.Run("normal error", func(t *testing.T) {
		var (
			cause = New("abc")
			desc  = "def"
		)

		err := WithDescription(cause, desc)
		assert.Equal(t, cause, err.cause)
		assert.Equal(t, desc, err.descString)
	})
}

func TestWithDescriptionf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescriptionf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("internal error", func(t *testing.T) {
		var (
			cause = new(InternalError)
			desc  = "abc def"
		)

		err := WithDescriptionf(cause, "abc %s", "def")
		assert.True(t, err == cause)
		assert.Equal(t, desc, cause.descString)
	})

	t.Run("normal error", func(t *testing.T) {
		var (
			cause = New("abc")
			desc  = "def ghi"
		)

		err := WithDescriptionf(cause, "def %s", "ghi")
		assert.Equal(t, cause, err.cause)
		assert.Equal(t, desc, err.descString)
	})
}

func TestWithDescriptionl(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescriptionl(nil, localization.Config{})
		assert.Nil(t, err)
	})

	t.Run("internal error", func(t *testing.T) {
		var (
			cause = new(InternalError)
			desc  = localization.NewTermConfig("abc")
		)

		err := WithDescriptionl(cause, desc)
		assert.True(t, err == cause)
		assert.Equal(t, desc, cause.descConfig)
	})

	t.Run("normal error", func(t *testing.T) {
		var (
			cause = New("abc")
			desc  = localization.NewTermConfig("def")
		)

		err := WithDescriptionl(cause, desc)
		assert.Equal(t, cause, err.cause)
		assert.Equal(t, desc, err.descConfig)
	})
}

func TestWithDescriptionlt(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescriptionlt(nil, "")
		assert.Nil(t, err)
	})

	t.Run("internal error", func(t *testing.T) {
		var (
			cause                   = new(InternalError)
			desc  localization.Term = "abc"
		)

		err := WithDescriptionlt(cause, desc)
		assert.True(t, err == cause)
		assert.Equal(t, desc.AsConfig(), cause.descConfig)
	})

	t.Run("normal error", func(t *testing.T) {
		var (
			cause                   = New("abc")
			desc  localization.Term = "def"
		)

		err := WithDescriptionlt(cause, desc)
		assert.Equal(t, cause, err.cause)
		assert.Equal(t, desc.AsConfig(), err.descConfig)
	})
}

func TestInternalError_Description(t *testing.T) {
	t.Run("string description", func(t *testing.T) {
		expect := "abc"

		err := WithDescription(New(""), expect)

		actual := err.Description(mock.NewNoOpLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		var term localization.Term = "abc"

		expect := "def"

		l := mock.
			NewLocalizer().
			On(term, expect).
			Build()

		err := WithDescriptionlt(New(""), term)

		actual := err.Description(l)
		assert.Equal(t, expect, actual)
	})

	t.Run("invalid description", func(t *testing.T) {
		err := WithDescription(New(""), "")

		actual := err.Description(mock.NewNoOpLocalizer())
		assert.NotEmpty(t, actual)
	})
}

func TestInternalError_Handle(t *testing.T) {
	expectDesc := "abc"

	m, s := state.NewMocker(t)

	ctx := plugin.NewContext(s)
	ctx.MessageCreateEvent = &state.MessageCreateEvent{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: discord.Message{
				ChannelID: 123,
			},
		},
	}
	ctx.Localizer = mock.NewNoOpLocalizer()

	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(expectDesc).
		MustBuild(ctx.Localizer)

	m.SendEmbed(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds:    []discord.Embed{embed},
	})

	e := WithDescription(New(""), expectDesc)

	err := e.Handle(s, ctx)
	require.NoError(t, err)

	m.Eval()
}
