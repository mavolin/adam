package restriction

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/plugin"
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
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
						},
					},
				},
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.GuildChannels,
									},
								},
							},
						},
					},
				},
			},
			allowed: plugin.GuildChannels,
			expect:  nil,
		},
		{
			name: "fail guild channels",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
						},
					},
				},
				Localizer:         mock.NewNoOpLocalizer(),
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.GuildTextChannels,
									},
								},
							},
						},
					},
				},
			},
			allowed: plugin.GuildChannels,
			expect:  newInvalidChannelTypeError(plugin.GuildTextChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name: "pass direct messages",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
						},
					},
				},
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.DirectMessages,
									},
								},
							},
						},
					},
				},
			},
			allowed: plugin.DirectMessages,
			expect:  nil,
		},
		{
			name: "fail direct messages",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
						},
					},
				},
				Localizer:         mock.NewNoOpLocalizer(),
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
			},
			allowed: plugin.DirectMessages,
			expect:  newInvalidChannelTypeError(plugin.DirectMessages, mock.NewNoOpLocalizer(), true),
		},
		{
			name: "all channels",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
						},
					},
				},
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.DirectMessages,
									},
								},
							},
						},
					},
				},
			},
			allowed: plugin.AllChannels,
			expect:  nil,
		},
		{
			name: "pass guild text",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
						},
					},
				},
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{
						Type: discord.GuildText,
					},
				},
			},
			allowed: plugin.GuildTextChannels,
			expect:  nil,
		},
		{
			name: "fail guild text - fatal",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
						},
					},
				},
				Localizer:         mock.NewNoOpLocalizer(),
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.GuildChannels,
									},
								},
							},
						},
					},
				},
			},
			allowed: plugin.GuildTextChannels,
			expect:  newInvalidChannelTypeError(plugin.GuildTextChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name: "fail guild text - not fatal",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
						},
					},
				},
				Localizer:         mock.NewNoOpLocalizer(),
				CommandIdentifier: ".abc",
				Provider: mock.PluginProvider{
					AllCommandsReturn: []plugin.CommandRepository{
						{
							Commands: []plugin.Command{
								mock.Command{
									MetaReturn: mock.CommandMeta{
										Name:         "abc",
										ChannelTypes: plugin.GuildChannels,
									},
								},
							},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{
						Type: discord.GuildNews,
					},
				},
			},
			allowed: plugin.GuildTextChannels,
			expect:  newInvalidChannelTypeError(plugin.GuildTextChannels, mock.NewNoOpLocalizer(), false),
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
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						GuildID: 123,
					},
				},
			},
			Localizer:         mock.NewNoOpLocalizer(),
			CommandIdentifier: ".abc",
			Provider: mock.PluginProvider{
				AllCommandsReturn: []plugin.CommandRepository{
					{
						Commands: []plugin.Command{
							mock.Command{
								MetaReturn: mock.CommandMeta{
									Name:         "abc",
									ChannelTypes: plugin.GuildChannels,
								},
							},
						},
					},
				},
			},
			ErrorHandler: mock.NewErrorHandler().
				ExpectSilentError(noRemainingError),
		}

		actual := assertChannelTypes(ctx, plugin.DirectMessages, noRemainingError)
		assert.Equal(t, errors.DefaultFatalRestrictionError, actual)
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
			name: "can not manage",
			target: discord.Role{
				Position: 2,
			},
			guild: &discord.Guild{
				Roles: []discord.Role{
					{
						ID:       123,
						Position: 1,
					},
				},
			},
			member: &discord.Member{
				RoleIDs: []discord.RoleID{123},
			},
			expect: false,
		},
		{
			name: "no manage roles permission",
			target: discord.Role{
				Position: 1,
			},
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
				User: discord.User{
					ID: 789,
				},
				RoleIDs: []discord.RoleID{123},
			},
			expect: false,
		},
		{
			name: "pass",
			target: discord.Role{
				Position: 1,
			},
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
				User: discord.User{
					ID: 789,
				},
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
			name: "empty roles",
			role: discord.Role{
				Position: 3,
			},
			roles: nil,
			expect: []discord.Role{
				{
					Position: 3,
				},
			},
		},
		{
			name: "roles filled",
			role: discord.Role{
				Position: 5,
			},
			roles: []discord.Role{
				{
					Position: 3,
				},
			},
			expect: []discord.Role{
				{
					Position: 3,
				},
				{
					Position: 5,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := insertRoleSorted(c.role, c.roles)
			assert.Equal(t, c.expect, actual)
		})
	}
}
