package response

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
)

func invokeMiddlewares(s *state.State, e *state.MessageCreateEvent, middlewares []interface{}) error {
	for _, m := range middlewares {
		switch m := m.(type) {
		case func(*state.State, interface{}):
			m(s, e)
		case func(*state.State, interface{}) error:
			err := handleErr(m(s, e))
			if err != nil {
				return err
			}
		case func(*state.State, *state.Base):
			m(s, e.Base)
		case func(*state.State, *state.Base) error:
			err := handleErr(m(s, e.Base))
			if err != nil {
				return err
			}
		case func(*state.State, *state.MessageCreateEvent):
			m(s, e)
		case func(*state.State, *state.MessageCreateEvent) error:
			err := handleErr(m(s, e))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func handleErr(err error) error {
	if err == nil {
		return nil
	} else if err == state.Filtered {
		return errors.Abort
	}

	return err
}
