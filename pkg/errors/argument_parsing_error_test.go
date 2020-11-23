package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestArgumentParsingError_Description(t *testing.T) {
	t.Run("string description", func(t *testing.T) {
		expect := "abc"

		e := NewArgumentParsingError(expect)

		actual, err := e.Description(nil)
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		var term i18n.Term = "abc"

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

func TestArgumentParsingError_Handle(t *testing.T) {
	expectDesc := "abc"

	var channelID discord.ChannelID = 123

	m, s := state.NewMocker(t)

	ctx := &plugin.Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: channelID,
				},
			},
		},
		Localizer: mock.NoOpLocalizer,
		Replier:   replierFromState(s, 123, 0),
	}

	m.SendEmbed(discord.Message{
		ChannelID: channelID,
		Embeds: []discord.Embed{
			ErrorEmbed.Clone().
				WithDescription(expectDesc).
				MustBuild(ctx.Localizer),
		},
	})

	e := NewArgumentParsingError(expectDesc)

	e.Handle(nil, ctx)

	m.Eval()
}
