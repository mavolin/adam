package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

func TestContext_IsBotOwner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var owner discord.UserID = 123

		ctx := &Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Author: discord.User{
							ID: owner,
						},
					},
				},
			},
			BotOwnerIDs: []discord.UserID{owner},
		}

		assert.True(t, ctx.IsBotOwner())
	})

	t.Run("failure", func(t *testing.T) {
		ctx := &Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Author: discord.User{
							ID: 123,
						},
					},
				},
			},
			BotOwnerIDs: []discord.UserID{465},
		}

		assert.False(t, ctx.IsBotOwner())
	})
}

func TestContext_Reply(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
				},
			},
		},
		s: s,
	}

	expect := &discord.Message{
		ID: 123,
		Author: discord.User{
			ID: 456,
		},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
	}

	m.SendText(*expect)

	actual, err := ctx.Reply(expect.Content)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_Replyl(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    localization.Term = "abc"
		content                   = "def"
	)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
				},
			},
		},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		s: s,
	}

	expect := &discord.Message{
		ID: 123,
		Author: discord.User{
			ID: 456,
		},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.SendText(*expect)

	actual, err := ctx.Replyl(localization.Config{
		Term: term,
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_Replylt(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    localization.Term = "abc"
		content                   = "def"
	)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
				},
			},
		},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		s: s,
	}

	expect := &discord.Message{
		ID: 123,
		Author: discord.User{
			ID: 456,
		},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.SendText(*expect)

	actual, err := ctx.Replylt(term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_ReplyEmbed(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
				},
			},
		},
		s: s,
	}

	expect := &discord.Message{
		ID: 123,
		Author: discord.User{
			ID: 456,
		},
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

	m.Eval()
}

func TestContext_ReplyEmbedBuilder(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
				},
			},
		},
		s: s,
	}

	builder := embedutil.
		NewBuilder().
		WithSimpleTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.MustBuild(nil)
	embed.Type = discord.NormalEmbed

	expect := &discord.Message{
		ID: 123,
		Author: discord.User{
			ID: 456,
		},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
		Embeds:    []discord.Embed{embed},
	}

	m.SendEmbed(*expect)

	actual, err := ctx.ReplyEmbedBuilder(builder)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_ReplyMessage(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
				},
			},
		},
		s: s,
	}

	expect := &discord.Message{
		ID: 123,
		Author: discord.User{
			ID: 456,
		},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
	}

	m.SendText(*expect)

	actual, err := ctx.ReplyMessage(api.SendMessageData{
		Content: expect.Content,
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_ReplyDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					Author: discord.User{
						ID: 123,
					},
				},
			},
		},
		s: s,
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID: 789,
		Author: discord.User{
			ID: 123,
		},
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

	m.Eval()
}

func TestContext_ReplyDMl(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    localization.Term = "abc"
		content                   = "def"
	)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					Author: discord.User{
						ID: 123,
					},
				},
			},
		},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		s: s,
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID: 789,
		Author: discord.User{
			ID: 123,
		},
		ChannelID: channelID,
		Content:   content,
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendText(*expect)

	actual, err := ctx.ReplyDMl(localization.Config{
		Term: term,
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_ReplyDMlt(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    localization.Term = "abc"
		content                   = "def"
	)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					Author: discord.User{
						ID: 123,
					},
				},
			},
		},
		Localizer: newMockedLocalizer(t).
			on(term, content).
			build(),
		s: s,
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID: 789,
		Author: discord.User{
			ID: 123,
		},
		ChannelID: channelID,
		Content:   content,
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendText(*expect)

	actual, err := ctx.ReplyDMlt(term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_ReplyEmbedDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					Author: discord.User{
						ID: 123,
					},
				},
			},
		},
		s: s,
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID: 789,
		Author: discord.User{
			ID: 123,
		},
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

	m.Eval()
}

func TestContext_ReplyEmbedBuilderDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					Author: discord.User{
						ID: 123,
					},
				},
			},
		},
		s: s,
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
		ID: 789,
		Author: discord.User{
			ID: 123,
		},
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

	m.Eval()
}

func TestContext_ReplyMessageDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					Author: discord.User{
						ID: 123,
					},
				},
			},
		},
		s: s,
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID: 789,
		Author: discord.User{
			ID: 123,
		},
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

	m.Eval()
}

