package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
)

func Test_mockLocalizer_Build(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		expect := "abc"

		l := newMockedLocalizerWithDefault(expect).
			build()

		actual, err := l.LocalizeTerm("")
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("on", func(t *testing.T) {
		var term localization.Term = "abc"

		expect := "def"

		l := newMockedLocalizer().
			on(term, expect).
			build()

		actual, err := l.LocalizeTerm(term)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("unexpected localization", func(t *testing.T) {
		l := newMockedLocalizer().
			build()

		assert.Panics(t, func() {
			l.LocalizeTerm("abc")
		})
	})
}
