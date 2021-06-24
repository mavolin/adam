package arg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestCommand_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := &plugin.ParseContext{
			Raw: "abc",
			Context: &plugin.Context{
				Provider: mock.NewPluginProvider([]plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Commands: []plugin.Command{
							mock.Command{CommandMeta: command.Meta{Name: "abc"}},
						},
					},
				}, nil),
			},
		}

		expect := ctx.FindCommand(ctx.Raw)
		assert.NotNil(t, expect)

		actual, err := Command.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("unknown commandType", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.NewPluginProvider([]plugin.Source{{Name: plugin.BuiltInSource}}, nil),
				},
			}

			expect := newArgumentError(commandNotFoundError, ctx, nil)

			_, actual := Command.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown commandType - some commands unavailable", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.NewPluginProvider(nil, []plugin.UnavailablePluginSource{
						{
							Name:  "abc",
							Error: errors.New("oh no, it didn't work"),
						},
					}),
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
				Provider: mock.NewPluginProvider([]plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{name: "abc"},
						},
					},
				}, nil),
			},
		}

		expect := ctx.FindModule(ctx.Raw)
		require.NotNil(t, expect)

		actual, err := Module.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("unknown commandType", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.NewPluginProvider([]plugin.Source{{Name: plugin.BuiltInSource}}, nil),
				},
			}

			expect := newArgumentError(moduleNotFoundError, ctx, nil)

			_, actual := Module.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown commandType - some commands unavailable", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.NewPluginProvider(nil, []plugin.UnavailablePluginSource{
						{
							Name:  "abc",
							Error: errors.New("oh no, it didn't work"),
						},
					}),
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
		t.Run("commandType", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.NewPluginProvider([]plugin.Source{
						{
							Name: plugin.BuiltInSource,
							Commands: []plugin.Command{
								mock.Command{CommandMeta: command.Meta{Name: "abc"}},
							},
						},
					}, nil),
				},
			}

			expect := ctx.FindCommand(ctx.Raw)
			assert.NotNil(t, expect)

			actual, err := Command.Parse(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})

		t.Run("moduleType", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.NewPluginProvider([]plugin.Source{
						{
							Name: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{name: "abc"},
							},
						},
					}, nil),
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
					Provider: mock.NewPluginProvider([]plugin.Source{{Name: plugin.BuiltInSource}}, nil),
				},
			}

			expect := newArgumentError(pluginNotFoundError, ctx, nil)

			_, actual := Plugin.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("unknown commandType - some commands unavailable", func(t *testing.T) {
			ctx := &plugin.ParseContext{
				Raw: "abc",
				Context: &plugin.Context{
					Provider: mock.NewPluginProvider(nil, []plugin.UnavailablePluginSource{
						{
							Name:  "abc",
							Error: errors.New("oh no, that didn't work"),
						},
					}),
				},
			}

			expect := newArgumentError(pluginNotFoundErrorProvidersUnavailable, ctx, nil)

			_, actual := Plugin.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}
