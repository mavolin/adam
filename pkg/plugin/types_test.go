package plugin

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/mock"
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

func TestChannelTypes_Names(t *testing.T) {
	testCases := []struct {
		name         string
		channelTypes ChannelTypes
		expect       []string
	}{
		{
			name:         "guild text",
			channelTypes: GuildText,
			expect:       []string{"text channel"},
		},
		{
			name:         "guild news",
			channelTypes: GuildNews,
			expect:       []string{"announcement channel"},
		},
		{
			name:         "direct message",
			channelTypes: DirectMessage,
			expect:       []string{"direct message"},
		},
		{
			name:         "multiple",
			channelTypes: Guild,
			expect:       []string{"text channel", "announcement channel"},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.channelTypes.Names(mock.NewNoOpLocalizer())
			assert.Equal(t, c.expect, actual)
		})
	}
}
