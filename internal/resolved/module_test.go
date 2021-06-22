package resolved

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
)

func TestModule_ShortDescription(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := "abc"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockModule{shortDesc: expect}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("fallback", func(t *testing.T) {
		expect := "def"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockModule{shortDesc: ""}}},
				{Modules: []plugin.Module{mockModule{shortDesc: expect}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("none", func(t *testing.T) {
		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockModule{shortDesc: ""}}},
				{Modules: []plugin.Module{mockModule{shortDesc: ""}}},
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
				{Modules: []plugin.Module{mockModule{longDesc: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("fallback", func(t *testing.T) {
		expect := "def"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockModule{}}},
				{Modules: []plugin.Module{mockModule{longDesc: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("short description", func(t *testing.T) {
		expect := "abc"

		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockModule{}}},
				{Modules: []plugin.Module{mockModule{shortDesc: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("none", func(t *testing.T) {
		rmod := &Module{
			sources: []plugin.SourceModule{
				{Modules: []plugin.Module{mockModule{}}},
				{Modules: []plugin.Module{mockModule{longDesc: ""}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Empty(t, actual)
	})
}

func TestModule_FindCommand(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		expect := &Command{source: mockCommand{name: "def"}}

		rmod := &Module{
			commands: []plugin.ResolvedCommand{
				&Command{source: mockCommand{name: "abc"}},
				expect,
				&Command{source: mockCommand{name: "ghi"}},
			},
		}

		actual := rmod.FindCommand(expect.Name())
		assert.Equal(t, expect, actual)
	})

	t.Run("alias", func(t *testing.T) {
		expect := &Command{
			source:  mockCommand{name: "def", aliases: []string{"mno"}},
			aliases: []string{"mno"},
		}

		rmod := &Module{
			commands: []plugin.ResolvedCommand{
				&Command{
					source: mockCommand{name: "abc", aliases: []string{"jkl"}},
				},
				expect,
				&Command{source: mockCommand{name: "ghi"}},
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
						mockModule{name: "def"},
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
								mockModule{name: "abc"},
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
								mockModule{name: "ghi"},
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
