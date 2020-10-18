package i18nutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
)

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
