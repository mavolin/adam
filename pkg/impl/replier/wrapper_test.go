package replier

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func Test_wrappedReplier_Reply(t *testing.T) {
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

	r := WrapState(s)

	data := api.SendMessageData{Content: "abc"}

	expect := discord.Message{
		ID:        012,
		ChannelID: ctx.ChannelID,
		Author:    ctx.Author,
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expect)

	actual, err := r.Reply(ctx, data)
	require.NoError(t, err)
	assert.Equal(t, expect, *actual)
}

func Test_wrappedReplier_ReplyDM(t *testing.T) {
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

		r := WrapState(s)

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
			ID:        012,
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

	r := WrapState(s)

	data := api.EditMessageData{Content: option.NewNullableString("abc")}

	expect := discord.Message{
		ID:        012,
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

		r := WrapState(s)

		data := api.EditMessageData{Content: option.NewNullableString("abc")}

		expect := discord.Message{
			ID:        012,
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

		r := &wrappedReplier{
			s:    s,
			dmID: dmID,
		}

		data := api.EditMessageData{Content: option.NewNullableString("abc")}

		expect := discord.Message{
			ID:        012,
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
