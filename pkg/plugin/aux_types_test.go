package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelType_Has(t *testing.T) {
	successCases := []struct {
		name        string
		channelType ChannelTypes
		target      discord.ChannelType
		expect      bool
	}{
		{
			name:        "all",
			channelType: AllChannels,
			target:      discord.GuildText,
		},
		{
			name:        "GuildTextChannels",
			channelType: GuildTextChannels,
			target:      discord.GuildText,
		},
		{
			name:        "DirectMessages",
			channelType: DirectMessages,
			target:      discord.DirectMessage,
		},
		{
			name:        "GuildNewsChannels",
			channelType: GuildNewsChannels,
			target:      discord.GuildNews,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				has := c.channelType.Has(c.target)
				assert.True(t, has)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		has := GuildTextChannels.Has(discord.DirectMessage)
		assert.False(t, has)
	})
}

func TestChannelTypes_Check(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     *Context
		allowed ChannelTypes
		expect  bool
	}{
		{
			name:    "pass guild channels",
			ctx:     &Context{Message: discord.Message{GuildID: 123}},
			allowed: GuildChannels,
			expect:  true,
		},
		{
			name:    "fail guild channels",
			ctx:     &Context{Message: discord.Message{GuildID: 0}},
			allowed: GuildChannels,
			expect:  false,
		},
		{
			name:    "pass direct messages",
			ctx:     &Context{Message: discord.Message{GuildID: 0}},
			allowed: DirectMessages,
			expect:  true,
		},
		{
			name:    "fail direct messages",
			ctx:     &Context{Message: discord.Message{GuildID: 123}},
			allowed: DirectMessages,
			expect:  false,
		},
		{
			name:    "all channels",
			ctx:     &Context{Message: discord.Message{GuildID: 0}},
			allowed: AllChannels,
			expect:  true,
		},
		{
			name: "pass guild text",
			ctx: &Context{
				Message: discord.Message{GuildID: 123},
				DiscordDataProvider: mockDiscordDataProvider{
					ChannelReturn: &discord.Channel{Type: discord.GuildText},
				},
			},
			allowed: GuildTextChannels,
			expect:  true,
		},
		{
			name: "fail guild text",
			ctx: &Context{
				Message: discord.Message{GuildID: 0},
			},
			allowed: GuildTextChannels,
			expect:  false,
		},
		{
			name: "fail guild text",
			ctx: &Context{
				Message: discord.Message{GuildID: 123},
				DiscordDataProvider: mockDiscordDataProvider{
					ChannelReturn: &discord.Channel{Type: discord.GuildNews},
				},
			},
			allowed: GuildTextChannels,
			expect:  false,
		},
		{
			name:    "0 channel types",
			ctx:     &Context{Message: discord.Message{GuildID: 123}},
			allowed: 0,
			expect:  false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.allowed.Check(c.ctx)
			require.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}
