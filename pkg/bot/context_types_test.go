package bot

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/state/store"
	"github.com/diamondburned/arikawa/v2/state/store/defaultstore"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

// =============================================================================
// plugin.Provider
// =====================================================================================

func TestCtxPluginProvider_PluginRepositories(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			repos: []plugin.Repository{
				{
					ProviderName: plugin.BuiltInProvider,
					Commands: []plugin.Command{
						mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
					},
				},
			},
			remProviders: []*pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return []plugin.Command{mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}}},
							[]plugin.Module{mock.Module{ModuleMeta: mock.ModuleMeta{Name: "ghi"}}},
							nil
					},
				},
			},
		}

		expect := []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands: []plugin.Command{
					mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
				},
			},
			{
				ProviderName: "another",
				Commands: []plugin.Command{
					mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
				},
				Modules: []plugin.Module{
					mock.Module{ModuleMeta: mock.ModuleMeta{Name: "ghi"}},
				},
			},
		}

		actual := p.PluginRepositories()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			repos: []plugin.Repository{
				{
					ProviderName: plugin.BuiltInProvider,
					Commands: []plugin.Command{
						mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
					},
				},
				{
					ProviderName: "another",
					Commands: []plugin.Command{
						mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
					},
					Modules: []plugin.Module{
						mock.Module{ModuleMeta: mock.ModuleMeta{Name: "ghi"}},
					},
				},
			},
			remProviders: nil,
		}

		expect := p.repos

		actual := p.PluginRepositories()
		assert.Equal(t, expect, actual)
	})
}

func TestCtxPluginProvider_Commands(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			repos: []plugin.Repository{
				{
					ProviderName: plugin.BuiltInProvider,
					Commands: []plugin.Command{
						mock.Command{
							CommandMeta: mock.CommandMeta{
								Name:         "abc",
								ChannelTypes: plugin.GuildChannels,
							},
						},
					},
				},
			},
			remProviders: []*pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return []plugin.Command{
							mock.Command{
								CommandMeta: mock.CommandMeta{
									Name:         "def",
									ChannelTypes: plugin.GuildChannels,
								},
							},
						}, nil, nil
					},
				},
			},
		}

		expect := []*plugin.RegisteredCommand{
			mock.GenerateRegisteredCommand(
				plugin.BuiltInProvider,
				mock.Command{
					CommandMeta: mock.CommandMeta{
						Name:         "abc",
						ChannelTypes: plugin.GuildChannels,
					},
				},
			),
			mock.GenerateRegisteredCommand(
				"another",
				mock.Command{
					CommandMeta: mock.CommandMeta{
						Name:         "def",
						ChannelTypes: plugin.GuildChannels,
					},
				},
			),
		}

		actual := p.Commands()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			commands: []*plugin.RegisteredCommand{
				mock.GenerateRegisteredCommand(
					plugin.BuiltInProvider,
					mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}}),
				mock.GenerateRegisteredCommand(
					"another",
					mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}}),
			},
		}

		expect := p.commands

		actual := p.Commands()
		assert.Equal(t, expect, actual)
	})
}

func TestCtxPluginProvider_Modules(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			repos: []plugin.Repository{
				{
					ProviderName: plugin.BuiltInProvider,
					Modules: []plugin.Module{
						mock.Module{
							ModuleMeta: mock.ModuleMeta{Name: "abc"},
							CommandsReturn: []plugin.Command{
								mock.Command{
									CommandMeta: mock.CommandMeta{
										Name:         "def",
										ChannelTypes: plugin.GuildChannels,
									},
								},
							},
						},
					},
				},
			},
			remProviders: []*pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return nil, []plugin.Module{
							mock.Module{
								ModuleMeta: mock.ModuleMeta{Name: "ghi"},
								CommandsReturn: []plugin.Command{
									mock.Command{
										CommandMeta: mock.CommandMeta{
											Name:         "jkl",
											ChannelTypes: plugin.GuildChannels,
										},
									},
								},
							},
						}, nil
					},
				},
			},
		}

		expect := []*plugin.RegisteredModule{
			mock.GenerateRegisteredModule(
				plugin.BuiltInProvider,
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "abc"},
					CommandsReturn: []plugin.Command{
						mock.Command{
							CommandMeta: mock.CommandMeta{
								Name:         "def",
								ChannelTypes: plugin.GuildChannels,
							},
						},
					},
				},
			),
			mock.GenerateRegisteredModule(
				"another",
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "ghi"},
					CommandsReturn: []plugin.Command{
						mock.Command{
							CommandMeta: mock.CommandMeta{
								Name:         "jkl",
								ChannelTypes: plugin.GuildChannels,
							},
						},
					},
				},
			),
		}

		actual := p.Modules()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			modules: []*plugin.RegisteredModule{
				mock.GenerateRegisteredModule(
					plugin.BuiltInProvider,
					mock.Module{
						ModuleMeta: mock.ModuleMeta{Name: "abc"},
						CommandsReturn: []plugin.Command{
							mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
						},
					}),
				mock.GenerateRegisteredModule(
					"another",
					mock.Module{
						ModuleMeta: mock.ModuleMeta{Name: "ghi"},
						CommandsReturn: []plugin.Command{
							mock.Command{CommandMeta: mock.CommandMeta{Name: "jkl"}},
						},
					}),
			},
		}

		expect := p.modules

		actual := p.Modules()
		assert.Equal(t, expect, actual)
	})
}

