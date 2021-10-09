package bot

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/internal/resolved"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestBot_AddPluginSource(t *testing.T) {
	t.Parallel()

	t.Run("name built_in", func(t *testing.T) {
		t.Parallel()

		b := &Bot{pluginResolver: resolved.NewPluginResolver(nil)}

		assert.Panics(t, func() {
			source := func(*event.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return nil, nil, nil
			}

			b.AddPluginSource(plugin.BuiltInSource, source)
		})
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		b := &Bot{pluginResolver: resolved.NewPluginResolver(nil)}

		b.AddPluginSource("abc", nil)
		assert.Len(t, b.pluginResolver.CustomSources, 0)
	})

	t.Run("replace", func(t *testing.T) {
		t.Parallel()

		b := &Bot{pluginResolver: resolved.NewPluginResolver(nil)}

		p := func(*event.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
			return nil, nil, nil
		}

		b.AddPluginSource("abc", p)
		b.AddPluginSource("def", p)

		assert.Len(t, b.pluginResolver.CustomSources, 2)

		var called bool

		b.AddPluginSource("abc",
			func(*event.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				called = true
				return nil, nil, nil
			})

		assert.Len(t, b.pluginResolver.CustomSources, 2)
		assert.Equal(t, b.pluginResolver.CustomSources[0].Name, "def")
		assert.Equal(t, b.pluginResolver.CustomSources[1].Name, "abc")

		_, _, _ = b.pluginResolver.CustomSources[1].Func(nil, nil)
		assert.True(t, called, "Bot.AddPluginSource did not replace abc")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		b := &Bot{pluginResolver: resolved.NewPluginResolver(nil)}

		p := func(*event.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
			return nil, nil, nil
		}

		b.AddPluginSource("abc", p)
		assert.Len(t, b.pluginResolver.CustomSources, 1)
	})
}
