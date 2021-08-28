package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
)

func TestLocalizer_Clone(t *testing.T) {
	t.Parallel()

	a := NewLocalizer(t).
		On("abc", "def")

	b := a.Clone(t)

	assert.Equal(t, a, b)
	assert.NotSame(t, a, b)

	a.On("ghi", "jkl")

	assert.NotEqual(t, a, b)
}

func TestLocalizer_Build(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		l := NewLocalizerWithDefault(t, expect).
			Build()

		actual, err := l.LocalizeTerm("def")
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("expected i18n", func(t *testing.T) {
		t.Parallel()

		t.Run("error on", func(t *testing.T) {
			t.Parallel()

			var term i18n.Term = "abc"

			l := NewLocalizer(t).
				ErrorOn(term).
				Build()

			_, err := l.LocalizeTerm(term)
			assert.Error(t, err)
		})

		t.Run("on", func(t *testing.T) {
			t.Parallel()

			var term i18n.Term = "abc"

			expect := "def"

			l := NewLocalizer(t).
				On(term, expect).
				Build()

			actual, err := l.LocalizeTerm(term)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("unexpected i18n", func(t *testing.T) {
		t.Parallel()

		tMock := new(testing.T)

		l := NewLocalizer(tMock).
			Build()

		actual, err := l.LocalizeTerm("unknown_term")
		assert.Empty(tMock, actual)
		assert.Error(t, err)

		assert.True(t, tMock.Failed())
	})
}
