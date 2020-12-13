package throttler

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// None is a no-op plugin.Throttler.
// It can be used to prevent inheritance from a parent.
var None plugin.Throttler = new(none)

type none struct{}

func (n none) Check(*state.State, *plugin.Context) (func(), error) { return func() {}, nil }
