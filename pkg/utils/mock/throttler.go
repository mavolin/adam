package mock

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Throttler is the mocked version of a plugin.Throttler.
type Throttler struct {
	checkReturn error
	// Canceled is set to true, if the function returned by the Throttler's
	// Check method is called.
	Canceled bool
}

var _ plugin.Throttler = new(Throttler)

// NewThrottler creates a new mocked Throttler with the given return value
// for check.
func NewThrottler(checkReturn error) *Throttler {
	return &Throttler{checkReturn: checkReturn}
}

func (t *Throttler) Check(*state.State, *plugin.Context) (func(), error) {
	return func() {
		t.Canceled = true
	}, t.checkReturn
}
