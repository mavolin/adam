package arg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
)

func TestInteger_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := 123

		ctx := &Context{Raw: "123"}

		actual, err := BasicInteger.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name     string
		min, max int

		raw string

		expectArg  i18n.Config
		expectFlag i18n.Config
	}{
		{
			name:       "invalid syntax",
			min:        0,
			max:        0,
			raw:        "abc",
			expectArg:  integerSyntaxError,
			expectFlag: integerSyntaxError,
		},
		{
			name:       "over bit range",
			min:        0,
			max:        0,
			raw:        "9999999999999999999999999999999999999999999999999999999",
			expectArg:  integerOverRangeErrorArg,
			expectFlag: integerOverRangeErrorFlag,
		},
		{
			name:       "under bit range",
			min:        0,
			max:        0,
			raw:        "-9999999999999999999999999999999999999999999999999999999",
			expectArg:  integerUnderRangeErrorArg,
			expectFlag: integerUnderRangeErrorFlag,
		},
		{
			name: "below min",
			min:  -3,
			max:  0,
			raw:  "-4",
			expectArg: integerBelowMinErrorArg.
				WithPlaceholders(map[string]interface{}{
					"min": -3,
				}),
			expectFlag: integerBelowMinErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"min": -3,
				}),
		},
		{
			name: "above max",
			min:  0,
			max:  5,
			raw:  "6",
			expectArg: integerAboveMaxErrorArg.
				WithPlaceholders(map[string]interface{}{
					"max": 5,
				}),
			expectFlag: integerAboveMaxErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"max": 5,
				}),
		},
	}

	for _, c := range failureCases {
		t.Run(c.name, func(t *testing.T) {
			var i Integer

			switch {
			case c.min != 0 && c.max != 0:
				i = IntegerWithBounds(c.min, c.max)
			case c.min != 0:
				i = IntegerWithMin(c.min)
			case c.max != 0:
				i = IntegerWithMax(c.max)
			}

			ctx := &Context{
				Raw:  c.raw,
				Kind: KindArgument,
			}

			c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, ctx)

			_, actual := i.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

			ctx = &Context{
				Raw:  c.raw,
				Kind: KindFlag,
			}

			c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, ctx)

			_, actual = i.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
		})
	}
}
