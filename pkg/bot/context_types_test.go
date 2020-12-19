package bot

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"

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
					Defaults: plugin.Defaults{BotPermissions: discord.PermissionSendMessages},
				},
			},
			remProviders: []pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return []plugin.Command{mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}}},
							[]plugin.Module{mock.Module{ModuleMeta: mock.ModuleMeta{Name: "ghi"}}},
							nil
					},
					defaults: plugin.Defaults{BotPermissions: discord.PermissionConnect},
				},
			},
		}

		expect := []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands: []plugin.Command{
					mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
				},
				Defaults: plugin.Defaults{BotPermissions: discord.PermissionSendMessages},
			},
			{
				ProviderName: "another",
				Commands: []plugin.Command{
					mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
				},
				Modules: []plugin.Module{
					mock.Module{ModuleMeta: mock.ModuleMeta{Name: "ghi"}},
				},
				Defaults: plugin.Defaults{BotPermissions: discord.PermissionConnect},
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
					Defaults: plugin.Defaults{BotPermissions: discord.PermissionSendMessages},
				},
				{
					ProviderName: "another",
					Commands: []plugin.Command{
						mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
					},
					Modules: []plugin.Module{
						mock.Module{ModuleMeta: mock.ModuleMeta{Name: "ghi"}},
					},
					Defaults: plugin.Defaults{BotPermissions: discord.PermissionConnect},
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
						mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
					},
					Defaults: plugin.Defaults{BotPermissions: discord.PermissionSendMessages},
				},
			},
			remProviders: []pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return []plugin.Command{mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}}}, nil, nil
					},
					defaults: plugin.Defaults{BotPermissions: discord.PermissionConnect},
				},
			},
		}

		expect := []*plugin.RegisteredCommand{
			mock.GenerateRegisteredCommandWithDefaults(
				plugin.BuiltInProvider,
				mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
				plugin.Defaults{BotPermissions: discord.PermissionSendMessages}),
			mock.GenerateRegisteredCommandWithDefaults(
				"another",
				mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
				plugin.Defaults{BotPermissions: discord.PermissionConnect}),
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
								mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
							},
						},
					},
					Defaults: plugin.Defaults{BotPermissions: discord.PermissionSendMessages},
				},
			},
			remProviders: []pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return nil, []plugin.Module{
							mock.Module{
								ModuleMeta: mock.ModuleMeta{Name: "ghi"},
								CommandsReturn: []plugin.Command{
									mock.Command{CommandMeta: mock.CommandMeta{Name: "jkl"}},
								},
							},
						}, nil
					},
					defaults: plugin.Defaults{BotPermissions: discord.PermissionConnect},
				},
			},
		}

		expect := []*plugin.RegisteredModule{
			mock.GenerateRegisteredModuleWithDefaults(
				plugin.BuiltInProvider,
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "abc"},
					CommandsReturn: []plugin.Command{
						mock.Command{CommandMeta: mock.CommandMeta{Name: "def"}},
					},
				},
				plugin.Defaults{BotPermissions: discord.PermissionSendMessages}),
			mock.GenerateRegisteredModuleWithDefaults(
				"another",
				mock.Module{
					ModuleMeta: mock.ModuleMeta{Name: "ghi"},
					CommandsReturn: []plugin.Command{
						mock.Command{CommandMeta: mock.CommandMeta{Name: "jkl"}},
					},
				},
				plugin.Defaults{BotPermissions: discord.PermissionConnect}),
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
						mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
					},
					Defaults: plugin.Defaults{BotPermissions: discord.PermissionSendMessages},
				},
			},
			remProviders: []pluginProvider{
				{
					name: "another",
					provider: func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
						return nil, nil, errors.New("abc")
					},
					defaults: plugin.Defaults{BotPermissions: discord.PermissionConnect},
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
						mock.Command{CommandMeta: mock.CommandMeta{Name: "abc"}},
					},
					Defaults: plugin.Defaults{BotPermissions: discord.PermissionSendMessages},
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

	var h ctxErrorHandler = func(err error) {
		actual = err
	}

	expect := errors.New("Abort. Retry. Fail.") //nolint:golint

	h.HandleError(expect)

	assert.Equal(t, expect, actual)
}
