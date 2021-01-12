package restriction

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func Test_assertChannelTypes(t *testing.T) {
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
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{ChannelTypes: plugin.GuildChannels},
				}),
			},
			allowed: plugin.GuildChannels,
			expect:  nil,
		},
		{
			name: "fail guild channels",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.FallbackLocalizer,
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{ChannelTypes: plugin.GuildTextChannels},
				}),
			},
			allowed: plugin.GuildChannels,
			expect:  newInvalidChannelTypeError(plugin.GuildTextChannels, i18n.FallbackLocalizer, true),
		},
		{
			name: "pass direct messages",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 0},
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{ChannelTypes: plugin.DirectMessages},
				}),
			},
			allowed: plugin.DirectMessages,
			expect:  nil,
		},
		{
			name: "fail direct messages",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 123},
				Localizer: i18n.FallbackLocalizer,
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{ChannelTypes: plugin.AllChannels},
				}),
			},
			allowed: plugin.DirectMessages,
			expect:  newInvalidChannelTypeError(plugin.DirectMessages, i18n.FallbackLocalizer, true),
		},
		{
			name: "all channels",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 0},
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{ChannelTypes: plugin.DirectMessages},
				}),
			},
			allowed: plugin.AllChannels,
			expect:  nil,
		},
		{
			name: "pass guild text",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 123},
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{
						ChannelTypes: plugin.AllChannels,
					},
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
				Localizer: i18n.FallbackLocalizer,
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{ChannelTypes: plugin.GuildChannels},
				}),
			},
			allowed: plugin.GuildTextChannels,
			expect:  newInvalidChannelTypeError(plugin.GuildTextChannels, i18n.FallbackLocalizer, true),
		},
		{
			name: "fail guild text - not fatal",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 123},
				Localizer: i18n.FallbackLocalizer,
				InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
					CommandMeta: mock.CommandMeta{ChannelTypes: plugin.GuildChannels},
				}),
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{Type: discord.GuildNews},
				},
			},
			allowed: plugin.GuildTextChannels,
			expect:  newInvalidChannelTypeError(plugin.GuildTextChannels, i18n.FallbackLocalizer, false),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := assertChannelTypes(c.ctx, c.allowed, nil)
			assert.Equal(t, c.expect, actual)
		})
	}

	t.Run("fail - no remaining", func(t *testing.T) {
		noRemainingError := errors.New("no remaining error")

		ctx := &plugin.Context{
			Message:   discord.Message{GuildID: 123},
			Localizer: i18n.FallbackLocalizer,
			InvokedCommand: mock.GenerateRegisteredCommand("built_in", mock.Command{
				CommandMeta: mock.CommandMeta{ChannelTypes: plugin.GuildChannels},
			}),
			ErrorHandler: mock.NewErrorHandler().
				ExpectSilentError(noRemainingError),
		}

		actual := assertChannelTypes(ctx, plugin.DirectMessages, noRemainingError)
		assert.Equal(t, plugin.DefaultFatalRestrictionError, actual)
	})
}

func Test_canManageRole(t *testing.T) {
	testCases := []struct {
		name   string
		target discord.Role
		guild  *discord.Guild
		member *discord.Member
		expect bool
	}{
		{
			name:   "can not manage",
			target: discord.Role{Position: 2},
			guild: &discord.Guild{
				Roles: []discord.Role{
					{
						ID:       123,
						Position: 1,
					},
				},
			},
			member: &discord.Member{RoleIDs: []discord.RoleID{123}},
			expect: false,
		},
		{
			name:   "no manage roles permission",
			target: discord.Role{Position: 1},
			guild: &discord.Guild{
				Roles: []discord.Role{
					{
						ID:          123,
						Position:    2,
						Permissions: 0,
					},
				},
				OwnerID: 456,
			},
			member: &discord.Member{
				User:    discord.User{ID: 789},
				RoleIDs: []discord.RoleID{123},
			},
			expect: false,
		},
		{
			name:   "pass",
			target: discord.Role{Position: 1},
			guild: &discord.Guild{
				Roles: []discord.Role{
					{
						ID:          123,
						Position:    2,
						Permissions: discord.PermissionManageRoles,
					},
				},
				OwnerID: 456,
			},
			member: &discord.Member{
				User:    discord.User{ID: 789},
				RoleIDs: []discord.RoleID{123},
			},
			expect: true,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := canManageRole(c.target, c.guild, c.member)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_insertRoleSorted(t *testing.T) {
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
		t.Run(c.name, func(t *testing.T) {
			actual := insertRoleSorted(c.role, c.roles)
			assert.Equal(t, c.expect, actual)
		})
	}
}
