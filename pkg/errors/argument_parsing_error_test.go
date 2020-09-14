package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestArgumentParsingError_WithReason(t *testing.T) {
	reason := "def"

	err1 := NewArgumentParsingError("abc")
	err2 := err1.WithReason(reason)

	assert.NotEqual(t, err1, err2)
	assert.Equal(t, reason, err2.reasonString)
	assert.Equal(t, err1.descString, err2.descString)
}

func TestArgumentParsingError_WithReasonl(t *testing.T) {
	reason := localization.NewTermConfig("def")

	err1 := NewArgumentParsingError("abc")
	err2 := err1.WithReasonl(reason)

	assert.NotEqual(t, err1, err2.reasonConfig)
	assert.Equal(t, err1.descString, err2.descString)
}

func TestArgumentParsingError_WithReasonlt(t *testing.T) {
	reason := localization.NewTermConfig("def")

	err1 := NewArgumentParsingError("abc")
	err2 := err1.WithReasonlt(reason.Term)

	assert.NotEqual(t, err1, err2.reasonConfig)
	assert.Equal(t, err1.descString, err2.descString)
}

func TestArgumentParsingError_Description(t *testing.T) {
	t.Run("string description", func(t *testing.T) {
		expect := "abc"

		e := NewArgumentParsingError(expect)

		actual, err := e.Description(nil)
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		var term localization.Term = "abc"

		expect := "def"

		l := mock.
			NewLocalizer(t).
			On(term, expect).
			Build()

		e := NewArgumentParsingErrorlt(term)

		actual, err := e.Description(l)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestArgumentParsingError_Reason(t *testing.T) {
	t.Run("string reason", func(t *testing.T) {
		expect := "abc"

		e := NewArgumentParsingError("").
			WithReason(expect)

		actual := e.Reason(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized reason", func(t *testing.T) {
		var term localization.Term = "abc"

		expect := "def"

		l := mock.
			NewLocalizer(t).
			On(term, expect).
			Build()

		e := NewArgumentParsingError("").
			WithReasonlt(term)

		actual := e.Reason(l)
		assert.Equal(t, expect, actual)
	})

	t.Run("no reason", func(t *testing.T) {
		e := NewArgumentParsingError("")

		actual := e.Reason(mock.NewNoOpLocalizer())
		assert.Empty(t, actual)
	})
}

func TestArgumentParsingError_Handle(t *testing.T) {
	t.Run("description only", func(t *testing.T) {
		expectDesc := "abc"

		var channelID discord.ChannelID = 123

		m, s := state.NewMocker(t)

		ctx := plugin.NewContext(s)
		ctx.MessageCreateEvent = &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: channelID,
				},
			},
		}
		ctx.Localizer = mock.NewNoOpLocalizer()

		m.SendEmbed(discord.Message{
			ChannelID: channelID,
			Embeds: []discord.Embed{
				ErrorEmbed.Clone().
					WithDescription(expectDesc).
					MustBuild(ctx.Localizer),
			},
		})

		e := NewArgumentParsingError(expectDesc)

		err := e.Handle(nil, ctx)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("with reason", func(t *testing.T) {
		var (
			expectDesc   = "abc"
			expectReason = "def"
		)

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

		embed := ErrorEmbed.Clone().
			WithDescription(expectDesc).
			WithField("Reason", expectReason).
			MustBuild(ctx.Localizer)
		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewArgumentParsingError(expectDesc).
			WithReason(expectReason)

		err := e.Handle(s, ctx)
		require.NoError(t, err)

		m.Eval()
	})
}
