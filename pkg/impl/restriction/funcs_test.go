package restriction

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestNSFW(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		ctx    *plugin.Context
		expect error
	}{
		{
			name: "not a guild",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			expect: plugin.NewFatalRestrictionErrorl(nsfwChannelError),
		},
		{
			name: "nsfw",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 123},
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.DirectMessages,
				}),
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{NSFW: true},
				},
			},
			expect: nil,
		},
		{
			name: "not nsfw",
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 123},
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.DirectMessages,
				}),
				DiscordDataProvider: mock.DiscordDataProvider{
					ChannelReturn: &discord.Channel{NSFW: false},
				},
			},
			expect: plugin.NewRestrictionErrorl(nsfwChannelError),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := NSFW(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestGuildOwner(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		ctx    *plugin.Context
		expect error
	}{
		{
			name: "not a guild",
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			expect: NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildChannels),
		},
		{
			name: "is owner",
			ctx: &plugin.Context{
				Message: discord.Message{
					GuildID: 123,
					Author:  discord.User{ID: 456},
				},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{OwnerID: 456},
				},
			},
			expect: nil,
		},
		{
			name: "is not owner",
			ctx: &plugin.Context{
				Message: discord.Message{
					GuildID: 123,
					Author:  discord.User{ID: 456},
				},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{OwnerID: 789},
				},
			},
			expect: plugin.NewFatalRestrictionErrorl(guildOwnerError),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := GuildOwner(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestBotOwner(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		ctx    *plugin.Context
		expect error
	}{
		{
			name: "bot owner",
			ctx: &plugin.Context{
				Message:     discord.Message{Author: discord.User{ID: 123}},
				BotOwnerIDs: []discord.UserID{123},
			},
			expect: nil,
		},
		{
			name: "not bot owner",
			ctx: &plugin.Context{
				Message:     discord.Message{Author: discord.User{ID: 123}},
				BotOwnerIDs: []discord.UserID{},
			},
			expect: plugin.NewFatalRestrictionErrorl(botOwnerError),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := BotOwner(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestUsers(t *testing.T) {
	t.Parallel()

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
				Message: discord.Message{Author: discord.User{ID: 123}},
			},
			expect: nil,
		},
		{
			name:    "prohibited",
			userIDs: []discord.UserID{123},
			ctx: &plugin.Context{
				Message: discord.Message{Author: discord.User{ID: 456}},
			},
			expect: plugin.DefaultFatalRestrictionError,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			f := Users(c.userIDs...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAllRoles(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		allowed []discord.RoleID
		ctx     *plugin.Context
		expect  error
	}{
		{
			name:    "no roles",
			allowed: nil,
			expect:  nil,
		},
		{
			name:    "not a guild",
			allowed: []discord.RoleID{123},
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			expect: NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildChannels),
		},
		{
			name:    "none missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member: &discord.Member{
					RoleIDs: []discord.RoleID{123, 456},
				},
			},
			expect: nil,
		},
		{
			name:    "none missing - no roles from guild",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{12}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{},
					},
				},
			},
			expect: plugin.DefaultFatalRestrictionError,
		},
		{
			name:    "missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 789},
				Member:    &discord.Member{RoleIDs: []discord.RoleID{12}},
				Localizer: i18n.NewFallbackLocalizer(),
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						OwnerID: 345,
						Roles:   []discord.Role{{ID: 12}, {ID: 456}},
					},
				},
			},
			expect: NewAllMissingRolesError(i18n.NewFallbackLocalizer(), discord.Role{ID: 456}),
		},
		{
			name:    "missing - can manage",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{12, 345}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{ID: 123, Position: 3},
							{ID: 456, Position: 5},
							{ID: 12, Position: 6},
							{ID: 345, Position: 2, Permissions: discord.PermissionManageRoles},
						},
					},
				},
			},
			expect: nil,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			f := AllRoles(c.allowed...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestMustAllRoles(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		allowed []discord.RoleID
		ctx     *plugin.Context
		expect  error
	}{
		{
			name:    "no roles",
			allowed: nil,
			expect:  nil,
		},
		{
			name:    "not a guild",
			allowed: []discord.RoleID{123},
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			expect: NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildChannels),
		},
		{
			name:    "none missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{123, 456}},
			},
			expect: nil,
		},
		{
			name:    "none missing - no roles from guild",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{12}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{Roles: []discord.Role{}},
				},
			},
			expect: plugin.DefaultFatalRestrictionError,
		},
		{
			name:    "missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{12}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{Roles: []discord.Role{{ID: 456}}},
				},
			},
			expect: NewAllMissingRolesError(i18n.NewFallbackLocalizer(), discord.Role{ID: 456}),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			f := MustAllRoles(c.allowed...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAnyRole(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		allowed []discord.RoleID
		ctx     *plugin.Context
		expect  error
	}{
		{
			name:    "no roles",
			allowed: nil,
			expect:  nil,
		},
		{
			name:    "not a guild",
			allowed: []discord.RoleID{123},
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			expect: NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildChannels),
		},
		{
			name:    "none missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{456}},
			},
			expect: nil,
		},
		{
			name:    "none missing from guild",
			allowed: []discord.RoleID{123},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 456},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{789}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{Roles: []discord.Role{{ID: 789}}},
				},
			},
			expect: plugin.DefaultFatalRestrictionError,
		},
		{
			name:    "missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{12}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles:   []discord.Role{{ID: 456, Position: 1}},
						OwnerID: 345,
					},
				},
			},
			expect: NewAnyMissingRolesError(i18n.NewFallbackLocalizer(), discord.Role{ID: 456, Position: 1}),
		},
		{
			name:    "missing - can manage",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{12, 345}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{ID: 123, Position: 3},
							{ID: 456, Position: 5},
							{ID: 12, Position: 4},
							{ID: 345, Position: 2, Permissions: discord.PermissionManageRoles},
						},
					},
				},
			},
			expect: nil,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			f := AnyRole(c.allowed...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestMustAnyRole(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		allowed []discord.RoleID
		ctx     *plugin.Context
		expect  error
	}{
		{
			name:    "no roles",
			allowed: nil,
			expect:  nil,
		},
		{
			name:    "not a guild",
			allowed: []discord.RoleID{123},
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			expect: NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildChannels),
		},
		{
			name:    "none missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{456}},
			},
			expect: nil,
		},
		{
			name:    "none missing from guild",
			allowed: []discord.RoleID{123},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 456},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{789}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{Roles: []discord.Role{{ID: 789}}},
				},
			},
			expect: plugin.DefaultFatalRestrictionError,
		},
		{
			name:    "missing",
			allowed: []discord.RoleID{123, 456},
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 789},
				Member:  &discord.Member{RoleIDs: []discord.RoleID{12}},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						Roles: []discord.Role{
							{ID: 456, Position: 2},
							{ID: 12, Position: 1},
						},
					},
				},
			},
			expect: NewAnyMissingRolesError(i18n.NewFallbackLocalizer(), discord.Role{ID: 456, Position: 2}),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			f := MustAnyRole(c.allowed...)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestChannels(t *testing.T) {
	t.Parallel()

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
			ctx:        &plugin.Context{Message: discord.Message{ChannelID: 123}},
			expect:     nil,
		},
		{
			name:       "prohibited - direct message",
			channelIDs: []discord.ChannelID{123},
			ctx: &plugin.Context{
				Message: discord.Message{
					ChannelID: 456,
					GuildID:   0,
				},
			},
			expect: plugin.DefaultFatalRestrictionError,
		},
		{
			name:       "prohibited - allowed channels in guild",
			channelIDs: []discord.ChannelID{123, 456, 789},
			channelsReturn: []discord.Channel{
				{ID: 456},
				{
					ID: 789,
					Overwrites: []discord.Overwrite{
						{
							ID:   678,
							Type: discord.OverwriteMember,
							Deny: discord.PermissionSendMessages,
						},
					},
				},
			},
			ctx: &plugin.Context{
				Message:   discord.Message{ChannelID: 12, GuildID: 345},
				Member:    &discord.Member{User: discord.User{ID: 678}},
				Localizer: i18n.NewFallbackLocalizer(),
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
			expect: NewChannelsError(i18n.NewFallbackLocalizer(), 456),
		},
		{
			name:           "prohibited - allowed channels not in guild",
			channelIDs:     []discord.ChannelID{123, 456},
			channelsReturn: []discord.Channel{{ID: 789}},
			ctx: &plugin.Context{
				Message: discord.Message{ChannelID: 789, GuildID: 12},
				Member:  &discord.Member{},
				DiscordDataProvider: mock.DiscordDataProvider{
					GuildReturn: &discord.Guild{
						ID: 12,
						Roles: []discord.Role{
							{
								ID:          12,
								Permissions: discord.PermissionViewChannel | discord.PermissionSendMessages,
							},
						},
					},
				},
			},
			expect: plugin.DefaultFatalRestrictionError,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			m, s := state.NewMocker(t)

			if c.channelsReturn != nil {
				m.Channels(c.ctx.GuildID, c.channelsReturn)
			}

			f := Channels(c.channelIDs...)

			actual := f(s, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestUserPermissions(t *testing.T) {
	t.Parallel()

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
			name:   "pass direct message",
			perms:  discord.PermissionSendMessages | discord.PermissionViewChannel,
			ctx:    &plugin.Context{Message: discord.Message{GuildID: 0}},
			expect: nil,
		},
		{
			name:  "fail direct message",
			perms: discord.PermissionAdministrator,
			ctx: &plugin.Context{
				Message:   discord.Message{GuildID: 0},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
			},
			expect: NewFatalChannelTypesError(i18n.NewFallbackLocalizer(), plugin.GuildChannels),
		},
		{
			name:  "pass guild",
			perms: discord.PermissionSendMessages,
			ctx: &plugin.Context{
				Message: discord.Message{GuildID: 123},
				Member:  &discord.Member{User: discord.User{ID: 456}},
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
				Message:   discord.Message{GuildID: 123},
				Member:    &discord.Member{User: discord.User{ID: 456}},
				Localizer: i18n.NewFallbackLocalizer(),
				InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					ChannelTypes: plugin.AllChannels,
				}),
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
			expect: NewUserPermissionsError(i18n.NewFallbackLocalizer(),
				discord.PermissionStream|discord.PermissionSendTTSMessages),
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			f := UserPermissions(c.perms)

			actual := f(nil, c.ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}
