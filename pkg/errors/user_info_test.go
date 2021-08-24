package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/require"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestUserInfo_Handle(t *testing.T) {
	t.Run("without Embed", func(t *testing.T) {
		expectDesc := "abc"

		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: i18n.NewFallbackLocalizer(),
			Replier:   mockplugin.NewWrappedReplier(s, 123, 0),
		}

		embed := NewInfoEmbed().
			WithDescription(expectDesc).
			MustBuild(ctx.Localizer)

		m.SendEmbeds(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserInfo(expectDesc)

		err := e.Handle(s, ctx)
		require.NoError(t, err)
	})

	t.Run("with Embed", func(t *testing.T) {
		var (
			expectDesc       = "abc"
			expectFieldName  = "def"
			expectFieldValue = "ghi"
		)

		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: i18n.NewFallbackLocalizer(),
			Replier:   mockplugin.NewWrappedReplier(s, 123, 0),
		}

		embed := NewInfoEmbed().
			WithDescription(expectDesc).
			WithField(expectFieldName, expectFieldValue).
			MustBuild(ctx.Localizer)

		m.SendEmbeds(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserInfo(expectDesc).
			WithField(expectFieldName, expectFieldValue)

		err := e.Handle(s, ctx)
		require.NoError(t, err)
	})
}
