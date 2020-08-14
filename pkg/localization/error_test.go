package localization

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoTranslationGeneratedError_Error(t *testing.T) {
	testCases := []struct {
		name   string
		term   string
		expect string
	}{
		{
			name:   "with term",
			term:   "abc",
			expect: "unable to generate a translation for term 'abc'",
		},
		{
			name:   "no term",
			term:   "",
			expect: "unable to generate a translation",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := &NoTranslationGeneratedError{
				Term: c.term,
			}

			actual := err.Error()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestNoTranslationGeneratedError_Is(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := new(NoTranslationGeneratedError)
		target := new(NoTranslationGeneratedError)

		is := err.Is(target)
		assert.True(t, is)
	})

	t.Run("different types", func(t *testing.T) {
		err := new(NoTranslationGeneratedError)
		target := io.EOF

		is := err.Is(target)
		assert.False(t, is)
	})

	t.Run("different terms", func(t *testing.T) {
		err := &NoTranslationGeneratedError{
			Term: "abc",
		}

		target := &NoTranslationGeneratedError{
			Term: "def",
		}

		is := err.Is(target)
		assert.False(t, is)
	})
}
