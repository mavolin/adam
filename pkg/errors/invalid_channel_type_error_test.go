package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestInvalidChannelError_Is(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		types := plugin.GuildChannels

		err1 := NewInvalidChannelTypeError(types)
		err2 := NewInvalidChannelTypeError(types)

		assert.True(t, err1.Is(err2))
	})

	t.Run("different types", func(t *testing.T) {
		err1 := NewInvalidChannelTypeError(plugin.DirectMessages)
		err2 := New("abc")

		assert.False(t, err1.Is(err2))
	})

	t.Run("different missing permissions", func(t *testing.T) {
		err1 := NewInvalidChannelTypeError(plugin.DirectMessages)
		err2 := NewInvalidChannelTypeError(plugin.GuildChannels)

		assert.False(t, err1.Is(err2))
	})
}

func TestInvalidChannelTypeError_Handle(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := plugin.NewContext(s)
	ctx.MessageCreateEvent = &state.MessageCreateEvent{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: discord.Message{
				ChannelID: 123,
			},
		},
	}
	ctx.Localizer = mock.NewLocalizer(t).
		On(errorTitle.Term, "title").
		On(channelTypeErrorGuild.Term, "guild").
		Build()

	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescriptionl(channelTypeErrorGuild).
		MustBuild(ctx.Localizer)

	m.SendEmbed(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds:    []discord.Embed{embed},
	})

	e := NewInvalidChannelTypeError(plugin.GuildChannels)

	err := e.Handle(s, ctx)
	require.NoError(t, err)

	m.Eval()
}
