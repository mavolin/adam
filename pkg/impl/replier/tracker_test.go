package replier

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestTracker_GuildMessages(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &plugin.Context{
		Message: discord.Message{
			ChannelID: 123,
			Author:    discord.User{ID: 456},
		},
		DiscordDataProvider: mock.DiscordDataProvider{
			ChannelReturn: &discord.Channel{},
			GuildReturn: &discord.Guild{
				Roles: []discord.Role{
					{ID: 789, Permissions: discord.PermissionAdministrator},
				},
			},
			SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{789}},
		},
	}

	r := NewTracker(s)

	data := api.SendMessageData{Content: "abc"}

	expectMessage := discord.Message{
		ID:        012,
		ChannelID: ctx.ChannelID,
		Author:    ctx.Author,
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expectMessage)

	actualMessage, err := r.ReplyMessage(ctx, data)
	require.NoError(t, err)
	assert.Equal(t, expectMessage, *actualMessage)

	expectGuildMessage := []discord.Message{expectMessage}

	actualGuildMessages := r.GuildMessages()
	assert.Equal(t, expectGuildMessage, actualGuildMessages)
}

func TestTracker_DMs(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &plugin.Context{
		Message: discord.Message{Author: discord.User{ID: 123}},
		DiscordDataProvider: mock.DiscordDataProvider{
			ChannelReturn: &discord.Channel{},
			GuildReturn: &discord.Guild{
				Roles: []discord.Role{
					{ID: 456, Permissions: discord.PermissionAdministrator},
				},
			},
			SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{456}},
		},
	}

	r := &Tracker{
		s:    s,
		dmID: 789,
	}

	data := api.SendMessageData{Content: "abc"}

	expectMessage := discord.Message{
		ID:        012,
		ChannelID: r.dmID,
		Author:    ctx.Author,
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expectMessage)

	actualMessage, err := r.ReplyDM(ctx, data)
	require.NoError(t, err)
	assert.Equal(t, expectMessage, *actualMessage)

	expectDMs := []discord.Message{expectMessage}

	actualDMs := r.DMs()
	assert.Equal(t, expectDMs, actualDMs)
}

func TestTracker_ReplyMessage(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &plugin.Context{
		Message: discord.Message{
			ChannelID: 123,
			Author:    discord.User{ID: 456},
		},
		DiscordDataProvider: mock.DiscordDataProvider{
			ChannelReturn: &discord.Channel{},
			GuildReturn: &discord.Guild{
				Roles: []discord.Role{
					{ID: 789, Permissions: discord.PermissionAdministrator},
				},
			},
			SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{789}},
		},
	}

	r := NewTracker(s)

	data := api.SendMessageData{Content: "abc"}

	expect := discord.Message{
		ID:        012,
		ChannelID: ctx.ChannelID,
		Author:    ctx.Author,
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expect)

	actual, err := r.ReplyMessage(ctx, data)
	require.NoError(t, err)
	assert.Equal(t, expect, *actual)
}

func TestTracker_ReplyDM(t *testing.T) {
	t.Run("unknown dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Message: discord.Message{Author: discord.User{ID: 123}},
			DiscordDataProvider: mock.DiscordDataProvider{
				ChannelReturn: &discord.Channel{},
				GuildReturn: &discord.Guild{
					Roles: []discord.Role{
						{ID: 456, Permissions: discord.PermissionAdministrator},
					},
				},
				SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{456}},
			},
		}

		var dmID discord.ChannelID = 789

		r := NewTracker(s)

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        012,
			ChannelID: dmID,
			Author:    ctx.Author,
			Content:   data.Content,
		}

		m.CreatePrivateChannel(discord.Channel{
			ID:           dmID,
			DMRecipients: []discord.User{ctx.Author},
		})
		m.SendMessageComplex(data, expect)

		actual, err := r.ReplyDM(ctx, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})

	t.Run("known dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Message: discord.Message{
				ChannelID: 123,
				Author:    discord.User{ID: 456},
			},
			DiscordDataProvider: mock.DiscordDataProvider{
				ChannelReturn: &discord.Channel{},
				GuildReturn: &discord.Guild{
					Roles: []discord.Role{
						{ID: 789, Permissions: discord.PermissionAdministrator},
					},
				},
				SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{789}},
			},
		}

		r := &Tracker{
			s:    s,
			dmID: ctx.ChannelID,
		}

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        012,
			ChannelID: ctx.ChannelID,
			Author:    ctx.Author,
			Content:   data.Content,
		}

		m.SendMessageComplex(data, expect)

		actual, err := r.ReplyDM(ctx, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}
