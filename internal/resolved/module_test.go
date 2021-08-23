package resolved

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestModule_ShortDescription(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := "abc"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockplugin.Module{ShortDescription: expect}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("fallback", func(t *testing.T) {
		expect := "def"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockplugin.Module{ShortDescription: ""}}},
				{Modules: []plugin.Module{mockplugin.Module{ShortDescription: expect}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("none", func(t *testing.T) {
		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockplugin.Module{ShortDescription: ""}}},
				{Modules: []plugin.Module{mockplugin.Module{ShortDescription: ""}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Empty(t, actual)
	})
}

func TestModule_LongDescription(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := "abc"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockplugin.Module{LongDescription: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("fallback", func(t *testing.T) {
		expect := "def"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockplugin.Module{}}},
				{Modules: []plugin.Module{mockplugin.Module{LongDescription: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("short description", func(t *testing.T) {
		expect := "abc"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockplugin.Module{}}},
				{Modules: []plugin.Module{mockplugin.Module{ShortDescription: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("none", func(t *testing.T) {
		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockplugin.Module{}}},
				{Modules: []plugin.Module{mockplugin.Module{LongDescription: ""}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Empty(t, actual)
	})
}

func TestModule_FindCommand(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		expect := &Command{source: mockplugin.Command{Name: "def"}}

		rmod := &Module{
			commands: []plugin.ResolvedCommand{
				&Command{source: mockplugin.Command{Name: "abc"}},
				expect,
				&Command{source: mockplugin.Command{Name: "ghi"}},
			},
		}

		actual := rmod.FindCommand(expect.Name())
		assert.Equal(t, expect, actual)
	})

	t.Run("alias", func(t *testing.T) {
		expect := &Command{
			source:  mockplugin.Command{Name: "def", Aliases: []string{"mno"}},
			aliases: []string{"mno"},
		}

		rmod := &Module{
			commands: []plugin.ResolvedCommand{
				&Command{
					source: mockplugin.Command{Name: "abc", Aliases: []string{"jkl"}},
				},
				expect,
				&Command{source: mockplugin.Command{Name: "ghi"}},
			},
		}

		actual := rmod.FindCommand(expect.Aliases()[0])
		assert.Equal(t, expect, actual)
	})

	t.Run("not found", func(t *testing.T) {
		actual := new(Module).FindCommand("abc")
		assert.Nil(t, actual)
	})
}

func TestResolvedModule_FindModule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := &Module{
			sources: []plugin.SourceModule{
				{
					SourceName: plugin.BuiltInSource,
					Modules: []plugin.Module{
						mockplugin.Module{Name: "def"},
					},
				},
			},
		}

		rmod := &Module{
			modules: []plugin.ResolvedModule{
				&Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockplugin.Module{Name: "abc"},
							},
						},
					},
				},
				expect,
				&Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockplugin.Module{Name: "ghi"},
							},
						},
					},
				},
			},
		}

		actual := rmod.FindModule(expect.Name())
		assert.Equal(t, expect, actual)
	})

	t.Run("not found", func(t *testing.T) {
		actual := new(Module).FindModule("abc")
		assert.Nil(t, actual)
	})
}
