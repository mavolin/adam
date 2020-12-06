package throttler

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// Conditional is a plugin.Throttler that uses different cooldowns based on
	// restrictions.
	// It will attempt to find the firstmost Condition, where a call to
	// Condition.Restrictions returns nil.
	// If there is no such Condition, the default throttler will be used.
	//
	// If the throttler of a condition or the default throttler is nil,
	// the encountering user will be regarded as exempt from throttling.
	Conditional struct {
		// Conditions contains the conditional throttlers.
		Conditions []Condition
		// Default is the fallback throttler.
		// If it is nil, users encountering this won't be throttled.
		Default plugin.Throttler
	}

	// Condition is a single conditional throttler.
	Condition struct {
		// Restrictions is the restriction a user must fulfill to use the
		// throttler in this condition.
		//
		// Restrictions is considered fulfilled, if it returns nil.
		Restrictions plugin.RestrictionFunc
		// Throttler is the plugin.Throttler used if the Restrictions are
		// fulfilled.
		Throttler plugin.Throttler
	}
)

func (c Conditional) Check(s *state.State, ctx *plugin.Context) (func(), error) {
	for _, con := range c.Conditions {
		if con.Restrictions(s, ctx) == nil {
			if con.Throttler != nil {
				return con.Throttler.Check(s, ctx)
			}

			return func() {}, nil
		}
	}

	if c.Default != nil {
		return c.Default.Check(s, ctx)
	}

	return func() {}, nil
}
