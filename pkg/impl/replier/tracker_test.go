package replier

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracker_GuildMessages(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	var channelID discord.ChannelID = 123

	r := NewTracker(s, 0, channelID)

	data := api.SendMessageData{Content: "abc"}

	expectMessage := discord.Message{
		ID:        456,
		ChannelID: channelID,
		Author:    discord.User{ID: 789},
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expectMessage)

	actualMessage, err := r.ReplyMessage(data)
	require.NoError(t, err)
	assert.Equal(t, expectMessage, *actualMessage)

	expectGuildMessage := []discord.Message{expectMessage}

	actualGuildMessages := r.GuildMessages()
	assert.Equal(t, expectGuildMessage, actualGuildMessages)
}

func TestTracker_DMs(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	var channelID discord.ChannelID = 123

	r := &Tracker{
		s:    s,
		dmID: channelID,
	}

	data := api.SendMessageData{Content: "abc"}

	expectMessage := discord.Message{
		ID:        456,
		ChannelID: channelID,
		Author:    discord.User{ID: 789},
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expectMessage)

	actualMessage, err := r.ReplyDM(data)
	require.NoError(t, err)
	assert.Equal(t, expectMessage, *actualMessage)

	expectDMs := []discord.Message{expectMessage}

	actualDMs := r.DMs()
	assert.Equal(t, expectDMs, actualDMs)
}

func TestTracker_ReplyMessage(t *testing.T) {
	m, s := state.NewMocker(t)
	defer m.Eval()

	var channelID discord.ChannelID = 123

	r := NewTracker(s, 0, channelID)

	data := api.SendMessageData{Content: "abc"}

	expect := discord.Message{
		ID:        456,
		ChannelID: channelID,
		Author:    discord.User{ID: 789},
		Content:   data.Content,
	}

	m.SendMessageComplex(data, expect)

	actual, err := r.ReplyMessage(data)
	require.NoError(t, err)
	assert.Equal(t, expect, *actual)
}

func TestTracker_ReplyDM(t *testing.T) {
	t.Run("unknown dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		var (
			channelID discord.ChannelID = 123
			userID    discord.UserID    = 456
		)

		r := NewTracker(s, userID, 0)

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        789,
			ChannelID: channelID,
			Author:    discord.User{ID: userID},
			Content:   data.Content,
		}

		m.CreatePrivateChannel(discord.Channel{
			ID:           channelID,
			DMRecipients: []discord.User{{ID: userID}},
		})
		m.SendMessageComplex(data, expect)

		actual, err := r.ReplyDM(data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})

	t.Run("known dm id", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		var channelID discord.ChannelID = 123

		r := &Tracker{
			s:    s,
			dmID: channelID,
		}

		data := api.SendMessageData{Content: "abc"}

		expect := discord.Message{
			ID:        456,
			ChannelID: channelID,
			Author:    discord.User{ID: 789},
			Content:   data.Content,
		}

		m.SendMessageComplex(data, expect)

		actual, err := r.ReplyDM(data)
		require.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}
