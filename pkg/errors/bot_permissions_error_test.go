package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestNewInsufficientBotPermissionsError(t *testing.T) {
	perms := discord.PermissionViewChannel | discord.PermissionManageEmojis

	expect := &BotPermissionsError{Missing: perms}
	actual := NewBotPermissionsError(perms)

	assert.Equal(t, expect, actual)
}

func TestInsufficientBotPermissionsError_Is(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var perms discord.Permissions = 123

		err1 := NewBotPermissionsError(perms)
		err2 := NewBotPermissionsError(perms)

		assert.True(t, err1.Is(err2))
	})

	t.Run("different types", func(t *testing.T) {
		err1 := NewBotPermissionsError(1)
		err2 := New("abc")

		assert.False(t, err1.Is(err2))
	})

	t.Run("different missing permissions", func(t *testing.T) {
		err1 := NewBotPermissionsError(discord.PermissionStream)
		err2 := NewBotPermissionsError(discord.PermissionUseVAD)

		assert.False(t, err1.Is(err2))
	})
}

func TestInsufficientBotPermissionsError_Handle(t *testing.T) {
	t.Run("single permission", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 123, 0),
		}

		embed := ErrorEmbed.Clone().
			WithDescription("It seems as if I don't have sufficient permissions to run this command. Please give me" +
				` the "Video" permission and try again.`).
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewBotPermissionsError(discord.PermissionStream)

		err := e.Handle(s, ctx)
		require.NoError(t, err)
	})

	t.Run("multiple permissions", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 123, 0),
		}

		embed := ErrorEmbed.Clone().
			WithDescription("It seems as if I don't have sufficient permissions to run this command. Please give me the "+
				"following permissions and try again:").
			WithField("Missing Permissions", "• Video\n• View Audit Log").
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewBotPermissionsError(discord.PermissionViewAuditLog | discord.PermissionStream)

		err := e.Handle(s, ctx)
		require.NoError(t, err)
	})
}