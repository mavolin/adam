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

func TestNSFW(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    *plugin.Context
		expect error
	}{
		{
			name: "not a guild",
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
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
			},
			expect: newInvalidChannelTypeError(plugin.GuildChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name: "nsfw",
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
										ChannelTypes: plugin.DirectMessages,
									},
								},
							},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{
						NSFW: true,
					},
				},
			},
			expect: nil,
		},
		{
			name: "not nsfw",
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
										ChannelTypes: plugin.DirectMessages,
									},
								},
							},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{
						NSFW: false,
					},
				},
			},
			expect: ErrNotNSFWChannel,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := NSFW(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestGuildOwner(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    *plugin.Context
		expect error
	}{
		{
			name: "not a guild",
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
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
			},
			expect: newInvalidChannelTypeError(plugin.GuildChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name: "is owner",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
							Author: discord.User{
								ID: 456,
							},
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
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						OwnerID: 456,
					},
				},
			},
			expect: nil,
		},
		{
			name: "is not owner",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
							Author: discord.User{
								ID: 456,
							},
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
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						OwnerID: 789,
					},
				},
			},
			expect: ErrNotGuildOwner,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := GuildOwner(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestBotOwner(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    *plugin.Context
		expect error
	}{
		{
			name: "bot owner",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							Author: discord.User{
								ID: 123,
							},
						},
					},
				},
				BotOwnerIDs: []discord.UserID{123},
			},
			expect: nil,
		},
		{
			name: "not bot owner",
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							Author: discord.User{
								ID: 123,
							},
						},
					},
				},
				BotOwnerIDs: []discord.UserID{},
			},
			expect: ErrNotBotOwner,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := BotOwner(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestUsers(t *testing.T) {
	testCases := []struct {
		name    string
		userIDs []discord.UserID
		ctx     *plugin.Context
		expect  error
	}{
		{
			name:    "no users",
			userIDs: nil,
			expect:  nil,
		},
		{
			name:    "allowed",
			userIDs: []discord.UserID{123},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							Author: discord.User{
								ID: 123,
							},
						},
					},
				},
			},
			expect: nil,
		},
		{
			name:    "prohibited",
			userIDs: []discord.UserID{123},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							Author: discord.User{
								ID: 456,
							},
						},
					},
				},
			},
			expect: errors.DefaultFatalRestrictionError,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := Users(c.userIDs...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAllRoles(t *testing.T) {
	testCases := []struct {
		name    string
		roleIDs []discord.RoleID
		ctx     *plugin.Context
		expect  error
	}{
		{
			name:    "no roles",
			roleIDs: nil,
			expect:  nil,
		},
		{
			name:    "not a guild",
			roleIDs: []discord.RoleID{123},
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
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
			},
			expect: newInvalidChannelTypeError(plugin.GuildChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name:    "none missing",
			roleIDs: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 789,
						},
						Member: &discord.Member{
							RoleIDs: []discord.RoleID{123, 456},
						},
					},
				},
			},
			expect: nil,
		},
		{
			name:    "none missing - no roles from guild",
			roleIDs: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 789,
						},
						Member: &discord.Member{
							RoleIDs: []discord.RoleID{789},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{},
					},
				},
			},
			expect: errors.DefaultFatalRestrictionError,
		},
		{
			name:    "missing",
			roleIDs: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 789,
						},
						Member: &discord.Member{
							RoleIDs: []discord.RoleID{789},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{
								ID: 456,
							},
						},
					},
				},
			},
			expect: newAllMissingRolesError([]discord.RoleID{456}, mock.NewNoOpLocalizer()),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := AllRoles(c.roleIDs...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAnyRoles(t *testing.T) {
	testCases := []struct {
		name    string
		roleIDs []discord.RoleID
		ctx     *plugin.Context
		expect  error
	}{
		{
			name:    "no roles",
			roleIDs: nil,
			expect:  nil,
		},
		{
			name:    "not a guild",
			roleIDs: []discord.RoleID{123},
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
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
			},
			expect: newInvalidChannelTypeError(plugin.GuildChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name:    "none missing",
			roleIDs: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 789,
						},
						Member: &discord.Member{
							RoleIDs: []discord.RoleID{456},
						},
					},
				},
			},
			expect: nil,
		},
		{
			name:    "none missing from guild",
			roleIDs: []discord.RoleID{123},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 456,
						},
						Member: &discord.Member{
							RoleIDs: []discord.RoleID{789},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{
								ID: 789,
							},
						},
					},
				},
			},
			expect: errors.DefaultFatalRestrictionError,
		},
		{
			name:    "missing",
			roleIDs: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 789,
						},
						Member: &discord.Member{
							RoleIDs: []discord.RoleID{789},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{
								ID: 456,
							},
						},
					},
				},
			},
			expect: newAllMissingRolesError([]discord.RoleID{456}, mock.NewNoOpLocalizer()),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := AnyRole(c.roleIDs...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestChannels(t *testing.T) {
	testCases := []struct {
		name           string
		channelIDs     []discord.ChannelID
		channelsReturn []discord.Channel
		ctx            *plugin.Context
		expect         error
	}{
		{
			name:       "no channels",
			channelIDs: nil,
			expect:     nil,
		},
		{
			name:       "allowed",
			channelIDs: []discord.ChannelID{123},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 123,
						},
					},
				},
			},
			expect: nil,
		},
		{
			name:       "prohibited - direct message",
			channelIDs: []discord.ChannelID{123},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 456,
							GuildID:   0,
						},
					},
				},
			},
			expect: errors.DefaultFatalRestrictionError,
		},
		{
			name:       "prohibited - allowed channels in guild",
			channelIDs: []discord.ChannelID{123, 456, 789},
			channelsReturn: []discord.Channel{
				{
					ID: 456,
				},
				{
					ID: 789,
					Permissions: []discord.Overwrite{
						{
							ID:   678,
							Type: discord.OverwriteMember,
							Deny: discord.PermissionSendMessages,
						},
					},
				},
			},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 012,
							GuildID:   345,
						},
						Member: &discord.Member{
							User: discord.User{
								ID: 678,
							},
						},
					},
				},
				Localizer: mock.NewNoOpLocalizer(),
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						ID: 345,
						Roles: []discord.Role{
							{
								ID:          345,
								Permissions: discord.PermissionViewChannel | discord.PermissionSendMessages,
							},
						},
					},
				},
			},
			expect: newChannelsError([]discord.ChannelID{456}, mock.NewNoOpLocalizer()),
		},
		{
			name:       "prohibited - allowed channels not in guild",
			channelIDs: []discord.ChannelID{123, 456},
			channelsReturn: []discord.Channel{
				{
					ID: 789,
				},
			},
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 789,
							GuildID:   012,
						},
						Member: &discord.Member{},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						ID: 012,
						Roles: []discord.Role{
							{
								ID:          012,
								Permissions: discord.PermissionViewChannel | discord.PermissionSendMessages,
							},
						},
					},
				},
			},
			expect: errors.DefaultFatalRestrictionError,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := state.NewMocker(t)

			if c.channelsReturn != nil {
				m.Channels(c.ctx.GuildID, c.channelsReturn)
			}

			f := Channels(c.channelIDs...)

			actual := f(s, c.ctx)
			assert.Equal(t, c.expect, actual)

			m.Eval()
		})
	}
}

