package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/localization"
)

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
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("expected localization", func(t *testing.T) {
		var term localization.Term = "abc"

		expect := "def"

		l := NewLocalizer().
			On(term, expect).
			Build()

		actual, err := l.LocalizeTerm(term)
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("unexpected localization", func(t *testing.T) {
		l := NewLocalizer().
			Build()

		assert.Panics(t, func() {
			l.LocalizeTerm("abc")
		})
	})
}
