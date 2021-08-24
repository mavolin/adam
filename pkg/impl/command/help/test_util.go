package help

import (
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

func mockHideFunc(ret HiddenLevel) HideFunc {
	return func(plugin.ResolvedCommand, *state.State, *plugin.Context) HiddenLevel {
		return ret
	}
}
