package plugin

import (
	"errors"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocki18n "github.com/mavolin/adam/internal/mock/i18n"
	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/i18n"
)

// =============================================================================
// ArgumentError
// =====================================================================================

func TestArgumentParsingError_Description(t *testing.T) {
	t.Parallel()

	t.Run("string description", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		e := NewArgumentError(expect)

		actual, err := e.Description(mocki18n.NewLocalizer(t).Build())
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		t.Parallel()

		var term i18n.Term = "abc"

		expect := "def"

		l := mocki18n.NewLocalizer(t).
			On(term, expect).
			Build()

		e := NewArgumentErrorlt(term)

		actual, err := e.Description(l)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestArgumentParsingError_Handle(t *testing.T) {
	t.Parallel()

	expectDesc := "abc"

	var channelID discord.ChannelID = 123

	m, s := state.NewMocker(t)

	ctx := &Context{
		Message:   discord.Message{ChannelID: channelID},
		Localizer: i18n.NewFallbackLocalizer(),
		Replier:   newMockedWrappedReplier(s, 123, 0),
	}

	expectEmbed, err := shared.ErrorEmbed.Clone().
		WithDescription(expectDesc).
		Build(ctx.Localizer)
	require.NoError(t, err)

	m.SendEmbeds(discord.Message{
		ChannelID: channelID,
		Embeds: []discord.Embed{
			expectEmbed,
		},
	})

	e := NewArgumentError(expectDesc)

	err = e.Handle(s, ctx)
	require.NoError(t, err)
}

// =============================================================================
// BotPermissionsError
// =====================================================================================

func TestNewBotPermissionsError(t *testing.T) {
	t.Parallel()

	perms := discord.PermissionViewChannel | discord.PermissionManageEmojisAndStickers

	expect := &BotPermissionsError{Missing: perms}
	actual := NewBotPermissionsError(perms)

	assert.Equal(t, expect, actual)
}

func TestBotPermissionsError_Is(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var perms discord.Permissions = 123

		err1 := NewBotPermissionsError(perms)
		err2 := NewBotPermissionsError(perms)

		assert.True(t, err1.Is(err2))
	})

	t.Run("different types", func(t *testing.T) {
		t.Parallel()

		err1 := NewBotPermissionsError(1)
		err2 := errors.New("abc")

		assert.False(t, err1.Is(err2))
	})

	t.Run("different missing permissions", func(t *testing.T) {
		t.Parallel()

		err1 := NewBotPermissionsError(discord.PermissionStream)
		err2 := NewBotPermissionsError(discord.PermissionUseVAD)

		assert.False(t, err1.Is(err2))
	})
}

func TestBotPermissionsError_Handle(t *testing.T) {
	t.Parallel()

	t.Run("single permission", func(t *testing.T) {
		t.Parallel()

		m, s := state.NewMocker(t)

		ctx := &Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: i18n.NewFallbackLocalizer(),
			Replier:   newMockedWrappedReplier(s, 123, 0),
		}

		e := NewBotPermissionsError(discord.PermissionStream)

		expectDesc := ctx.MustLocalize(botPermissionsDescSingle.
			WithPlaceholders(botPermissionsDescSinglePlaceholders{
				MissingPermission: "Video",
			}))

		expectEmbed, err := shared.ErrorEmbed.Clone().
			WithDescription(expectDesc).
			Build(ctx.Localizer)
		require.NoError(t, err)

		m.SendEmbeds(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{expectEmbed},
		})

		err = e.Handle(s, ctx)
		require.NoError(t, err)
	})

	t.Run("multiple permissions", func(t *testing.T) {
		t.Parallel()

		m, s := state.NewMocker(t)

		ctx := &Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: i18n.NewFallbackLocalizer(),
			Replier:   newMockedWrappedReplier(s, 123, 0),
		}

		expectEmbed, err := shared.ErrorEmbed.Clone().
			WithDescriptionl(botPermissionsDescMulti).
			WithField("Missing Permissions", "• Video\n• View Audit Log").
			Build(ctx.Localizer)
		require.NoError(t, err)

		m.SendEmbeds(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{expectEmbed},
		})

		e := NewBotPermissionsError(discord.PermissionViewAuditLog | discord.PermissionStream)

		err = e.Handle(s, ctx)
		require.NoError(t, err)
	})
}

