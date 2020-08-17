package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
)

func TestNewNoOpLocalizer(t *testing.T) {
	l := NewNoOpLocalizer()

	_, err := l.LocalizeTerm("abc")
	assert.Error(t, err)
}

func TestLocalizer_Clone(t *testing.T) {
	l1 := NewLocalizer().
		On("abc", "def")

	l2 := l1.Clone()

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

	t.Run("expected localization", func(t *testing.T) {
		t.Run("error on", func(t *testing.T) {
			var term localization.Term = "abc"

			l := NewLocalizer().
				ErrorOn(term).
				Build()

			_, err := l.LocalizeTerm(term)
			assert.Error(t, err)
		})

		t.Run("on", func(t *testing.T) {
			var term localization.Term = "abc"

			expect := "def"

			l := NewLocalizer().
				On(term, expect).
				Build()

			actual, err := l.LocalizeTerm(term)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("unexpected localization", func(t *testing.T) {
		l := NewLocalizer().
			Build()

		assert.Panics(t, func() {
			l.LocalizeTerm("abc")
		})
	})
}
