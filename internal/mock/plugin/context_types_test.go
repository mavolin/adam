package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WrappedReplier_Reply(t *testing.T) {
	t.Parallel()

	m, s := state.NewMocker(t)

	r := NewWrappedReplier(s, 123, 456)

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
	t.Parallel()

	t.Run("unknown dm id", func(t *testing.T) {
		t.Parallel()

		m, s := state.NewMocker(t)

		r := NewWrappedReplier(s, 123, 456)

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
		t.Parallel()

		m, s := state.NewMocker(t)

		r := NewWrappedReplier(s, 123, 456)
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
	t.Parallel()

	m, s := state.NewMocker(t)

	r := NewWrappedReplier(s, 123, 456)

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
	t.Parallel()

	t.Run("unknown dm id", func(t *testing.T) {
		t.Parallel()

		m, s := state.NewMocker(t)

		r := NewWrappedReplier(s, 123, 456)

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
		t.Parallel()

		m, s := state.NewMocker(t)

		r := NewWrappedReplier(s, 123, 456)
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
