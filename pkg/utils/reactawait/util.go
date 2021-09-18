package reactawait

import (
	"context"

	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"
)

func invokeReactionAddMiddlewares(s *state.State, e *event.MessageReactionAdd, middlewares []interface{}) error {
	for _, m := range middlewares {
		switch m := m.(type) {
		case func(*state.State, interface{}):
			m(s, e)
		case func(*state.State, interface{}) error:
			if err := m(s, e); err != nil {
				return err
			}
		case func(*state.State, *event.Base):
			m(s, e.Base)
		case func(*state.State, *event.Base) error:
			if err := m(s, e.Base); err != nil {
				return err
			}
		case func(*state.State, *event.MessageReactionAdd):
			m(s, e)
		case func(*state.State, *event.MessageReactionAdd) error:
			if err := m(s, e); err != nil {
				return err
			}
		}
	}

	return nil
}

// sendResult blocks until it can send a result or the passed context.Context
// gets canceled.
func sendResult(ctx context.Context, result chan<- interface{}, val interface{}) {
	select {
	case <-ctx.Done():
	case result <- val:
	}
}
