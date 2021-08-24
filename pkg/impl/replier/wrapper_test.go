package replier

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func Test_wrappedReplier_Reply(t *testing.T) {
	testCases := []struct {
		name        string
		ctx         *plugin.Context
		inlineReply bool
		inData      api.SendMessageData

		expectData    api.SendMessageData
		expectMessage discord.Message
	}{
		{
			name: "normal reply",
			ctx: &plugin.Context{
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
			},
			inlineReply: false,
			inData:      api.SendMessageData{Content: "abc"},
			expectData:  api.SendMessageData{Content: "abc"},
			expectMessage: discord.Message{
				ID:        12,
				ChannelID: 123,
				Author:    discord.User{ID: 456},
				Content:   "abc",
			},
		},
		{
			name: "inline reply",
			ctx: &plugin.Context{
				Message: discord.Message{
					ID:        123,
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{},
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{ID: 12, Permissions: discord.PermissionAdministrator},
						},
					},
					SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{12}},
				},
			},
			inlineReply: true,
			inData:      api.SendMessageData{Content: "abc"},
			expectData: api.SendMessageData{
				Content:   "abc",
				Reference: &discord.MessageReference{MessageID: 123},
			},
			expectMessage: discord.Message{
				ID:        345,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
				Reference: &discord.MessageReference{MessageID: 123},
			},
		},
		{
			name: "blocked inline reply",
			ctx: &plugin.Context{
				Message: discord.Message{
					ID:        123,
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{},
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{ID: 12, Permissions: discord.PermissionAdministrator},
						},
					},
					SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{12}},
				},
			},
			inlineReply: true,
			inData: api.SendMessageData{
				Content:   "abc",
				Reference: new(discord.MessageReference),
			},
			expectData: api.SendMessageData{
				Content:   "abc",
				Reference: new(discord.MessageReference),
			},
			expectMessage: discord.Message{
				ID:        345,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
				Reference: &discord.MessageReference{MessageID: 123},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := state.NewMocker(t)

			r := WrapState(s, c.inlineReply)

			m.SendMessageComplex(c.expectData, c.expectMessage)

			actual, err := r.Reply(c.ctx, c.inData)
			require.NoError(t, err)
			assert.Equal(t, c.expectMessage, *actual)
		})
	}
}

func Test_wrappedReplier_ReplyDM(t *testing.T) {
	t.Run("unknown dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)

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

		r := WrapState(s, false)

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        12,
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

		r := &wrappedReplier{
			s:    s,
			dmID: dmID,
		}

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        12,
			ChannelID: dmID,
			Author:    ctx.Author,
			Content:   data.Content,
		}

		m.SendMessageComplex(data, expect)

		actual, err := r.ReplyDM(ctx, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}

func Test_wrappedReplier_Edit(t *testing.T) {
	m, s := state.NewMocker(t)

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

	r := WrapState(s, false)

	data := api.EditMessageData{Content: option.NewNullableString("abc")}

	expect := discord.Message{
		ID:        12,
		ChannelID: ctx.ChannelID,
		Author:    ctx.Author,
		Content:   data.Content.Val,
	}

	m.EditMessageComplex(data, expect)

	actual, err := r.Edit(ctx, expect.ID, data)
	require.NoError(t, err)
	assert.Equal(t, expect, *actual)
}

func Test_wrappedReplier_EditDM(t *testing.T) {
	t.Run("unknown dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)

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

		r := WrapState(s, false)

		data := api.EditMessageData{Content: option.NewNullableString("abc")}

		expect := discord.Message{
			ID:        12,
			ChannelID: dmID,
			Author:    ctx.Author,
			Content:   data.Content.Val,
		}

		m.CreatePrivateChannel(discord.Channel{
			ID:           dmID,
			DMRecipients: []discord.User{ctx.Author},
		})
		m.EditMessageComplex(data, expect)

		actual, err := r.EditDM(ctx, expect.ID, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})

	t.Run("known dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)

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

		r := &wrappedReplier{
			s:    s,
			dmID: dmID,
		}

		data := api.EditMessageData{Content: option.NewNullableString("abc")}

		expect := discord.Message{
			ID:        12,
			ChannelID: dmID,
			Author:    ctx.Author,
			Content:   data.Content.Val,
		}

		m.EditMessageComplex(data, expect)

		actual, err := r.EditDM(ctx, expect.ID, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}