func TestContext_DeleteInvoke(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := NewContext(s)
	ctx.MessageCreateEvent = &state.MessageCreateEvent{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: discord.Message{
				ID:        123,
				ChannelID: 456,
			},
		},
	}

	m.DeleteMessage(ctx.ChannelID, ctx.ID)

	err := ctx.DeleteInvoke()
	assert.NoError(t, err)

	m.Eval()
}

func TestContext_HasSelfPermission(t *testing.T) {
	testCases := []struct {
		name    string
		check   discord.Permissions
		dm      bool
		guild   *discord.Guild
		channel *discord.Channel
		self    *discord.Member
		expect  bool
	}{
		{
			name:    "pass dm",
			check:   constant.DMPermissions,
			dm:      true,
			guild:   nil,
			channel: nil,
			self:    nil,
			expect:  true,
		},
		{
			name:    "fail dm",
			check:   constant.DMPermissions + 1,
			dm:      true,
			guild:   nil,
			channel: nil,
			self:    nil,
			expect:  false,
		},
		{
			name:  "pass guild",
			check: discord.PermissionAdministrator,
			dm:    false,
			guild: &discord.Guild{
				OwnerID: 123,
				Roles: []discord.Role{
					{
						ID:          456,
						Permissions: discord.PermissionAdministrator,
					},
				},
			},
			channel: &discord.Channel{},
			self: &discord.Member{
				RoleIDs: []discord.RoleID{456},
			},
			expect: true,
		},
		{
			name:  "fail guild",
			check: discord.PermissionAdministrator,
			dm:    false,
			guild: &discord.Guild{
				OwnerID: 123,
				Roles: []discord.Role{
					{
						ID:          456,
						Permissions: discord.PermissionViewChannel,
					},
				},
			},
			channel: &discord.Channel{},
			self: &discord.Member{
				RoleIDs: []discord.RoleID{456},
			},
			expect: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := &Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: new(gateway.MessageCreateEvent),
				},
				DiscordDataProvider: mockDiscordDataProvider{
					GuildReturn:   c.guild,
					ChannelReturn: c.channel,
					SelfReturn:    c.self,
				},
			}

			if !c.dm {
				ctx.GuildID = 123
			}

			actual, err := ctx.HasSelfPermission(c.check)
			require.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestContext_HasUserPermission(t *testing.T) {
	testCases := []struct {
		name    string
		check   discord.Permissions
		dm      bool
		guild   *discord.Guild
		channel *discord.Channel
		member  *discord.Member
		expect  bool
	}{
		{
			name:    "pass dm",
			check:   constant.DMPermissions,
			dm:      true,
			guild:   nil,
			channel: nil,
			member:  nil,
			expect:  true,
		},
		{
			name:    "fail dm",
			check:   constant.DMPermissions + 1,
			dm:      true,
			guild:   nil,
			channel: nil,
			member:  nil,
			expect:  false,
		},
		{
			name:  "pass guild",
			check: discord.PermissionAdministrator,
			dm:    false,
			guild: &discord.Guild{
				OwnerID: 123,
				Roles: []discord.Role{
					{
						ID:          456,
						Permissions: discord.PermissionAdministrator,
					},
				},
			},
			channel: &discord.Channel{},
			member: &discord.Member{
				RoleIDs: []discord.RoleID{456},
			},
			expect: true,
		},
		{
			name:  "fail guild",
			check: discord.PermissionAdministrator,
			dm:    false,
			guild: &discord.Guild{
				OwnerID: 123,
				Roles: []discord.Role{
					{
						ID:          456,
						Permissions: discord.PermissionViewChannel,
					},
				},
			},
			channel: &discord.Channel{},
			member: &discord.Member{
				RoleIDs: []discord.RoleID{456},
			},
			expect: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := &Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: new(gateway.MessageCreateEvent),
				},
				DiscordDataProvider: mockDiscordDataProvider{
					GuildReturn:   c.guild,
					ChannelReturn: c.channel,
				},
			}

			if !c.dm {
				ctx.GuildID = 123
				ctx.Member = c.member
			}

			actual, err := ctx.HasUserPermission(c.check)
			require.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}
