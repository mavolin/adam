package help

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

func mockHideFunc(ret HiddenLevel) HideFunc {
	return func(*plugin.RegisteredCommand, *state.State, *plugin.Context) HiddenLevel {
		return ret
	}
}
