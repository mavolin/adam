package messageutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestReactionWaiter_Await(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						GuildID:   650048604110585858,
						ChannelID: 651147777631584286,
						Author: discord.User{
							ID: 256827968133791744,
						},
					},
				},
			},
			DiscordDataProvider: mock.DiscordDataProvider{
				ChannelReturn: &discord.Channel{},
				ChannelError:  nil,
				GuildReturn: &discord.Guild{
					Roles: []discord.Role{
						{
							ID:          012,
							Permissions: discord.PermissionSendMessages,
						},
					},
				},
				GuildError: nil,
				SelfReturn: &discord.Member{
					RoleIDs: []discord.RoleID{012},
				},
			},
		}

		expect := &TimeoutError{UserID: ctx.Author.ID}

		_, actual := NewReactionWaiter(s, ctx, 123).
			Await(1)
		assert.Equal(t, expect, actual)

		m.Eval()
	})
}
