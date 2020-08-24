package errors

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
	"github.com/mavolin/adam/pkg/plugin"
)

func TestUserInfo_Description(t *testing.T) {
	t.Run("string description", func(t *testing.T) {
		expect := "abc"

		e := NewUserInfo(expect)

		actual, err := e.Description(mock.NewNoOpLocalizer())
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		var term localization.Term = "abc"

		expect := "def"

		l := mock.
			NewLocalizer().
			On(term, expect).
			Build()

		e := NewUserInfolt(term)

		actual, err := e.Description(l)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestUserInfo_Handle(t *testing.T) {
	t.Run("without fields", func(t *testing.T) {
		expectDesc := "abc"

		m, s := state.NewMocker(t)

		ctx := plugin.NewContext(s)
		ctx.MessageCreateEvent = &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
					Author: discord.User{
						ID: 456,
					},
				},
			},
		}
		ctx.Localizer = mock.NewNoOpLocalizer()

		embed := newInfoEmbedBuilder(ctx.Localizer).
			WithDescription(expectDesc).
			MustBuild(ctx.Localizer)

		m.SendMessageComplex(api.SendMessageData{
			Embed: &embed,
			AllowedMentions: &api.AllowedMentions{
				Users: []discord.UserID{ctx.Author.ID},
			},
		}, discord.Message{
			ChannelID: ctx.ChannelID,
		})

		e := NewUserInfo(expectDesc)

		err := e.Handle(s, ctx)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("with fields", func(t *testing.T) {
		var (
			expectDesc       = "abc"
			expectFieldName  = "def"
			expectFieldValue = "ghi"
		)

		m, s := state.NewMocker(t)

		ctx := plugin.NewContext(s)
		ctx.MessageCreateEvent = &state.MessageCreateEvent{
			MessageCreateEvent: &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 123,
					Author: discord.User{
						ID: 456,
					},
				},
			},
		}
		ctx.Localizer = mock.NewNoOpLocalizer()

		embed := newInfoEmbedBuilder(ctx.Localizer).
			WithDescription(expectDesc).
			WithField(expectFieldName, expectFieldValue).
			MustBuild(ctx.Localizer)

		m.SendMessageComplex(api.SendMessageData{
			Embed: &embed,
			AllowedMentions: &api.AllowedMentions{
				Users: []discord.UserID{ctx.Author.ID},
			},
		}, discord.Message{
			ChannelID: ctx.ChannelID,
		})

		e := NewUserInfo(expectDesc).
			WithField(expectFieldName, expectFieldValue)

		err := e.Handle(s, ctx)
		require.NoError(t, err)

		m.Eval()
	})
}
