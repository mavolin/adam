package plugin

import (
	"fmt"
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocki18n "github.com/mavolin/adam/internal/mock/i18n"
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

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   fmt.Sprint("abc", "def"),
	}

	m.SendMessage(*expect)

	actual, err := ctx.Reply("abc", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Replyf(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   fmt.Sprintf("abc %s", "def"),
	}

	m.SendMessage(*expect)

	actual, err := ctx.Replyf("abc %s", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Replyl(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.SendMessage(*expect)

	actual, err := ctx.Replyl(term.AsConfig())
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Replylt(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.SendMessage(*expect)

	actual, err := ctx.Replylt(term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbed(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
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

	m.SendEmbeds(*expect)

	actual, err := ctx.ReplyEmbeds(expect.Embeds[0])
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbedBuilder(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message:   discord.Message{ChannelID: 123},
		Replier:   newMockedWrappedReplier(s, 123, 0),
		Localizer: mocki18n.NewLocalizer(t).Build(),
	}

	builder := embedutil.
		NewBuilder().
		WithTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.MustBuild(ctx.Localizer)
	embed.Type = discord.NormalEmbed

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
		Embeds:    []discord.Embed{embed},
	}

	m.SendEmbeds(*expect)

	actual, err := ctx.ReplyEmbedBuilders(builder)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyMessage(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
	}

	m.SendMessage(*expect)

	actual, err := ctx.ReplyMessage(api.SendMessageData{
		Content: expect.Content,
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   fmt.Sprint("abc", "def"),
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendMessage(*expect)

	actual, err := ctx.ReplyDM("abc", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyfDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   fmt.Sprintf("abc %s", "def"),
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.SendMessage(*expect)

	actual, err := ctx.ReplyfDM("abc %s", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplylDM(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{

		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.SendMessage(*expect)

	actual, err := ctx.ReplylDM(term.AsConfig())
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyltDM(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.SendMessage(*expect)

	actual, err := ctx.ReplyltDM(term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbedDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.SendEmbeds(*expect)

	actual, err := ctx.ReplyEmbedsDM(expect.Embeds[0])
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyEmbedBuilderDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier:   newMockedWrappedReplier(s, 0, 123),
		Localizer: mocki18n.NewLocalizer(t).Build(),
	}

	builder := embedutil.
		NewBuilder().
		WithTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.MustBuild(ctx.Localizer)
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
	m.SendEmbeds(*expect)

	actual, err := ctx.ReplyEmbedBuildersDM(builder)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_ReplyMessageDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{

		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.SendMessage(*expect)

	actual, err := ctx.ReplyMessageDM(api.SendMessageData{
		Content: expect.Content,
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Edit(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   fmt.Sprint("abc", "def"),
	}

	m.EditText(*expect)

	actual, err := ctx.Edit(expect.ID, "abc", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Editf(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   fmt.Sprintf("abc %s", "def"),
	}

	m.EditText(*expect)

	actual, err := ctx.Editf(expect.ID, "abc %s", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Editl(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.EditText(*expect)

	actual, err := ctx.Editl(expect.ID, term.AsConfig())
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_Editlt(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   content,
	}

	m.EditText(*expect)

	actual, err := ctx.Editlt(expect.ID, term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditEmbed(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
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

	m.EditEmbeds(*expect)

	actual, err := ctx.EditEmbeds(expect.ID, expect.Embeds[0])
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditEmbedBuilder(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message:   discord.Message{ChannelID: 123},
		Replier:   newMockedWrappedReplier(s, 123, 0),
		Localizer: mocki18n.NewLocalizer(t).Build(),
	}

	builder := embedutil.
		NewBuilder().
		WithTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.MustBuild(ctx.Localizer)
	embed.Type = discord.NormalEmbed

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
		Embeds:    []discord.Embed{embed},
	}

	m.EditEmbeds(*expect)

	actual, err := ctx.EditEmbedBuilders(expect.ID, builder)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditMessage(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expect := &discord.Message{
		ID:        123,
		Author:    discord.User{ID: 456},
		ChannelID: ctx.ChannelID,
		Content:   "abc",
	}

	m.EditText(*expect)

	actual, err := ctx.EditMessage(expect.ID, api.EditMessageData{
		Content: option.NewNullableString(expect.Content),
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   fmt.Sprint("abc", "def"),
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.EditText(*expect)

	actual, err := ctx.EditDM(expect.ID, "abc", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditfDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
	}

	var channelID discord.ChannelID = 456

	expect := &discord.Message{
		ID:        789,
		Author:    discord.User{ID: 123},
		ChannelID: channelID,
		Content:   fmt.Sprintf("abc %s", "def"),
	}

	m.CreatePrivateChannel(discord.Channel{
		ID:           channelID,
		DMRecipients: []discord.User{{ID: ctx.Author.ID}},
	})
	m.EditText(*expect)

	actual, err := ctx.EditfDM(expect.ID, "abc %s", "def")
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditlDM(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{

		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.EditText(*expect)

	actual, err := ctx.EditlDM(expect.ID, term.AsConfig())
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditltDM(t *testing.T) {
	m, s := state.NewMocker(t)

	var (
		term    i18n.Term = "abc"
		content           = "def"
	)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Localizer: mocki18n.NewLocalizer(t).
			On(term, content).
			Build(),
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.EditText(*expect)

	actual, err := ctx.EditltDM(expect.ID, term)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditEmbedDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.EditEmbeds(*expect)

	actual, err := ctx.EditEmbedsDM(expect.ID, expect.Embeds[0])
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditEmbedBuilderDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier:   newMockedWrappedReplier(s, 0, 123),
		Localizer: mocki18n.NewLocalizer(t).Build(),
	}

	builder := embedutil.
		NewBuilder().
		WithTitle("abc").
		WithDescription("def").
		WithColor(discord.DefaultEmbedColor)

	embed := builder.MustBuild(ctx.Localizer)
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
	m.EditEmbeds(*expect)

	actual, err := ctx.EditEmbedBuildersDM(expect.ID, builder)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_EditMessageDM(t *testing.T) {
	m, s := state.NewMocker(t)

	ctx := &Context{

		Message: discord.Message{
			GuildID: 1,
			Author:  discord.User{ID: 123},
		},
		Replier: newMockedWrappedReplier(s, 0, 123),
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
	m.EditText(*expect)

	actual, err := ctx.EditMessageDM(expect.ID, api.EditMessageData{
		Content: option.NewNullableString(expect.Content),
	})
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestContext_SelfPermissions(t *testing.T) {
	t.Run("dm", func(t *testing.T) {
		ctx := &Context{Message: discord.Message{GuildID: 0}}

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
					Roles:   []discord.Role{{ID: 456, Permissions: expect}},
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
					Roles:   []discord.Role{{ID: 456, Permissions: expect}},
				},
				ChannelReturn: new(discord.Channel),
			},
		}

		actual, err := ctx.UserPermissions()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}
