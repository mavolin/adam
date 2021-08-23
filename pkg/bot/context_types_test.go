package bot

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/state/store"
	"github.com/diamondburned/arikawa/v2/state/store/defaultstore"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// plugin.ErrorHandler
// =====================================================================================

func Test_newCtxErrorHandler(t *testing.T) {
	var called bool

	f := func(error, *state.State, *plugin.Context) { called = true }

	h := newCtxErrorHandler(nil, nil, f)
	h(errors.New("abc"))

	assert.True(t, called, "wrapped error handler was not called")
}

func TestCtxErrorHandler_HandleError(t *testing.T) {
	var actual error

	var h ctxErrorHandler = func(err error) { actual = err }

	expect := errors.New("boom")

	h.HandleError(expect)

	assert.Equal(t, expect, actual)
}

// =============================================================================
// plugin.DiscordDataProvider
// =====================================================================================

func TestDiscordDataProvider_GuildAsync(t *testing.T) {
	t.Run("cached", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Guild{ID: 123}

		s.Cabinet = store.Cabinet{GuildStore: defaultstore.NewGuild()}

		err := s.Cabinet.GuildSet(*expect)
		require.NoError(t, err)

		p := &discordDataProvider{
			s:       s,
			guildID: expect.ID,
		}

		actual, err := p.GuildAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("fetch", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Guild{
			ID:                     123,
			OwnerID:                1,
			RulesChannelID:         1,
			PublicUpdatesChannelID: 1,
		}

		m.Guild(*expect)

		p := &discordDataProvider{
			s:       s,
			guildID: expect.ID,
		}

		actual, err := p.GuildAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestDiscordDataProvider_ChannelAsync(t *testing.T) {
	t.Run("cached", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Channel{ID: 123, GuildID: 456}

		s.Cabinet = store.Cabinet{ChannelStore: defaultstore.NewChannel()}

		err := s.Cabinet.ChannelSet(*expect)
		require.NoError(t, err)

		p := &discordDataProvider{
			s:         s,
			channelID: expect.ID,
		}

		actual, err := p.ChannelAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("fetch", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Channel{
			ID: 123,
		}

		m.Channel(*expect)

		p := &discordDataProvider{
			s:         s,
			channelID: expect.ID,
		}

		actual, err := p.ChannelAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestDiscordDataProvider_MemberAsync(t *testing.T) {
	t.Run("cached", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := &discord.Member{
			User: discord.User{ID: 456},
		}

		s.Cabinet = store.Cabinet{MemberStore: defaultstore.NewMember()}

		err := s.Cabinet.MemberSet(guildID, *expect)
		require.NoError(t, err)

		p := &discordDataProvider{
			s:       s,
			guildID: guildID,
			selfID:  expect.User.ID,
		}

		actual, err := p.SelfAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("fetch", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := &discord.Member{
			User: discord.User{ID: 456},
		}

		m.Member(guildID, *expect)

		p := &discordDataProvider{
			s:       s,
			guildID: guildID,
			selfID:  expect.User.ID,
		}

		actual, err := p.SelfAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}
