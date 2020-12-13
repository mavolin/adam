package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/embedutil"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

func TestContext_IsBotOwner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var owner discord.UserID = 123

		ctx := &Context{
			Message:     discord.Message{Author: discord.User{ID: owner}},
			BotOwnerIDs: []discord.UserID{owner},
		}

		assert.True(t, ctx.IsBotOwner())
	})

	t.Run("failure", func(t *testing.T) {
		ctx := &Context{
			Message:     discord.Message{Author: discord.User{ID: 123}},
			BotOwnerIDs: []discord.UserID{465},
		}

		assert.False(t, ctx.IsBotOwner())
	})
}

func TestContext_Reply(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: replierFromState(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
	}

	m.SendText(*expect)

	actual, err := ctx.Reply(expect.Content)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Replyl(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		Replier: replierFromState(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.SendText(*expect)

	actual, err := ctx.Replyl(term.AsConfig())
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Replylt(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		Replier: replierFromState(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.SendText(*expect)

	actual, err := ctx.Replylt(term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbed(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: replierFromState(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
		Embeds: []discord.Embed{
			{
				Type:  discord.NormalEmbed,
				Color: discord.DefaultEmbedColor,
			},
		},
	}

	m.SendEmbed(*expect)

	actual, err := ctx.ReplyEmbed(expect.Embeds[0])
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbedBuilder(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: replierFromState(s, 123, 0),
	}

	builder := embedutil.
		NewBuilder().
		WithSimpleTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.MustBuild(nil)
	embed.Type = discord.NormalEmbed

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
		Embeds:    []discord.Embed{embed},
	}

	m.SendEmbed(*expect)

	actual, err := ctx.ReplyEmbedBuilder(builder)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyMessage(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: replierFromState(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
	}

	m.SendText(*expect)

	actual, err := ctx.ReplyMessage(api.SendMessageData{
		Content: expect.Content,
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyDM(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: replierFromState(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   "abc",
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendText(*expect)

	actual, err := ctx.ReplyDM(expect.Content)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplylDM(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{

		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		Replier: replierFromState(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   content,
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendText(*expect)

	actual, err := ctx.ReplylDM(term.AsConfig())
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyltDM(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		Replier: replierFromState(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   content,
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendText(*expect)

	actual, err := ctx.ReplyltDM(term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbedDM(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: replierFromState(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   "abc",
		Embeds: []discord.Embed{
			{
				Type:  discord.NormalEmbed,
				Color: discord.DefaultEmbedColor,
			},
		},
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendEmbed(*expect)

	actual, err := ctx.ReplyEmbedDM(expect.Embeds[0])
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbedBuilderDM(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: replierFromState(s, 0, 123),
	}

	builder := embedutil.
		NewBuilder().
		WithSimpleTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.MustBuild(nil)
	embed.Type = discord.NormalEmbed

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   "abc",
		Embeds:    []discord.Embed{embed},
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendEmbed(*expect)

	actual, err := ctx.ReplyEmbedBuilderDM(builder)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyMessageDM(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &Context{

		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: replierFromState(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   "abc",
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendText(*expect)

	actual, err := ctx.ReplyMessageDM(api.SendMessageData{
		Content: expect.Content,
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_SelfPermissions(t *testing.T) {
	t.Run("dm", func(t *testing.T) {
		ctx := &Context{
			Message: discord.Message{GuildID: 0},
		}

		actual, err := ctx.SelfPermissions()
		require.NoError(t, err)
		assert.Equal(t, permutil.DMPermissions, actual)
	})

	t.Run("guild", func(t *testing.T) {
		expect := discord.PermissionViewChannel | discord.PermissionAddReactions

		ctx := &Context{
			Message: discord.Message{GuildID: 123},
			DiscordDataProvider: mockDiscordDataProvider{
				GuildReturn: &discord.Guild{
					OwnerID: 123,
					Roles: []discord.Role{
						{
							ID:          456,
							Permissions: expect,
						},
					},
				},
				ChannelReturn: new(discord.Channel),
				SelfReturn:    &discord.Member{RoleIDs: []discord.RoleID{456}},
			},
		}

		actual, err := ctx.SelfPermissions()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestContext_UserPermissions(t *testing.T) {
	t.Run("dm", func(t *testing.T) {
		ctx := &Context{
			Message: discord.Message{GuildID: 0},
		}

		actual, err := ctx.UserPermissions()
		require.NoError(t, err)
		assert.Equal(t, permutil.DMPermissions, actual)
	})

	t.Run("guild", func(t *testing.T) {
		expect := discord.PermissionViewChannel | discord.PermissionAddReactions

		ctx := &Context{
			Message: discord.Message{GuildID: 123},
			Member:  &discord.Member{RoleIDs: []discord.RoleID{456}},
			DiscordDataProvider: mockDiscordDataProvider{
				GuildReturn: &discord.Guild{
					OwnerID: 123,
					Roles: []discord.Role{
						{
							ID:          456,
							Permissions: expect,
						},
					},
				},
				ChannelReturn: new(discord.Channel),
			},
		}

		actual, err := ctx.UserPermissions()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}
