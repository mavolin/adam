package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestUserError_Handle(t *testing.T) {
	t.Run("without embed", func(t *testing.T) {
		expectDesc := "abc"

		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						ChannelID: 123,
					},
				},
			},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 123, 0),
		}

		embed := ErrorEmbed.Clone().
			WithDescription(expectDesc).
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserError(expectDesc)

		e.Handle(s, ctx)
	})

	t.Run("with embed", func(t *testing.T) {
		var (
			expectDesc       = "abc"
			expectFieldName  = "def"
			expectFieldValue = "ghi"
		)

		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						ChannelID: 123,
					},
				},
			},
			Localizer: mock.NoOpLocalizer,
			Replier:   replierFromState(s, 123, 0),
		}

		embed := ErrorEmbed.Clone().
			WithDescription(expectDesc).
			WithField(expectFieldName, expectFieldValue).
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := NewUserError(expectDesc).
			WithField(expectFieldName, expectFieldValue)

		e.Handle(s, ctx)
	})
}
