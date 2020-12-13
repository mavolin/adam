package mock

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Throttler is the mocked version of a plugin.Throttler.
type Throttler struct {
	checkReturn  error
	CancelCalled bool
}

var _ plugin.Throttler = new(Throttler)

// NewThrottler creates a new mocked Throttler with the given return value
// for check.
func NewThrottler(checkReturn error) *Throttler {
	return &Throttler{checkReturn: checkReturn}
}

func (t *Throttler) Check(*state.State, *plugin.Context) (func(), error) {
	return func() {
		t.CancelCalled = true
	}, t.checkReturn
}
