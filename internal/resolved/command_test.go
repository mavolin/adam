package resolved

import (
	"testing"

	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestCommand_LongDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		longDescription  string
		shortDescription string

		expect string
	}{
		{
			name:             "long description",
			longDescription:  "abc",
			shortDescription: "def",
			expect:           "abc",
		},
		{name: "short description", shortDescription: "abc", expect: "abc"},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			rcmd := &Command{
				source: mockplugin.Command{
					ShortDescription: c.shortDescription,
					LongDescription:  c.longDescription,
				},
			}

			actual := rcmd.LongDescription(nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestResolvedCommand_IsRestricted(t *testing.T) {
	t.Parallel()

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
}
