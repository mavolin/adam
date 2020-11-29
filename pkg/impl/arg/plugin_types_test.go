package arg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestCommand_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := &Context{
			Raw: "abc",
			Context: &plugin.Context{
				Provider: mock.PluginProvider{
					PluginRepositoriesReturn: []plugin.Repository{
						{
							ProviderName: "built_in",
							Commands: []plugin.Command{
								mock.Command{
									CommandMeta: mock.CommandMeta{Name: "abc"},
								},
							},
						},
					},
				},
			},
		}

		expect := ctx.FindCommand(ctx.Raw)
		assert.NotNil(t, expect)

		actual, err := Command.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("unknown command", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						PluginRepositoriesReturn: []plugin.Repository{
							{ProviderName: "built_in"},
						},
					},
				},
			}

			expect := newArgParsingErr(commandNotFound, ctx, map[string]interface{}{
				"invoke": ctx.Raw,
			})

			_, actual := Command.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown command - some commands unavailable", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						UnavailablePluginProvidersReturn: []plugin.UnavailablePluginProvider{
							{
								Name:  "abc",
								Error: errors.New("oh no, it didn't work"),
							},
						},
					},
				},
			}

			expect := newArgParsingErr(commandNotFoundProvidersUnavailable, ctx, map[string]interface{}{
				"invoke": ctx.Raw,
			})

			_, actual := Command.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}

func TestModule_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := &Context{
			Raw: "abc",
			Context: &plugin.Context{
				Provider: mock.PluginProvider{
					PluginRepositoriesReturn: []plugin.Repository{
						{
							ProviderName: "built_in",
							Modules: []plugin.Module{
								mock.Module{
									ModuleMeta: mock.ModuleMeta{Name: "abc"},
								},
							},
						},
					},
				},
			},
		}

		expect := ctx.FindModule(ctx.Raw)
		require.NotNil(t, expect)

		actual, err := Module.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("unknown command", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						PluginRepositoriesReturn: []plugin.Repository{
							{ProviderName: "built_in"},
						},
					},
				},
			}

			expect := newArgParsingErr(moduleNotFound, ctx, map[string]interface{}{
				"invoke": ctx.Raw,
			})

			_, actual := Module.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown command - some commands unavailable", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						UnavailablePluginProvidersReturn: []plugin.UnavailablePluginProvider{
							{
								Name:  "abc",
								Error: errors.New("oh no, it didn't work"),
							},
						},
					},
				},
			}

			expect := newArgParsingErr(moduleNotFoundProvidersUnavailable, ctx, map[string]interface{}{
				"invoke": ctx.Raw,
			})

			_, actual := Module.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}

func TestPlugin_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("command", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						PluginRepositoriesReturn: []plugin.Repository{
							{
								ProviderName: "built_in",
								Commands: []plugin.Command{
									mock.Command{
										CommandMeta: mock.CommandMeta{Name: "abc"},
									},
								},
							},
						},
					},
				},
			}

			expect := ctx.FindCommand(ctx.Raw)
			assert.NotNil(t, expect)

			actual, err := Command.Parse(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})

		t.Run("module", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						PluginRepositoriesReturn: []plugin.Repository{
							{
								ProviderName: "built_in",
								Modules: []plugin.Module{
									mock.Module{
										ModuleMeta: mock.ModuleMeta{Name: "abc"},
									},
								},
							},
						},
					},
				},
			}

			expect := ctx.FindModule(ctx.Raw)
			require.NotNil(t, expect)

			actual, err := Module.Parse(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("unknown plugin", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						PluginRepositoriesReturn: []plugin.Repository{
							{ProviderName: "built_in"},
						},
					},
				},
			}

			expect := newArgParsingErr(pluginNotFound, ctx, map[string]interface{}{
				"invoke": ctx.Raw,
			})

			_, actual := Plugin.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown command - some commands unavailable", func(t *testing.T) {
			ctx := &Context{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.PluginProvider{
						UnavailablePluginProvidersReturn: []plugin.UnavailablePluginProvider{
							{
								Name:  "abc",
								Error: errors.New("oh no, it didn't work"),
							},
						},
					},
				},
			}

			expect := newArgParsingErr(pluginNotFoundProvidersUnavailable, ctx, map[string]interface{}{
				"invoke": ctx.Raw,
			})

			_, actual := Plugin.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}
