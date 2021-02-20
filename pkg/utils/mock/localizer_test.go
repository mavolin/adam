package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
)

func TestLocalizer_Clone(t *testing.T) {
	l1 := NewLocalizer(t).
		On("abc", "def")

	l2 := l1.Clone(t)

	assert.Equal(t, l1, l2)

	l1.On("ghi", "jkl")

	assert.NotEqual(t, l1, l2)
}

func TestLocalizer_Build(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		expect := "abc"

		l := NewLocalizerWithDefault(expect).
			Build()

		actual, err := l.LocalizeTerm("")
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("expected i18n", func(t *testing.T) {
		t.Run("error on", func(t *testing.T) {
			var term i18n.Term = "abc"

			l := NewLocalizer(t).
				ErrorOn(term).
				Build()

			_, err := l.LocalizeTerm(term)
			assert.Error(t, err)
		})

		t.Run("on", func(t *testing.T) {
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
		tMock := new(testing.T)

		l := NewLocalizer(tMock).
			Build()

		actual, err := l.LocalizeTerm("unknown_term")
		assert.Empty(tMock, actual)
		assert.Error(t, err)

		assert.True(t, tMock.Failed())
	})
}
