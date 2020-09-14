package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
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
