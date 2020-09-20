package throttling

import "github.com/mavolin/adam/pkg/plugin"

// None is a no-op plugin.Throttler.
// It can be used to override the throttler of a parent.
var None plugin.Throttler = new(none)

type none struct{}

func (n none) Check(*plugin.Context) (func(), error) { return func() {}, nil }
