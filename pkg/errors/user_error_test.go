package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestUserError_Handle(t *testing.T) {
	t.Run("without Embed", func(t *testing.T) {
		expectDesc := "abc"

		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 123, 0),
		}

		embed := NewErrorEmbed().
			WithDescription(expectDesc).
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserError(expectDesc)

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
		defer m.Eval()

		ctx := &plugin.Context{
			Message:   discord.Message{ChannelID: 123},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 123, 0),
		}

		embed := NewErrorEmbed().
			WithDescription(expectDesc).
			WithField(expectFieldName, expectFieldValue).
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserError(expectDesc).
			WithField(expectFieldName, expectFieldValue)

		err := e.Handle(s, ctx)
		require.NoError(t, err)
	})
}