func TestCtxPluginProvider_Command(t *testing.T) {
	t.Run("top-level", func(t *testing.T) {
		p := &ctxPluginProvider{
			commands: []*plugin.RegisteredCommand{
				mock.GenerateRegisteredCommand(
					plugin.BuiltInProvider,
					mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}}),
				mock.GenerateRegisteredCommand(
					"another",
					mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}}),
			},
		}

		expect := p.commands[1]
		actual := p.Command(".def")
		assert.Equal(t, expect, actual)
	})

	t.Run("nested", func(t *testing.T) {
		p := &ctxPluginProvider{
			modules: []*plugin.RegisteredModule{
				mock.GenerateRegisteredModule(
					plugin.BuiltInProvider,
					mock.Module{
						ModuleMeta: mock.ModuleMeta{Name: "abc"},
						CommandsReturn: []plugin.Command{
							mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
						},
					}),
				mock.GenerateRegisteredModule(
					"another",
					mock.Module{
						ModuleMeta: mock.ModuleMeta{Name: "ghi"},
						CommandsReturn: []plugin.Command{
							mock.Command{CommandMeta: mock.CommandMeta{Name: "jkl"}},
							mock.Command{CommandMeta: mock.CommandMeta{Name: "mno"}},
						},
					}),
			},
		}

		expect := p.modules[1].Commands[1]
		actual := p.Command(".ghi.mno")
		assert.Equal(t, expect, actual)
	})
}

func TestCtxPluginProvider_Module(t *testing.T) {
	p := &ctxPluginProvider{
		modules: []*plugin.RegisteredModule{
			mock.GenerateRegisteredModule(
				plugin.BuiltInProvider,
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "abc"},
					ModulesReturn: []plugin.Module{
						mock.Module{
							ModuleMeta: mock.ModuleMeta{Name: "def"},
							CommandsReturn: []plugin.Command{
								mock.Command{CommandMeta: mock.CommandMeta{Name: "ghi"}},
							},
						},
					},
				}),
			mock.GenerateRegisteredModule(
				"another",
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "jkl"},
					ModulesReturn: []plugin.Module{
						mock.Module{
							ModuleMeta: mock.ModuleMeta{Name: "mno"},
							CommandsReturn: []plugin.Command{
								mock.Command{CommandMeta: mock.CommandMeta{Name: "pqr"}},
							},
						},
					},
				}),
		},
	}

	expect := p.modules[1].Modules[0]
	actual := p.Module(".jkl.mno")
	assert.Equal(t, expect, actual)
}

func TestCtxPluginProvider_FindCommand(t *testing.T) {
	t.Run("top-level", func(t *testing.T) {
		p := &ctxPluginProvider{
			commands: []*plugin.RegisteredCommand{
				mock.GenerateRegisteredCommand(
					plugin.BuiltInProvider,
					mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}}),
				mock.GenerateRegisteredCommand(
					"another",
					mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}}),
			},
		}

		expect := p.commands[1]
		actual := p.FindCommand("def  \t")
		assert.Equal(t, expect, actual)
	})

	t.Run("nested", func(t *testing.T) {
		p := &ctxPluginProvider{
			modules: []*plugin.RegisteredModule{
				mock.GenerateRegisteredModule(
					plugin.BuiltInProvider,
					mock.Module{
						ModuleMeta: mock.ModuleMeta{Name: "abc"},
						CommandsReturn: []plugin.Command{
							mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
						},
					}),
				mock.GenerateRegisteredModule(
					"another",
					mock.Module{
						ModuleMeta: mock.ModuleMeta{Name: "ghi"},
						CommandsReturn: []plugin.Command{
							mock.Command{CommandMeta: mock.CommandMeta{Name: "jkl"}},
							mock.Command{CommandMeta: mock.CommandMeta{Name: "mno"}},
						},
					}),
			},
		}

		expect := p.modules[1].Commands[1]
		actual := p.FindCommand(" ghi  mno")
		assert.Equal(t, expect, actual)
	})
}

