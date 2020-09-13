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
		name        string
		ctx         *plugin.Context
		assertTypes plugin.ChannelTypes
		expect      error
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
			assertTypes: plugin.GuildChannels,
			expect:      nil,
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
			assertTypes: plugin.GuildChannels,
			expect:      newInvalidChannelTypeError(plugin.GuildTextChannels, mock.NewNoOpLocalizer(), true),
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
			assertTypes: plugin.DirectMessages,
			expect:      nil,
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
			assertTypes: plugin.DirectMessages,
			expect:      newInvalidChannelTypeError(plugin.DirectMessages, mock.NewNoOpLocalizer(), true),
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
			assertTypes: plugin.AllChannels,
			expect:      nil,
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
			assertTypes: plugin.GuildTextChannels,
			expect:      nil,
		},
		{
			name: "fail guild text",
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
			assertTypes: plugin.GuildTextChannels,
			expect:      newInvalidChannelTypeError(plugin.GuildTextChannels, mock.NewNoOpLocalizer(), true),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := assertChannelTypes(c.ctx, c.assertTypes, nil)
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
