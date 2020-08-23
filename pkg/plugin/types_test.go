package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
)

func TestChannelType_Has(t *testing.T) {
	successCases := []struct {
		name        string
		channelType ChannelType
		target      discord.ChannelType
		expect      bool
	}{
		{
			name:        "all",
			channelType: All,
			target:      discord.GuildText,
		},
		{
			name:        "GuildText",
			channelType: GuildText,
			target:      discord.GuildText,
		},
		{
			name:        "DirectMessage",
			channelType: DirectMessage,
			target:      discord.DirectMessage,
		},
		{
			name:        "GuildNews",
			channelType: GuildNews,
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
		has := GuildText.Has(discord.DirectMessage)
		assert.False(t, has)
	})
}