func TestCtxPluginProvider_FindModule(t *testing.T) {
	p := &ctxPluginProvider{
		modules: []*plugin.RegisteredModule{
			mock.GenerateRegisteredModule(
				plugin.BuiltInProvider,
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "abc"},
					ModulesReturn: []plugin.Module{
						mock.Module{
							ModuleMeta: mock.ModuleMeta{Name: "def"},
							CommandsReturn: []plugin.Command{
								mock.Command{CommandMeta: mock.CommandMeta{Name: "ghi"}},
							},
						},
					},
				}),
			mock.GenerateRegisteredModule(
				"another",
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "jkl"},
					ModulesReturn: []plugin.Module{
						mock.Module{
							ModuleMeta: mock.ModuleMeta{Name: "mno"},
							CommandsReturn: []plugin.Command{
								mock.Command{CommandMeta: mock.CommandMeta{Name: "pqr"}},
							},
						},
					},
				}),
		},
	}

	expect := p.modules[1].Modules[0]
	actual := p.FindModule("jkl \nmno")
	assert.Equal(t, expect, actual)
}

func TestCtxPluginProvider_UnavailablePluginProviders(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			repos: []plugin.Repository{
				{
					ProviderName: plugin.BuiltInProvider,
					Commands: []plugin.Command{
						mock.Command{
							CommandMeta: mock.CommandMeta{
								Name:         "abc",
								ChannelTypes: plugin.GuildChannels,
							},
						},
					},
				},
			},
			remProviders: []*pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return nil, nil, errors.New("abc")
					},
				},
			},
		}

		expect := []plugin.UnavailablePluginProvider{
			{
				Name:  "another",
				Error: errors.New("abc"),
			},
		}

		actual := p.UnavailablePluginProviders()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := &ctxPluginProvider{
			repos: []plugin.Repository{
				{
					ProviderName: plugin.BuiltInProvider,
					Commands: []plugin.Command{
						mock.Command{
							CommandMeta: mock.CommandMeta{
								Name:         "abc",
								ChannelTypes: plugin.GuildChannels,
							},
						},
					},
				},
			},
			remProviders: nil,
			unavailableProviders: []plugin.UnavailablePluginProvider{
				{
					Name:  "another",
					Error: errors.New("abc"),
				},
			},
		}

		expect := p.unavailableProviders

		actual := p.UnavailablePluginProviders()
		assert.Equal(t, expect, actual)
	})
}

// =============================================================================
// plugin.ErrorHandler
// =====================================================================================

func Test_newCtxErrorHandler(t *testing.T) {
	var called bool

	f := func(error, *state.State, *plugin.Context) { called = true }

	h := newCtxErrorHandler(nil, nil, f)
	h(errors.New("abc"))

	assert.True(t, called, "wrapped error handler was not called")
}

func TestCtxErrorHandler_HandleError(t *testing.T) {
	var actual error

	var h ctxErrorHandler = func(err error) { actual = err }

	expect := errors.New("Abort. Retry. Fail.") //nolint:golint

	h.HandleError(expect)

	assert.Equal(t, expect, actual)
}

// =============================================================================
// plugin.DiscordDataProvider
// =====================================================================================

func TestDiscordDataProvider_GuildAsync(t *testing.T) {
	t.Run("cached", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Guild{ID: 123}

		s.Cabinet = store.Cabinet{GuildStore: defaultstore.NewGuild()}

		err := s.Cabinet.GuildSet(*expect)
		require.NoError(t, err)

		p := &discordDataProvider{
			s:       s,
			guildID: expect.ID,
		}

		actual, err := p.GuildAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("fetch", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Guild{
			ID:                     123,
			OwnerID:                1,
			RulesChannelID:         1,
			PublicUpdatesChannelID: 1,
		}

		m.Guild(*expect)

		p := &discordDataProvider{
			s:       s,
			guildID: expect.ID,
		}

		actual, err := p.GuildAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestDiscordDataProvider_ChannelAsync(t *testing.T) {
	t.Run("cached", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Channel{ID: 123, GuildID: 456}

		s.Cabinet = store.Cabinet{ChannelStore: defaultstore.NewChannel()}

		err := s.Cabinet.ChannelSet(*expect)
		require.NoError(t, err)

		p := &discordDataProvider{
			s:         s,
			channelID: expect.ID,
		}

		actual, err := p.ChannelAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("fetch", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		expect := &discord.Channel{
			ID: 123,
		}

		m.Channel(*expect)

		p := &discordDataProvider{
			s:         s,
			channelID: expect.ID,
		}

		actual, err := p.ChannelAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func TestDiscordDataProvider_MemberAsync(t *testing.T) {
	t.Run("cached", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := &discord.Member{
			User: discord.User{ID: 456},
		}

		s.Cabinet = store.Cabinet{MemberStore: defaultstore.NewMember()}

		err := s.Cabinet.MemberSet(guildID, *expect)
		require.NoError(t, err)

		p := &discordDataProvider{
			s:       s,
			guildID: guildID,
			selfID:  expect.User.ID,
		}

		actual, err := p.SelfAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("fetch", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := &discord.Member{
			User: discord.User{ID: 456},
		}

		m.Member(guildID, *expect)

		p := &discordDataProvider{
			s:       s,
			guildID: guildID,
			selfID:  expect.User.ID,
		}

		actual, err := p.SelfAsync()()
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}
