package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestRestrictionError_Description(t *testing.T) {
	t.Run("string description", func(t *testing.T) {
		expect := "abc"

		e := NewRestrictionError(expect)

		actual, err := e.Description(mock.NoOpLocalizer)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		var term i18n.Term = "abc"

		expect := "def"

		l := mock.
			NewLocalizer(t).
			On(term, expect).
			Build()

		e := NewRestrictionErrorlt(term)

		actual, err := e.Description(l)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestRestrictionError_Handle(t *testing.T) {
	expectDesc := "abc"

	m, s := state.NewMocker(t)
	defer m.Eval()

	ctx := &plugin.Context{
		Message:   discord.Message{ChannelID: 123},
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

	e := NewRestrictionError(expectDesc)

	err := e.Handle(s, ctx)
	require.NoError(t, err)
}
