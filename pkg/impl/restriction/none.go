package restriction

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// None is a no-op plugin.RestrictionFunc.
// It can be used to prevent inheritance of the RestrictionFunc of the parent.
var None plugin.RestrictionFunc = func(*state.State, *plugin.Context) error { return nil }
