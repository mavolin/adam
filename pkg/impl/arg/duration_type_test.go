package arg

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/duration"
)

func TestDuration_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := 1*duration.Week + 3*duration.Day

		ctx := &Context{Raw: "1w 3d"}

		actual, err := SimpleDuration.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name string

		duration Type
		raw      string

		expectArg, expectFlag *i18n.Config
	}{
		{
			name:       "size",
			duration:   SimpleDuration,
			raw:        fmt.Sprintf("%dh", int64(math.MaxInt64)),
			expectArg:  durationSizeErrorArg,
			expectFlag: durationSizeErrorFlag,
		},
		{
			name:       "syntax",
			duration:   SimpleDuration,
			raw:        "abc",
			expectArg:  durationInvalidError,
			expectFlag: durationInvalidError,
		},
		{
			name:       "missing unit",
			duration:   SimpleDuration,
			raw:        "123 456h",
			expectArg:  durationMissingUnitErrorArg,
			expectFlag: durationMissingUnitErrorFlag,
		},
		{
			name:     "invalid unit",
			duration: SimpleDuration,
			raw:      "123abc",
			expectArg: durationInvalidUnitError.
				WithPlaceholders(map[string]interface{}{
					"unit": "abc",
				}),
			expectFlag: durationInvalidUnitError.
				WithPlaceholders(map[string]interface{}{
					"unit": "abc",
				}),
		},
		{
			name:     "below min",
			duration: Duration{Min: 5 * time.Second},
			raw:      "4s",
			expectArg: durationBelowMinErrorArg.
				WithPlaceholders(map[string]interface{}{
					"min": "5s",
				}),
			expectFlag: durationBelowMinErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"min": "5s",
				}),
		},
		{
			name:     "above max",
			duration: Duration{Max: 5 * time.Second},
			raw:      "6s",
			expectArg: durationAboveMaxErrorArg.
				WithPlaceholders(map[string]interface{}{
					"max": "5s",
				}),
			expectFlag: durationAboveMaxErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"max": "5s",
				}),
		},
	}

	for _, c := range failureCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := &Context{
				Raw:  c.raw,
				Kind: KindArg,
			}

			expect := c.expectArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual := c.duration.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			ctx.Kind = KindFlag

			expect = c.expectFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual = c.duration.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
		})
	}
}
