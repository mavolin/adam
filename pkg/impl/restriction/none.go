package restriction

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// None is a no-op plugin.RestrictionFunc.
// It can be used to override the restrictions of a parent.
var None plugin.RestrictionFunc = func(*state.State, *plugin.Context) error { return nil }
