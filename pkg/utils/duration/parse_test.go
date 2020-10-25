package duration

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := "50min 3y 12.5 min6s3h"

		expect := 50*Minute + 3*Year + 12*Minute + 30*time.Second + 6*Second + 3*Hour

		actual, err := Parse(s)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name string

		raw string

		expect *ParseError
	}{
		{
			name:   "syntax",
			raw:    "abc",
			expect: &ParseError{Code: ErrSyntax},
		},
		{
			name:   "size",
			raw:    fmt.Sprintf("%dh", int64(math.MaxInt64)),
			expect: &ParseError{Code: ErrSize},
		},
		{
			name:   "missing unit",
			raw:    "123 456",
			expect: &ParseError{Code: ErrMissingUnit},
		},
		{
			name: "invalid unit",
			raw:  "123abc",
			expect: &ParseError{
				Code: ErrInvalidUnit,
				Val:  "abc",
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				c.expect.RawDuration = c.raw

				_, actual := Parse(c.raw)
				assert.Equal(t, c.expect, actual)
			})
		}
	})
}
