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
			expect:  NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildTextChannels),
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
			expect:  NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.DirectMessages),
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
			expect:  NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildTextChannels),
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
			expect:  NewChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildTextChannels),
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
