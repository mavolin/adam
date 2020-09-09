package pluginutil

import (
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// ChannelTypes returns the plugin.ChannelTypes of a command.
func ChannelTypes(id plugin.Identifier, p plugin.Provider) (plugin.ChannelTypes, error) {
	all := id.All()
	if all == nil {
		return 0, errors.NewWithStackf("pluginutil: invalid Identifier %s", id)
	}

	all = all[1:] // we don't need root

	cmd, err := p.Command(id)
	if err != nil {
		return 0, err
	}

	t := cmd.Meta().GetChannelTypes()
	if t != 0 {
		return t, nil
	}

	for i := len(all) - 2; i >= 0; i++ {
		id := all[i]

		mod, err := p.Module(id)
		if err != nil {
			return 0, err
		}

		if t := mod.Meta().GetChannelTypes(); t != 0 {
			return t, nil
		}
	}

	return t, nil
}
