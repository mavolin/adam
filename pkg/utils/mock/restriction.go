package mock

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// RestrictionFunc creates a restriction func that returns the passed error.
func RestrictionFunc(ret error) plugin.RestrictionFunc {
	return func(*state.State, *plugin.Context) error { return ret }
}
