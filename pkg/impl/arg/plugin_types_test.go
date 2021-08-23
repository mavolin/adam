package arg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestCommand_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := &plugin.ParseContext{
			Raw: "abc",
			Context: &plugin.Context{
				Provider: &mock.PluginProvider{
					Sources: []plugin.Source{
						{
							Name: plugin.BuiltInSource,
							Commands: []plugin.Command{
								mock.Command{Name: "abc"},
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
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						Sources: []plugin.Source{{Name: plugin.BuiltInSource}},
					},
				},
			}

			expect := newArgumentError(commandNotFoundError, ctx, nil)

			_, actual := Command.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown command - some commands unavailable", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						UnavailableSources: []plugin.UnavailableSource{
							{
								Name:  "abc",
								Error: errors.New("oh no, it didn't work"),
							},
						},
					},
				},
			}

			expect := newArgumentError(commandNotFoundErrorProvidersUnavailable, ctx, nil)

			_, actual := Command.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}

func TestModule_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := &plugin.ParseContext{
			Raw: "abc",
			Context: &plugin.Context{
				Provider: &mock.PluginProvider{
					Sources: []plugin.Source{
						{
							Name: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockplugin.Module{Name: "abc"},
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
		t.Run("unknown module", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						Sources: []plugin.Source{{Name: plugin.BuiltInSource}},
					},
				},
			}

			expect := newArgumentError(moduleNotFoundError, ctx, nil)

			_, actual := Module.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown module - some commands unavailable", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						UnavailableSources: []plugin.UnavailableSource{
							{
								Name:  "abc",
								Error: errors.New("oh no, it didn't work"),
							},
						},
					},
				},
			}

			expect := newArgumentError(moduleNotFoundErrorProvidersUnavailable, ctx, nil)

			_, actual := Module.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}

func TestPlugin_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("command", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						Sources: []plugin.Source{
							{
								Name: plugin.BuiltInSource,
								Commands: []plugin.Command{
									mock.Command{Name: "abc"},
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
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						Sources: []plugin.Source{
							{
								Name: plugin.BuiltInSource,
								Modules: []plugin.Module{
									mockplugin.Module{Name: "abc"},
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
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						Sources: []plugin.Source{{Name: plugin.BuiltInSource}},
					},
				},
			}

			expect := newArgumentError(pluginNotFoundError, ctx, nil)

			_, actual := Plugin.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown plugin - some commands unavailable", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: &mock.PluginProvider{
						UnavailableSources: []plugin.UnavailableSource{
							{
								Name:  "abc",
								Error: errors.New("oh no, that didn't work"),
							},
						},
					},
				},
			}

			expect := newArgumentError(pluginNotFoundErrorProvidersUnavailable, ctx, nil)

			_, actual := Plugin.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}
