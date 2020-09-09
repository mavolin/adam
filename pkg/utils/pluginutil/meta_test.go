package pluginutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestChannelTypes(t *testing.T) {
	testCases := []struct {
		name           string
		id             plugin.Identifier
		pluginProvider plugin.Provider
		expect         plugin.ChannelTypes
	}{
		{
			name: "top-level command",
			id:   ".abc",
			pluginProvider: mock.PluginProvider{
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
			expect: plugin.GuildChannels,
		},
		{
			name: "command in module",
			id:   ".abc.def",
			pluginProvider: mock.PluginProvider{
				AllModulesReturn: []plugin.ModuleRepository{
					{
						Modules: []plugin.Module{
							mock.Module{
								MetaReturn: mock.ModuleMeta{
									Name:         "abc",
									ChannelTypes: plugin.DirectMessages,
								},
								CommandsReturn: []plugin.Command{
									mock.Command{
										MetaReturn: mock.CommandMeta{
											Name:         "def",
											ChannelTypes: plugin.GuildTextChannels,
										},
									},
								},
							},
						},
					},
				},
			},
			expect: plugin.GuildTextChannels,
		},
		{
			name: "module",
			id:   ".abc.def",
			pluginProvider: mock.PluginProvider{
				AllModulesReturn: []plugin.ModuleRepository{
					{
						Modules: []plugin.Module{
							mock.Module{
								MetaReturn: mock.ModuleMeta{
									Name:         "abc",
									ChannelTypes: plugin.DirectMessages,
								},
								CommandsReturn: []plugin.Command{
									mock.Command{
										MetaReturn: mock.CommandMeta{
											Name:         "def",
											ChannelTypes: 0,
										},
									},
								},
							},
						},
					},
				},
			},
			expect: plugin.DirectMessages,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := ChannelTypes(c.id, c.pluginProvider)
			require.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}
