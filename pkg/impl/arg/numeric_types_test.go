package arg

import (
	"fmt"
	"math"
	"strconv"
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

		actual, err := SimpleInteger.Parse(nil, ctx)
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
			raw:        strconv.Itoa(math.MaxInt64) + "9",
			expectArg:  numberOverRangeErrorArg,
			expectFlag: numberOverRangeErrorFlag,
		},
		{
			name:       "under bit range",
			min:        0,
			max:        0,
			raw:        strconv.Itoa(math.MinInt64) + "9",
			expectArg:  numberUnderRangeErrorArg,
			expectFlag: numberUnderRangeErrorFlag,
		},
		{
			name: "below min",
			min:  -3,
			max:  0,
			raw:  "-4",
			expectArg: numberBelowMinErrorArg.
				WithPlaceholders(map[string]interface{}{
					"min": -3,
				}),
			expectFlag: numberBelowMinErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"min": -3,
				}),
		},
		{
			name: "above max",
			min:  0,
			max:  5,
			raw:  "6",
			expectArg: numberAboveMaxErrorArg.
				WithPlaceholders(map[string]interface{}{
					"max": 5,
				}),
			expectFlag: numberAboveMaxErrorFlag.
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

func TestDecimal_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := 123.456

		ctx := &Context{Raw: "123.456"}

		actual, err := SimpleDecimal.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name     string
		min, max float64

		raw string

		expectArg  i18n.Config
		expectFlag i18n.Config
	}{
		{
			name:       "invalid syntax",
			min:        0,
			max:        0,
			raw:        "abc",
			expectArg:  decimalSyntaxError,
			expectFlag: decimalSyntaxError,
		},
		{
			name:       "over bit range",
			min:        0,
			max:        0,
			raw:        fmt.Sprint(math.MaxFloat64) + "9",
			expectArg:  numberOverRangeErrorArg,
			expectFlag: numberOverRangeErrorFlag,
		},
		{
			name:       "under bit range",
			min:        0,
			max:        0,
			raw:        fmt.Sprint(-1*math.MaxFloat64) + "9",
			expectArg:  numberUnderRangeErrorArg,
			expectFlag: numberUnderRangeErrorFlag,
		},
		{
			name: "below min",
			min:  -3.4,
			max:  0,
			raw:  "-3.5",
			expectArg: numberBelowMinErrorArg.
				WithPlaceholders(map[string]interface{}{
					"min": -3.4,
				}),
			expectFlag: numberBelowMinErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"min": -3.4,
				}),
		},
		{
			name: "above max",
			min:  0,
			max:  5.2,
			raw:  "5.3",
			expectArg: numberAboveMaxErrorArg.
				WithPlaceholders(map[string]interface{}{
					"max": 5.2,
				}),
			expectFlag: numberAboveMaxErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"max": 5.2,
				}),
		},
	}

	for _, c := range failureCases {
		t.Run(c.name, func(t *testing.T) {
			var d Decimal

			switch {
			case c.min != 0 && c.max != 0:
				d = DecimalWithBounds(c.min, c.max)
			case c.min != 0:
				d = DecimalWithMin(c.min)
			case c.max != 0:
				d = DecimalWithMax(c.max)
			}

			ctx := &Context{
				Raw:  c.raw,
				Kind: KindArgument,
			}

			c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, ctx)

			_, actual := d.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

			ctx = &Context{
				Raw:  c.raw,
				Kind: KindFlag,
			}

			c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, ctx)

			_, actual = d.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
		})
	}
}
