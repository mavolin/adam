package messageutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mavolin/adam/pkg/errors"
)

func Test_invokeMiddlewares(t *testing.T) {
	_, s := state.NewMocker(t)

	e := &state.MessageCreateEvent{
		Base: state.NewBase(),
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: discord.Message{
				Content: "abc",
			},
		},
	}
	e.Base.Set("abc", "def")

	testErr := errors.New("abc")

	t.Run("state, interface", func(t *testing.T) {
		m := mock.Mock{}

		middleware := func(s *state.State, e interface{}) {
			m.Called(s, e)
		}

		m.On("1", s, e)

		err := invokeMiddlewares(s, e, []interface{}{middleware})
		assert.NoError(t, err)

		m.AssertExpectations(t)
	})

	t.Run("state, interface returns error", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			m := mock.Mock{}

			middleware := func(s *state.State, e interface{}) error {
				return m.
					Called(s, e).
					Error(0)
			}

			m.
				On("1", s, e).
				Return(nil)

			err := invokeMiddlewares(s, e, []interface{}{middleware})
			assert.NoError(t, err)

			m.AssertExpectations(t)
		})

		t.Run("error", func(t *testing.T) {
			m := mock.Mock{}

			middleware := func(s *state.State, b interface{}) error {
				return m.
					Called(s, b).
					Error(0)
			}

			m.
				On("1", s, e).
				Return(testErr)

			err := invokeMiddlewares(s, e, []interface{}{middleware})
			assert.Equal(t, testErr, err)

			m.AssertExpectations(t)
		})
	})

	t.Run("state, base", func(t *testing.T) {
		m := mock.Mock{}

		middleware := func(s *state.State, e *state.Base) {
			m.Called(s, e)
		}

		m.On("1", s, e.Base)

		err := invokeMiddlewares(s, e, []interface{}{middleware})
		assert.NoError(t, err)

		m.AssertExpectations(t)
	})

	t.Run("state, base returns error", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			m := mock.Mock{}

			middleware := func(s *state.State, b *state.Base) error {
				return m.
					Called(s, b).
					Error(0)
			}

			m.
				On("1", s, e.Base).
				Return(nil)

			err := invokeMiddlewares(s, e, []interface{}{middleware})
			assert.NoError(t, err)

			m.AssertExpectations(t)
		})

		t.Run("error", func(t *testing.T) {
			m := mock.Mock{}

			middleware := func(s *state.State, b *state.Base) error {
				return m.
					Called(s, b).
					Error(0)
			}

			m.
				On("1", s, e.Base).
				Return(testErr)

			err := invokeMiddlewares(s, e, []interface{}{middleware})
			assert.Equal(t, testErr, err)

			m.AssertExpectations(t)
		})
	})

	t.Run("state, message create event", func(t *testing.T) {
		m := mock.Mock{}

		middleware := func(s *state.State, e *state.MessageCreateEvent) {
			m.Called(s, e)
		}

		m.On("1", s, e)

		err := invokeMiddlewares(s, e, []interface{}{middleware})
		assert.NoError(t, err)

		m.AssertExpectations(t)
	})

	t.Run("state, message create event returns error", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			m := mock.Mock{}

			middleware := func(s *state.State, e *state.MessageCreateEvent) error {
				return m.
					Called(s, e).
					Error(0)
			}

			m.
				On("1", s, e).
				Return(nil)

			err := invokeMiddlewares(s, e, []interface{}{middleware})
			assert.NoError(t, err)

			m.AssertExpectations(t)
		})

		t.Run("error", func(t *testing.T) {
			m := mock.Mock{}

			middleware := func(s *state.State, e *state.MessageCreateEvent) error {
				return m.
					Called(s, e).
					Error(0)
			}

			m.On("1", s, e).Return(testErr)

			err := invokeMiddlewares(s, e, []interface{}{middleware})
			assert.Equal(t, testErr, err)

			m.AssertExpectations(t)
		})
	})
}

func Test_handleErr(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "filtered",
			err:    state.Filtered,
			expect: errors.Abort,
		},
		{
			name:   "not filtered",
			err:    errors.New("abc"),
			expect: errors.New("abc"),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := handleErr(c.err)
			assert.Equal(t, c.expect, actual)
		})
	}
}
