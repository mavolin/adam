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

func TestInsufficientBotPermissionsError_Is(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var perms discord.Permissions = 123

		err1 := NewInsufficientBotPermissionsError(perms)
		err2 := NewInsufficientBotPermissionsError(perms)

		assert.True(t, err1.Is(err2))
	})

	t.Run("different types", func(t *testing.T) {
		err1 := NewInsufficientBotPermissionsError(1)
		err2 := New("abc")

		assert.False(t, err1.Is(err2))
	})

	t.Run("different missing permissions", func(t *testing.T) {
		err1 := NewInsufficientBotPermissionsError(123)
		err2 := NewInsufficientBotPermissionsError(456)

		assert.False(t, err1.Is(err2))
	})
}

func TestInsufficientBotPermissionsError_Handle(t *testing.T) {
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
		WithDescription("It seems as if I don't have sufficient permissions to run this command. Please give me the "+
			"following permissions and try again.").
		WithField("Missing Permissions", "• Administrator\n• Video").
		Build()

	m.SendEmbed(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds:    []discord.Embed{embed},
	})

	e := NewInsufficientBotPermissionsError(discord.PermissionAdministrator | discord.PermissionStream)

	err := e.Handle(s, ctx)
	require.NoError(t, err)

	m.Eval()
}
