package help

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// checkHideFuncs checks the passed HideFuncs and returns the highest
// HiddenLevel found.
func checkHideFuncs(cmd *plugin.RegisteredCommand, s *state.State, ctx *plugin.Context, f ...HideFunc) HiddenLevel {
	var lvl HiddenLevel

	for _, f := range f {
		lvl2 := f(cmd, s, ctx)
		if lvl2 >= Hide {
			return Hide
		} else if lvl2 > lvl {
			lvl = lvl2
		}
	}

	return lvl
}
