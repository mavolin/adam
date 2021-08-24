package bot

import (
	"reflect"
	"testing"

	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestMiddlewareManager_AddMiddleware(t *testing.T) {
	successCases := []interface{}{
		func(*state.State, interface{}) {},
		func(*state.State, interface{}) error { return nil },
		func(*state.State, *event.Base) {},
		func(*state.State, *event.Base) error { return nil },
		func(*state.State, *event.MessageCreate) {},
		func(*state.State, *event.MessageCreate) error { return nil },
		func(*state.State, *event.MessageUpdate) {},
		func(*state.State, *event.MessageUpdate) error { return nil },
		func(next CommandFunc) CommandFunc {
			return func(*state.State, *plugin.Context) error { return nil }
		},
		Middleware(func(next CommandFunc) CommandFunc {
			return func(*state.State, *plugin.Context) error { return nil }
		}),
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			funcType := reflect.TypeOf(c)

			t.Run(funcType.String(), func(t *testing.T) {
				var m MiddlewareManager

				err := m.TryAddMiddleware(c)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		var m MiddlewareManager

		err := m.TryAddMiddleware("invalid")
		assert.True(t, errors.Is(err, ErrMiddleware))
	})
}
