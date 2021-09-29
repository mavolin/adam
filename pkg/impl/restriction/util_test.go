package restriction

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func Test_assertChannelTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		ctx     *plugin.Context
		allowed plugin.ChannelTypes
		expect  error
	}{
		{
			name: "pass guild channels",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 123},
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.GuildChannels,
				}),
			},
			allowed: plugin.GuildChannels,
			expect:  nil,
		},
		{
			name: "fail guild channels",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.GuildTextChannels,
				}),
			},
			allowed: plugin.GuildChannels,
			expect:  NewFatalChannelTypesError(plugin.GuildTextChannels, i18n.NewFallbackLocalizer()),
		},
		{
			name: "pass direct messages",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 0},
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.DirectMessages,
				}),
			},
			allowed: plugin.DirectMessages,
			expect:  nil,
		},
		{
			name: "fail direct messages",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 123},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			allowed: plugin.DirectMessages,
			expect:  NewFatalChannelTypesError(plugin.DirectMessages, i18n.NewFallbackLocalizer()),
		},
		{
			name: "all channels",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 0},
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.DirectMessages,
				}),
			},
			allowed: plugin.AllChannels,
			expect:  nil,
		},
		{
			name: "pass guild text",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 123},
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{Type: discord.GuildText},
				},
			},
			allowed: plugin.GuildTextChannels,
			expect:  nil,
		},
		{
			name: "fail guild text - fatal",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.GuildChannels,
				}),
			},
			allowed: plugin.GuildTextChannels,
			expect:  NewFatalChannelTypesError(plugin.GuildTextChannels, i18n.NewFallbackLocalizer()),
		},
		{
			name: "fail guild text - not fatal",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 123},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.GuildChannels,
				}),
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{Type: discord.GuildNews},
				},
			},
			allowed: plugin.GuildTextChannels,
			expect:  NewChannelTypesError(plugin.GuildTextChannels, i18n.NewFallbackLocalizer()),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := assertChannelTypes(c.ctx, c.allowed)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_insertRoleSorted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		role   discord.Role
		roles  []discord.Role
		expect []discord.Role
	}{
		{
			name:   "empty roles",
			role:   discord.Role{Position: 3},
			roles:  nil,
			expect: []discord.Role{{Position: 3}},
		},
		{
			name:   "append role",
			role:   discord.Role{Position: 5},
			roles:  []discord.Role{{Position: 3}},
			expect: []discord.Role{{Position: 3}, {Position: 5}},
		},
		{
			name:   "insert role",
			role:   discord.Role{Position: 3},
			roles:  []discord.Role{{Position: 5}},
			expect: []discord.Role{{Position: 3}, {Position: 5}},
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := insertRoleSorted(c.role, c.roles)
			assert.Equal(t, c.expect, actual)
		})
	}
}
