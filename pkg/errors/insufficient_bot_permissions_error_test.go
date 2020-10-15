package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestNewInsufficientBotPermissionsError(t *testing.T) {
	t.Run("regular permissions", func(t *testing.T) {
		perms := discord.PermissionViewChannel | discord.PermissionManageEmojis

		expect := &InsufficientBotPermissionsError{
			MissingPermissions: perms,
		}

		actual := NewInsufficientBotPermissionsError(perms)

		assert.Equal(t, expect, actual)
	})

	t.Run("regular permissions", func(t *testing.T) {
		expect := &InsufficientBotPermissionsError{
			MissingPermissions: discord.PermissionAdministrator,
		}

		actual := NewInsufficientBotPermissionsError(discord.PermissionViewChannel | discord.PermissionAdministrator)

		assert.Equal(t, expect, actual)
	})
}

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
		err1 := NewInsufficientBotPermissionsError(discord.PermissionStream)
		err2 := NewInsufficientBotPermissionsError(discord.PermissionUseVAD)

		assert.False(t, err1.Is(err2))
	})
}

func TestInsufficientBotPermissionsError_Handle(t *testing.T) {
	t.Run("single permission", func(t *testing.T) {
		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						ChannelID: 123,
					},
				},
			},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 0),
		}

		embed := ErrorEmbed.Clone().
			WithDescription("It seems as if I don't have sufficient permissions to run this command. Please give me" +
				` the "Video" permission and try again.`).
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewInsufficientBotPermissionsError(discord.PermissionStream)

		err := e.Handle(s, ctx)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("multiple permissions", func(t *testing.T) {
		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						ChannelID: 123,
					},
				},
			},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 0),
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

		e := NewInsufficientBotPermissionsError(discord.PermissionViewAuditLog | discord.PermissionStream)

		err := e.Handle(s, ctx)
		require.NoError(t, err)

		m.Eval()
	})
}
