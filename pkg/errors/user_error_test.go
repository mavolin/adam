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

func TestUserError_Handle(t *testing.T) {
	t.Parallel()

	t.Run("without description", func(t *testing.T) {
		t.Parallel()

		expectDesc := "abc"

		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: i18n.NewFallbackLocalizer(),
			Replier:   mockplugin.NewWrappedReplier(s, 123, 0),
		}

		embed := NewErrorEmbed().
			WithDescription(expectDesc).
			MustBuild(ctx.Localizer)

		m.SendEmbeds(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserError(expectDesc)

		err := e.Handle(s, ctx)
		require.NoError(t, err)
	})

	t.Run("with description", func(t *testing.T) {
		t.Parallel()

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

		embed := NewErrorEmbed().
			WithDescription(expectDesc).
			WithField(expectFieldName, expectFieldValue).
			MustBuild(ctx.Localizer)

		m.SendEmbeds(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserError(expectDesc).
			WithField(expectFieldName, expectFieldValue)

		err := e.Handle(s, ctx)
		require.NoError(t, err)
	})
}
