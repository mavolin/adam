package embedutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
)

func Test_mockLocalizer_build(t *testing.T) {
	t.Run("expected localization", func(t *testing.T) {
		t.Run("on", func(t *testing.T) {
			var term localization.Term = "abc"

			expect := "def"

			l := newMockedLocalizer(t).
				on(term, expect).
				build()

			actual, err := l.LocalizeTerm(term)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("unexpected localization", func(t *testing.T) {
		var term localization.Term = "unknown_term"

		tMock := new(testing.T)

		l := newMockedLocalizer(tMock).
			build()

		actualTerm, err := l.LocalizeTerm(term)
		assert.Equal(t, string(term), actualTerm)
		assert.Error(t, err)

		assert.True(t, tMock.Failed())
	})
}
