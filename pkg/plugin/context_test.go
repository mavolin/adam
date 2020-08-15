package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/mock"
)

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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_Replyl(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    = "abc"
		content = "def"
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
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestContext_Replylt(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    = "abc"
		content = "def"
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)

	m.Eval()
}
