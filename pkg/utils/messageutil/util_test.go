package messageutil

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mavolin/adam/pkg/errors"
)

func Test_invokeMiddlewares(t *testing.T) {
	t.Parallel()

	_, s := state.NewMocker(t)

	e := &event.MessageCreate{
		Base: event.NewBase(),
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: discord.Message{Content: "abc"},
		},
	}
	e.Base.Set("abc", "def")

	testErr := errors.New("abc")

	t.Run("state, interface", func(t *testing.T) {
		t.Parallel()

		m := mock.Mock{}

		middleware := func(s *state.State, e interface{}) {
			m.Called(s, e)
		}

		m.On("1", s, e)

		err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
		assert.NoError(t, err)

		m.AssertExpectations(t)
	})

	t.Run("state, interface returns error", func(t *testing.T) {
		t.Parallel()

		t.Run("no error", func(t *testing.T) {
			t.Parallel()

			m := mock.Mock{}

			middleware := func(s *state.State, e interface{}) error {
				return m.
					Called(s, e).
					Error(0)
			}

			m.
				On("1", s, e).
				Return(nil)

			err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
			assert.NoError(t, err)

			m.AssertExpectations(t)
		})

		t.Run("error", func(t *testing.T) {
			t.Parallel()

			m := mock.Mock{}

			middleware := func(s *state.State, b interface{}) error {
				return m.
					Called(s, b).
					Error(0)
			}

			m.
				On("1", s, e).
				Return(testErr)

			err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
			assert.Equal(t, testErr, err)

			m.AssertExpectations(t)
		})
	})

	t.Run("state, base", func(t *testing.T) {
		t.Parallel()

		m := mock.Mock{}

		middleware := func(s *state.State, e *event.Base) {
			m.Called(s, e)
		}

		m.On("1", s, e.Base)

		err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
		assert.NoError(t, err)

		m.AssertExpectations(t)
	})

	t.Run("state, base returns error", func(t *testing.T) {
		t.Parallel()

		t.Run("no error", func(t *testing.T) {
			t.Parallel()

			m := mock.Mock{}

			middleware := func(s *state.State, b *event.Base) error {
				return m.
					Called(s, b).
					Error(0)
			}

			m.
				On("1", s, e.Base).
				Return(nil)

			err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
			assert.NoError(t, err)

			m.AssertExpectations(t)
		})

		t.Run("error", func(t *testing.T) {
			t.Parallel()

			m := mock.Mock{}

			middleware := func(s *state.State, b *event.Base) error {
				return m.
					Called(s, b).
					Error(0)
			}

			m.
				On("1", s, e.Base).
				Return(testErr)

			err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
			assert.Equal(t, testErr, err)

			m.AssertExpectations(t)
		})
	})

	t.Run("state, message create event", func(t *testing.T) {
		t.Parallel()

		m := mock.Mock{}

		middleware := func(s *state.State, e *event.MessageCreate) {
			m.Called(s, e)
		}

		m.On("1", s, e)

		err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
		assert.NoError(t, err)

		m.AssertExpectations(t)
	})

	t.Run("state, message create event returns error", func(t *testing.T) {
		t.Parallel()

		t.Run("no error", func(t *testing.T) {
			t.Parallel()

			m := mock.Mock{}

			middleware := func(s *state.State, e *event.MessageCreate) error {
				return m.
					Called(s, e).
					Error(0)
			}

			m.
				On("1", s, e).
				Return(nil)

			err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
			assert.NoError(t, err)

			m.AssertExpectations(t)
		})

		t.Run("error", func(t *testing.T) {
			t.Parallel()

			m := mock.Mock{}

			middleware := func(s *state.State, e *event.MessageCreate) error {
				return m.
					Called(s, e).
					Error(0)
			}

			m.On("1", s, e).Return(testErr)

			err := invokeMessageMiddlewares(s, e, []interface{}{middleware})
			assert.Equal(t, testErr, err)

			m.AssertExpectations(t)
		})
	})
}
