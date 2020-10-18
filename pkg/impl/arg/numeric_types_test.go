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
	"github.com/mavolin/adam/pkg/utils/i18nutil"
	"github.com/mavolin/adam/pkg/utils/mock"
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

		expectArg, expectFlag *i18n.Config
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
			var i *Integer

			switch {
			case c.min != 0 && c.max != 0:
				i = IntegerWithBounds(c.min, c.max)
			case c.min != 0:
				i = IntegerWithMin(c.min)
			case c.max != 0:
				i = IntegerWithMax(c.max)
			default:
				i = SimpleInteger
			}

			ctx := &Context{
				Raw:  c.raw,
				Kind: KindArg,
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

		expectArg, expectFlag *i18n.Config
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
				Kind: KindArg,
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

func TestNumericID_Name(t *testing.T) {
	t.Run("default name", func(t *testing.T) {
		expect := mock.NoOpLocalizer.MustLocalize(idName)

		id := SimpleNumericID

		actual := id.Name(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})

	t.Run("custom name", func(t *testing.T) {
		expect := "abc"

		id := NumericID{
			CustomName: i18nutil.NewText(expect),
		}

		actual := id.Name(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})
}

func TestNumericID_Description(t *testing.T) {
	t.Run("default description", func(t *testing.T) {
		expect := mock.NoOpLocalizer.MustLocalize(idDescription)

		id := SimpleNumericID

		actual := id.Description(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})

	t.Run("custom description", func(t *testing.T) {
		expect := "abc"

		id := NumericID{
			CustomDescription: i18nutil.NewText(expect),
		}

		actual := id.Description(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})
}

func TestNumericID_Parse(t *testing.T) {
	sucessCases := []struct {
		name string
		text NumericID

		raw string

		expect uint64
	}{
		{
			name:   "simple text",
			text:   SimpleNumericID,
			raw:    "123",
			expect: 123,
		},
		{
			name:   "min length",
			text:   NumericID{MinLength: 3},
			raw:    "123",
			expect: 123,
		},
		{
			name:   "max length",
			text:   NumericID{MaxLength: 3},
			raw:    "123",
			expect: 123,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range sucessCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{Raw: c.raw}

				actual, err := c.text.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name string
		text NumericID

		raw string

		expectArg, expectFlag *i18n.Config
	}{
		{
			name: "below min",
			text: NumericID{MinLength: 3},
			raw:  "12",
			expectArg: idBelowMinLengthErrorArg.
				WithPlaceholders(map[string]interface{}{
					"min": uint(3),
				}),
			expectFlag: idBelowMinLengthErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"min": uint(3),
				}),
		},
		{
			name: "above max",
			text: NumericID{MaxLength: 3},
			raw:  "1234",
			expectArg: idAboveMaxLengthErrorArg.
				WithPlaceholders(map[string]interface{}{
					"max": uint(3),
				}),
			expectFlag: idAboveMaxLengthErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"max": uint(3),
				}),
		},
		{
			name:       "not a number",
			text:       SimpleNumericID,
			raw:        "abc",
			expectArg:  idNotANumberErrorArg,
			expectFlag: idNotANumberErrorFlag,
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{
					Raw:  c.raw,
					Kind: KindArg,
				}

				c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, ctx)

				_, actual := c.text.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

				ctx = &Context{
					Raw:  c.raw,
					Kind: KindFlag,
				}

				c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, ctx)

				_, actual = c.text.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
			})
		}
	})
}