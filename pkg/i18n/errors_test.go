package i18n

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalizationError_Error(t *testing.T) {
	testCases := []struct {
		name   string
		term   Term
		expect string
	}{
		{
			name:   "with term",
			term:   "abc",
			expect: "i18n: unable to generate a translation for term 'abc'",
		},
		{
			name:   "no term",
			term:   "",
			expect: "i18n: unable to generate a translation",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := &LocalizationError{
				Term: c.term,
			}

			actual := err.Error()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestNoTranslationGeneratedError_Is(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := new(LocalizationError)
		target := new(LocalizationError)

		is := err.Is(target)
		assert.True(t, is)
	})

	t.Run("different types", func(t *testing.T) {
		err := new(LocalizationError)
		target := io.EOF

		is := err.Is(target)
		assert.False(t, is)
	})

	t.Run("different terms", func(t *testing.T) {
		err := &LocalizationError{
			Term: "abc",
		}

		target := &LocalizationError{
			Term: "def",
		}

		is := err.Is(target)
		assert.False(t, is)
	})
}
