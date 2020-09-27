package locutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
)

func TestText_IsEmpty(t *testing.T) {
	testCases := []struct {
		name   string
		text   Text
		expect bool
	}{
		{
			name:   "empty",
			text:   Text{},
			expect: true,
		},
		{
			name:   "string",
			text:   NewStaticText("abc"),
			expect: false,
		},
		{
			name:   "config",
			text:   NewLocalizedText(localization.NewTermConfig("abc")),
			expect: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.text.IsEmpty()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestText_Get(t *testing.T) {
	t.Run("static", func(t *testing.T) {
		expect := "abc"

		text := NewStaticText(expect)

		actual, err := text.Get(nil)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized", func(t *testing.T) {
		expect := "abc"

		var term localization.Term = "def"

		text := NewLocalizedText(term.AsConfig())

		actual, err := text.Get(newMockedLocalizer(t).
			on(term, expect).
			build())
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}
