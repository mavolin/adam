package messageutil

import (
	"context"

	"github.com/mavolin/disstate/v3/pkg/state"
)

//nolint:dupl
func invokeMessageMiddlewares(s *state.State, e *state.MessageCreateEvent, middlewares []interface{}) error {
	for _, m := range middlewares {
		switch m := m.(type) {
		case func(*state.State, interface{}):
			m(s, e)
		case func(*state.State, interface{}) error:
			if err := m(s, e); err != nil {
				return err
			}
		case func(*state.State, *state.Base):
			m(s, e.Base)
		case func(*state.State, *state.Base) error:
			if err := m(s, e.Base); err != nil {
				return err
			}
		case func(*state.State, *state.MessageCreateEvent):
			m(s, e)
		case func(*state.State, *state.MessageCreateEvent) error:
			if err := m(s, e); err != nil {
				return err
			}
		}
	}

	return nil
}

//nolint:dupl
func invokeReactionAddMiddlewares(s *state.State, e *state.MessageReactionAddEvent, middlewares []interface{}) error {
	for _, m := range middlewares {
		switch m := m.(type) {
		case func(*state.State, interface{}):
			m(s, e)
		case func(*state.State, interface{}) error:
			if err := m(s, e); err != nil {
				return err
			}
		case func(*state.State, *state.Base):
			m(s, e.Base)
		case func(*state.State, *state.Base) error:
			if err := m(s, e.Base); err != nil {
				return err
			}
		case func(*state.State, *state.MessageReactionAddEvent):
			m(s, e)
		case func(*state.State, *state.MessageReactionAddEvent) error:
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
