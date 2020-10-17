package i18nutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
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
			expect: false,
		},
		{
			name:   "string",
			text:   NewText("abc"),
			expect: true,
		},
		{
			name:   "config",
			text:   NewTextl(i18n.NewTermConfig("abc")),
			expect: true,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.text.IsValid()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestText_Get(t *testing.T) {
	t.Run("static", func(t *testing.T) {
		expect := "abc"

		text := NewText(expect)

		actual, err := text.Get(nil)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("localized", func(t *testing.T) {
		expect := "abc"

		var term i18n.Term = "def"

		text := NewTextl(term.AsConfig())

		actual, err := text.Get(newMockedLocalizer(t).
			on(term, expect).
			build())
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}
