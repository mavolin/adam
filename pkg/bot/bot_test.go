package bot

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
)

func TestBot_AddPluginProvider(t *testing.T) {
	t.Run("name built_in", func(t *testing.T) {
		b, err := New(Options{Token: "abc"})
		require.NoError(t, err)

		assert.Panics(t, func() {
			pfunc := func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return nil, nil, nil
			}

			b.AddPluginProvider(plugin.BuiltInProvider, pfunc)
		})
	})

	t.Run("nil", func(t *testing.T) {
		b, err := New(Options{Token: "abc"})
		require.NoError(t, err)

		b.AddPluginProvider("abc", nil)
		assert.Len(t, b.pluginProviders, 0)
	})

	t.Run("replace", func(t *testing.T) {
		b, err := New(Options{Token: "abc"})
		require.NoError(t, err)

		p := mockPluginProvider(nil, nil, nil)

		b.AddPluginProvider("abc", p)
		b.AddPluginProvider("def", p)

		assert.Len(t, b.pluginProviders, 2)

		var called bool

		b.AddPluginProvider("abc",
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				called = true
				return nil, nil, nil
			})

		assert.Len(t, b.pluginProviders, 2)
		assert.Equal(t, b.pluginProviders[0].name, "def")
		assert.Equal(t, b.pluginProviders[1].name, "abc")

		_, _, _ = b.pluginProviders[1].provider(nil, nil)
		assert.True(t, called, "Bot.AddPluginProvider did not replace abc")
	})

	t.Run("success", func(t *testing.T) {
		b, err := New(Options{Token: "abc"})
		require.NoError(t, err)

		p := mockPluginProvider(nil, nil, nil)

		b.AddPluginProvider("abc", p)
		assert.Len(t, b.pluginProviders, 1)
	})
}