func TestBotPermissions(t *testing.T) {
	testCases := []struct {
		name   string
		perms  discord.Permissions
		ctx    *plugin.Context
		expect error
	}{
		{
			name:   "perms are 0",
			perms:  0,
			expect: nil,
		},
		{
			name:  "pass direct message",
			perms: discord.PermissionSendMessages | discord.PermissionViewChannel,
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
						},
					},
				},
			},
			expect: nil,
		},
		{
			name:  "fail direct message",
			perms: discord.PermissionAdministrator,
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
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
			},
			expect: newInvalidChannelTypeError(plugin.GuildChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name:  "pass guild",
			perms: discord.PermissionSendMessages,
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						ID: 123,
						Roles: []discord.Role{
							{
								ID:          123,
								Permissions: discord.PermissionViewChannel | discord.PermissionSendMessages,
							},
						},
					},
					ChannelReturn: &discord.Channel{},
					SelfReturn: &discord.Member{
						User: discord.User{
							ID: 456,
						},
					},
				},
			},
			expect: nil,
		},
		{
			name:  "fail guild",
			perms: discord.PermissionStream | discord.PermissionSendTTSMessages | discord.PermissionSendMessages,
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
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						ID: 123,
						Roles: []discord.Role{
							{
								ID:          123,
								Permissions: discord.PermissionViewChannel | discord.PermissionSendMessages,
							},
						},
					},
					ChannelReturn: &discord.Channel{},
					SelfReturn: &discord.Member{
						User: discord.User{
							ID: 456,
						},
					},
				},
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
			expect: newInsufficientBotPermissionsError(discord.PermissionStream|discord.PermissionSendTTSMessages,
				mock.NewNoOpLocalizer()),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := BotPermissions(c.perms)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestUserPermissions(t *testing.T) {
	testCases := []struct {
		name   string
		perms  discord.Permissions
		ctx    *plugin.Context
		expect error
	}{
		{
			name:   "perms are 0",
			perms:  0,
			expect: nil,
		},
		{
			name:  "pass direct message",
			perms: discord.PermissionSendMessages | discord.PermissionViewChannel,
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
						},
					},
				},
			},
			expect: nil,
		},
		{
			name:  "fail direct message",
			perms: discord.PermissionAdministrator,
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
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				},
			},
			expect: newInvalidChannelTypeError(plugin.GuildChannels, mock.NewNoOpLocalizer(), true),
		},
		{
			name:  "pass guild",
			perms: discord.PermissionSendMessages,
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
						},
						Member: &discord.Member{
							User: discord.User{
								ID: 456,
							},
						},
					},
				},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						ID: 123,
						Roles: []discord.Role{
							{
								ID:          123,
								Permissions: discord.PermissionViewChannel | discord.PermissionSendMessages,
							},
						},
					},
					ChannelReturn: &discord.Channel{},
				},
			},
			expect: nil,
		},
		{
			name:  "fail guild",
			perms: discord.PermissionStream | discord.PermissionSendTTSMessages | discord.PermissionSendMessages,
			ctx: &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 123,
						},
						Member: &discord.Member{
							User: discord.User{
								ID: 456,
							},
						},
					},
				},
				Localizer:         mock.NewNoOpLocalizer(),
				CommandIdentifier: ".abc",
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						ID: 123,
						Roles: []discord.Role{
							{
								ID:          123,
								Permissions: discord.PermissionViewChannel | discord.PermissionSendMessages,
							},
						},
					},
					ChannelReturn: &discord.Channel{},
				},
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
			expect: newInsufficientUserPermissionsError(discord.PermissionStream|discord.PermissionSendTTSMessages,
				mock.NewNoOpLocalizer()),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := UserPermissions(c.perms)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}
