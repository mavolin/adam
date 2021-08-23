package resolved

import (
	"testing"

	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestCommand_LongDescription(t *testing.T) {
	t.Run("long description", func(t *testing.T) {
		expect := "abc"

		rcmd := &Command{
			source: mockplugin.Command{
				LongDescription: expect,
			},
		}

		actual := rcmd.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("short description", func(t *testing.T) {
		expect := "abc"

		rcmd := &Command{
			source: mockplugin.Command{
				ShortDescription: expect,
			},
		}

		actual := rcmd.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})
}

func TestResolvedCommand_IsRestricted(t *testing.T) {
	t.Run("regular error", func(t *testing.T) {
		expect := errors.New("abc")

		rcmd := &Command{
			source: mockplugin.Command{
				Restrictions: func(*state.State, *plugin.Context) error {
					return expect
				},
			},
		}

		actual := rcmd.IsRestricted(nil, nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("plugin.RestrictionErrorWrapper", func(t *testing.T) {
		expect := errors.New("abc")

		rcmd := &Command{
			source: mockplugin.Command{
				Restrictions: func(*state.State, *plugin.Context) error {
					return &mockplugin.RestrictionErrorWrapper{Return: expect}
				},
			},
		}

		actual := rcmd.IsRestricted(nil, nil)
		assert.Equal(t, expect, actual)
	})
}
