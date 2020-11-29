package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
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
	defer m.Eval()

	ctx := &plugin.Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: mock.NewLocalizer(t).
			On(errorTitle.Term, "title").
			On(channelTypeErrorGuild.Term, "guild").
			Build(),
		Replier: replierFromState(s, 123, 0),
	}

	embed := ErrorEmbed.Clone().
		WithDescriptionl(channelTypeErrorGuild).
		MustBuild(ctx.Localizer)

	m.SendEmbed(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds:    []discord.Embed{embed},
	})

	e := NewInvalidChannelTypeError(plugin.GuildChannels)

	e.Handle(s, ctx)
}
