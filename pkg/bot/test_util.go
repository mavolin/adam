package bot

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

func mockPluginProvider(cmds []plugin.Command, mods []plugin.Module, err error) PluginProvider {
	return func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
		return cmds, mods, err
	}
}
