package bot

import (
	"reflect"
	"testing"

	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestMiddlewareManager_AddMiddleware(t *testing.T) {
	successCases := []interface{}{
		func(*state.State, interface{}) {},
		func(*state.State, interface{}) error { return nil },
		func(*state.State, *state.Base) {},
		func(*state.State, *state.Base) error { return nil },
		func(*state.State, *state.MessageCreateEvent) {},
		func(*state.State, *state.MessageCreateEvent) error { return nil },
		func(*state.State, *state.MessageUpdateEvent) {},
		func(*state.State, *state.MessageUpdateEvent) error { return nil },
		func(next CommandFunc) CommandFunc {
			return func(*state.State, *plugin.Context) error {
				return nil
			}
		},
		MiddlewareFunc(func(next CommandFunc) CommandFunc {
			return func(*state.State, *plugin.Context) error {
				return nil
			}
		}),
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			funcType := reflect.TypeOf(c)

			t.Run(funcType.String(), func(t *testing.T) {
				var m MiddlewareManager

				err := m.AddMiddleware(c)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		var m MiddlewareManager

		err := m.AddMiddleware("invalid")
		assert.True(t, errors.Is(err, ErrNotAMiddleware))
	})
}
