package msgawait

import (
	"context"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestReactionWaiter_Await(t *testing.T) {
	t.Parallel()

	t.Run("timeout", func(t *testing.T) {
		t.Parallel()

		_, s := state.NewMocker(t)

		ctx := &plugin.Context{
			Message: discord.Message{
				GuildID:   650048604110585858,
				ChannelID: 651147777631584286,
				Author:    discord.User{ID: 256827968133791744},
			},
			DiscordDataProvider: mock.DiscordDataProvider{
				ChannelReturn: &discord.Channel{},
				ChannelError:  nil,
				GuildReturn: &discord.Guild{
					Roles: []discord.Role{
						{ID: 12, Permissions: discord.PermissionSendMessages},
					},
				},
				GuildError: nil,
				SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{12}},
			},
		}

		expect := &TimeoutError{UserID: ctx.Author.ID, Cause: context.DeadlineExceeded}

		rctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		_, actual := Reaction(s, ctx, 123).
			AwaitContext(rctx)
		assert.Equal(t, expect, actual)
	})
}
