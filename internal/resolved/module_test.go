package resolved

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestModule_ShortDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		shortDescriptions []string

		expect string
	}{
		{
			name:              "first hit",
			shortDescriptions: []string{"abc", "def"},
			expect:            "abc",
		},
		{
			name:              "fallback",
			shortDescriptions: []string{"", "abc"},
			expect:            "abc",
		},
		{
			name:              "none",
			shortDescriptions: []string{"", ""},
			expect:            "",
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			rmod := &Module{sources: make([]plugin.SourceModule, len(c.shortDescriptions))}

			for i, sdesc := range c.shortDescriptions {
				rmod.sources[i].Modules = []plugin.Module{
					mockplugin.Module{ShortDescription: sdesc},
				}
			}

			actual := rmod.ShortDescription(nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestModule_LongDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		modules []plugin.Module

		expect string
	}{
		{
			name: "first hit",
			modules: []plugin.Module{
				mockplugin.Module{LongDescription: "abc"},
				mockplugin.Module{LongDescription: "def"},
			},
			expect: "abc",
		},
		{
			name: "second hit",
			modules: []plugin.Module{
				mockplugin.Module{},
				mockplugin.Module{LongDescription: "abc"},
			},
			expect: "abc",
		},
		{
			name: "fallback to short description",
			modules: []plugin.Module{
				mockplugin.Module{},
				mockplugin.Module{ShortDescription: "abc"},
			},
			expect: "abc",
		},
		{
			name: "empty",
			modules: []plugin.Module{
				mockplugin.Module{},
				mockplugin.Module{},
			},
			expect: "",
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			rmod := &Module{sources: make([]plugin.SourceModule, len(c.modules))}

			for i, smod := range c.modules {
				rmod.sources[i].Modules = []plugin.Module{smod}
			}

			actual := rmod.LongDescription(nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestModule_FindCommand(t *testing.T) {
	t.Parallel()

	smod := mockplugin.Module{
		Name: "abc",
		Commands: []plugin.Command{
			mockplugin.Command{Name: "def", Aliases: []string{"ghi", "jkl"}},
			mockplugin.Command{Name: "mno"},
			mockplugin.Command{Name: "pqr", Aliases: []string{"stu"}},
		},
	}

	resolver := NewPluginResolver(nil)
	resolver.AddBuiltInModule(smod)

	rmod := resolver.NewProvider(event.NewBase(), &discord.Message{}).Modules()[0]

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, expect := range rmod.Commands() {
			actual := rmod.FindCommand(expect.Name())
			require.NotNilf(t, actual,
				"expected query %s to yield command %s, but found nil", expect.Name(), expect.ID())
			assert.Samef(t, expect, actual,
				"expected query %s to yield command %s, but found %s", expect.Name(), expect.ID(), actual.ID())

			for _, alias := range expect.Aliases() {
				actual = rmod.FindCommand(alias)
				require.NotNilf(t, actual,
					"expected query %s to yield command %s, but found nil", expect.Name(), expect.ID())
				assert.Samef(t, expect, actual,
					"expected query %s to yield command %s, but found %s", expect.Name(), expect.ID(), actual.ID())
			}
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		actual := rmod.FindCommand("cba")
		assert.Nil(t, actual)
	})
}

func TestResolvedModule_FindModule(t *testing.T) {
	t.Parallel()

	smod := mockplugin.Module{
		Name: "abc",
		Modules: []plugin.Module{
			mockplugin.Module{Name: "def"},
			mockplugin.Module{Name: "ghi"},
			mockplugin.Module{Name: "jkl"},
		},
	}

	resolver := NewPluginResolver(nil)
	resolver.AddBuiltInModule(smod)

	rmod := resolver.NewProvider(event.NewBase(), &discord.Message{}).Modules()[0]

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, expect := range rmod.Modules() {
			actual := rmod.FindModule(expect.Name())
			require.NotNilf(t, actual,
				"expected query %s to yield module %s, but found nil", expect.Name(), expect.ID())
			assert.Samef(t, expect, actual,
				"expected query %s to yield module %s, but found %s", expect.Name(), expect.ID(), actual.ID())
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		actual := rmod.FindModule("cba")
		assert.Nil(t, actual)
	})
}
