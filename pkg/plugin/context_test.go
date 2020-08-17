package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/utils/discordutil"
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

	builder := discordutil.
		NewEmbedBuilder().
		WithSimpleTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.Build()
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

func TestContext_ReplyLocalizedEmbedBuilder(t *testing.T) {
	m, s := state.NewMocker(t)

	l := mock.NewLocalizerWithDefault("abc").Build()

	ctx := &Context{
		MessageCreateEvent: &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
				},
			},
		},
		Localizer: l,
		s:         s,
	}

	builder := discordutil.
		NewLocalizedEmbedBuilder().
		WithSimpleTitlet("abc").
		WithDescriptiont("def").
		WithColor(discord.DefaultEmbedColor)

	embed, err := builder.Build(l)
	require.NoError(t, err)

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

	actual, err := ctx.ReplyLocalizedEmbedBuilder(builder)
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
		Localizer: mock.NewLocalizer().
			On(term, content).
			Build(),
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
		Localizer: mock.NewLocalizer().
			On(term, content).
			Build(),
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
