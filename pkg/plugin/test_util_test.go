package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
)

func Test_mockLocalizer_build(t *testing.T) {
	t.Run("expected i18n", func(t *testing.T) {
		t.Run("on", func(t *testing.T) {
			var term i18n.Term = "abc"

			expect := "def"

			l := newMockedLocalizer(t).
				on(term, expect).
				build()

			actual, err := l.LocalizeTerm(term)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("unexpected i18n", func(t *testing.T) {
		tMock := new(testing.T)

		l := newMockedLocalizer(tMock).
			build()

		actual, err := l.LocalizeTerm("unknown_term")
		assert.Empty(tMock, actual)
		assert.Error(t, err)

		assert.True(t, tMock.Failed())
	})
}

func Test_wrappedReplier_Reply(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	r := replierFromState(s, 123, 456)

	data := api.SendMessageData{Content: "abc"}

	expect := discord.Message{
		ID:        12,
		ChannelID: r.channelID,
		Author:    discord.User{ID: r.userID},
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expect)

	actual, err := r.Reply(nil, data)
	require.NoError(t, err)
	assert.Equal(t, expect, *actual)
}

func Test_wrappedReplier_ReplyDM(t *testing.T) {
	t.Run("unknown dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		r := replierFromState(s, 123, 456)

		var dmID discord.ChannelID = 789

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        12,
			ChannelID: dmID,
			Author:    discord.User{ID: r.userID},
			Content:   data.Content,
		}

		m.CreatePrivateChannel(discord.Channel{
			ID:           dmID,
			DMRecipients: []discord.User{{ID: r.userID}},
		})
		m.SendMessageComplex(data, expect)

		actual, err := r.ReplyDM(nil, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})

	t.Run("known dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		r := replierFromState(s, 123, 456)
		r.dmID = 789

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        12,
			ChannelID: r.dmID,
			Author:    discord.User{ID: r.userID},
			Content:   data.Content,
		}

		m.SendMessageComplex(data, expect)

		actual, err := r.ReplyDM(nil, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}

func Test_wrappedReplier_Edit(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	r := replierFromState(s, 123, 456)

	data := api.EditMessageData{Content: option.NewNullableString("abc")}

	expect := discord.Message{
		ID:        12,
		ChannelID: r.channelID,
		Author:    discord.User{ID: r.userID},
		Content:   data.Content.Val,
	}

	m.EditMessageComplex(data, expect)

	actual, err := r.Edit(nil, expect.ID, data)
	require.NoError(t, err)
	assert.Equal(t, expect, *actual)
}

func Test_wrappedReplier_EditDM(t *testing.T) {
	t.Run("unknown dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		r := replierFromState(s, 123, 456)

		var dmID discord.ChannelID = 789

		data := api.EditMessageData{Content: option.NewNullableString("abc")}

		expect := discord.Message{
			ID:        12,
			ChannelID: dmID,
			Author:    discord.User{ID: r.userID},
			Content:   data.Content.Val,
		}

		m.CreatePrivateChannel(discord.Channel{
			ID:           dmID,
			DMRecipients: []discord.User{{ID: r.userID}},
		})
		m.EditMessageComplex(data, expect)

		actual, err := r.EditDM(nil, expect.ID, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})

	t.Run("known dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		r := replierFromState(s, 123, 456)
		r.dmID = 789

		data := api.EditMessageData{Content: option.NewNullableString("abc")}

		expect := discord.Message{
			ID:        12,
			ChannelID: r.dmID,
			Author:    discord.User{ID: r.userID},
			Content:   data.Content.Val,
		}

		m.EditMessageComplex(data, expect)

		actual, err := r.EditDM(nil, expect.ID, data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}