// =============================================================================
// ChannelTypeError
// =====================================================================================

func TestChannelTypeError_Is(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		types := GuildChannels

		err1 := NewChannelTypeError(types)
		err2 := NewChannelTypeError(types)

		assert.True(t, err1.Is(err2))
	})

	t.Run("different types", func(t *testing.T) {
		t.Parallel()

		err1 := NewChannelTypeError(DirectMessages)
		err2 := errors.New("abc")

		assert.False(t, err1.Is(err2))
	})

	t.Run("different missing permissions", func(t *testing.T) {
		t.Parallel()

		err1 := NewChannelTypeError(DirectMessages)
		err2 := NewChannelTypeError(GuildChannels)

		assert.False(t, err1.Is(err2))
	})
}

func TestChannelTypeError_Handle(t *testing.T) {
	t.Parallel()

	m, s := state.NewMocker(t)

	ctx := &Context{
		Message: discord.Message{ChannelID: 123},
		Localizer: mocki18n.NewLocalizer(t).
			On("error.title", "title").
			On(channelTypeErrorGuild.Term, "guild").
			Build(),
		Replier: newMockedWrappedReplier(s, 123, 0),
	}

	expectEmbed, err := shared.ErrorEmbed.Clone().
		WithDescriptionl(channelTypeErrorGuild).
		Build(ctx.Localizer)
	require.NoError(t, err)

	m.SendEmbeds(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds:    []discord.Embed{expectEmbed},
	})

	e := NewChannelTypeError(GuildChannels)

	err = e.Handle(s, ctx)
	require.NoError(t, err)
}

// =============================================================================
// RestrictionError
// =====================================================================================

func TestRestrictionError_Description(t *testing.T) {
	t.Parallel()

	t.Run("string description", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		e := NewRestrictionError(expect)

		actual, err := e.Description(i18n.NewFallbackLocalizer())
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		t.Parallel()

		var term i18n.Term = "abc"

		expect := "def"

		l := mocki18n.NewLocalizer(t).
			On(term, expect).
			Build()

		e := NewRestrictionErrorlt(term)

		actual, err := e.Description(l)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestRestrictionError_Handle(t *testing.T) {
	t.Parallel()

	expectDesc := "abc"

	m, s := state.NewMocker(t)

	ctx := &Context{
		Message:   discord.Message{ChannelID: 123},
		Localizer: i18n.NewFallbackLocalizer(),
		Replier:   newMockedWrappedReplier(s, 123, 0),
	}

	expectEmbed, err := shared.ErrorEmbed.Clone().
		WithDescription(expectDesc).
		Build(ctx.Localizer)
	require.NoError(t, err)

	m.SendEmbeds(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds:    []discord.Embed{expectEmbed},
	})

	e := NewRestrictionError(expectDesc)

	err = e.Handle(s, ctx)
	require.NoError(t, err)
}

// =============================================================================
// ThrottlingError
// =====================================================================================

func TestThrottlingError_Description(t *testing.T) {
	t.Parallel()

	t.Run("string description", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		e := NewThrottlingError(expect)

		actual, err := e.Description(i18n.NewFallbackLocalizer())
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		t.Parallel()

		var term i18n.Term = "abc"

		expect := "def"

		l := mocki18n.NewLocalizer(t).
			On(term, expect).
			Build()

		e := NewThrottlingErrorlt(term)

		actual, err := e.Description(l)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestThrottlingError_Handle(t *testing.T) {
	t.Parallel()

	expectDesc := "abc"

	m, s := state.NewMocker(t)

	ctx := &Context{
		Message:   discord.Message{ChannelID: 123},
		Localizer: i18n.NewFallbackLocalizer(),
		Replier:   newMockedWrappedReplier(s, 123, 0),
	}

	expectEmbed, err := shared.InfoEmbed.Clone().
		WithDescription(expectDesc).
		Build(ctx.Localizer)
	require.NoError(t, err)

	m.SendEmbeds(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds:    []discord.Embed{expectEmbed},
	})

	e := NewThrottlingError(expectDesc)

	err = e.Handle(s, ctx)
	require.NoError(t, err)
}
