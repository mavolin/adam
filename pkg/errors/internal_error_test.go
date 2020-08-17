package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/getsentry/sentry-go"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestInternalError_Description(t *testing.T) {
	t.Run("string description", func(t *testing.T) {
		expect := "abc"

		e := WithDescription(New(""), expect).(*InternalError)

		actual := e.Description(mock.NewNoOpLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		var term localization.Term = "abc"

		expect := "def"

		l := mock.
			NewLocalizer().
			On(term, expect).
			Build()

		e := WithDescriptionlt(New(""), term).(*InternalError)

		actual := e.Description(l)
		assert.Equal(t, expect, actual)
	})

	t.Run("invalid description", func(t *testing.T) {
		e := WithDescription(New(""), "").(*InternalError)

		actual := e.Description(mock.NewNoOpLocalizer())
		assert.NotEmpty(t, actual)
	})
}

func TestInternalError_Handle(t *testing.T) {
	expectDesc := "abc"

	m, s := state.NewMocker(t)

	ctx := plugin.NewContext(s)
	ctx.MessageCreateEvent = &state.MessageCreateEvent{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: discord.Message{
				ChannelID: 123,
			},
		},
	}
	ctx.Hub = sentry.NewHub(nil, sentry.NewScope())
	ctx.Localizer = mock.NewNoOpLocalizer()

	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(expectDesc).
		Build()

	m.SendEmbed(discord.Message{
		ChannelID: ctx.ChannelID,
		Embeds: []discord.Embed{
			embed,
		},
	})

	e := WithDescription(New(""), expectDesc).(Handler)

	err := e.Handle(s, ctx)
	assert.NoError(t, err)

	m.Eval()
}
