package arg

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestInteger_Parse(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expect := 123

		ctx := &plugin.ParseContext{Raw: "123"}

		actual, err := SimpleInteger.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name     string
		min, max int

		raw string

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
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
			expectArg:  numberOverRangeError,
			expectFlag: numberOverRangeError,
		},
		{
			name:       "under bit range",
			min:        0,
			max:        0,
			raw:        strconv.Itoa(math.MinInt64) + "9",
			expectArg:  numberBelowRangeError,
			expectFlag: numberBelowRangeError,
		},
		{
			name:         "below min",
			min:          -3,
			max:          0,
			raw:          "-4",
			expectArg:    numberBelowMinErrorArg,
			expectFlag:   numberBelowMinErrorFlag,
			placeholders: map[string]interface{}{"min": -3},
		},
		{
			name:         "above max",
			min:          0,
			max:          5,
			raw:          "6",
			expectArg:    numberAboveMaxErrorArg,
			expectFlag:   numberAboveMaxErrorFlag,
			placeholders: map[string]interface{}{"max": 5},
		},
	}

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		for _, c := range failureCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				var i plugin.ArgType

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

				ctx := &plugin.ParseContext{
					Raw:  c.raw,
					Kind: plugin.KindArg,
				}

				expect := newArgumentError(c.expectArg, ctx, c.placeholders)

				_, actual := i.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = plugin.KindFlag
				expect = newArgumentError(c.expectFlag, ctx, c.placeholders)

				_, actual = i.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func TestDecimal_Parse(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expect := 123.456

		ctx := &plugin.ParseContext{Raw: "123.456"}

		actual, err := SimpleDecimal.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name     string
		min, max float64

		raw string

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
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
			expectArg:  numberOverRangeError,
			expectFlag: numberOverRangeError,
		},
		{
			name:       "under bit range",
			min:        0,
			max:        0,
			raw:        fmt.Sprint(-1*math.MaxFloat64) + "9",
			expectArg:  numberBelowRangeError,
			expectFlag: numberBelowRangeError,
		},
		{
			name:         "below min",
			min:          -3.4,
			max:          0,
			raw:          "-3.5",
			expectArg:    numberBelowMinErrorArg,
			expectFlag:   numberBelowMinErrorFlag,
			placeholders: map[string]interface{}{"min": -3.4},
		},
		{
			name:         "above max",
			min:          0,
			max:          5.2,
			raw:          "5.3",
			expectArg:    numberAboveMaxErrorArg,
			expectFlag:   numberAboveMaxErrorFlag,
			placeholders: map[string]interface{}{"max": 5.2},
		},
	}

	for _, c := range failureCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			var d Decimal

			switch {
			case c.min != 0 && c.max != 0:
				d = DecimalWithBounds(c.min, c.max)
			case c.min != 0:
				d = DecimalWithMin(c.min)
			case c.max != 0:
				d = DecimalWithMax(c.max)
			}

			ctx := &plugin.ParseContext{
				Raw:  c.raw,
				Kind: plugin.KindArg,
			}

			expect := newArgumentError(c.expectArg, ctx, c.placeholders)

			_, actual := d.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = plugin.KindFlag
			expect = newArgumentError(c.expectFlag, ctx, c.placeholders)

			_, actual = d.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})
	}
}

func TestNumericID_Name(t *testing.T) {
	t.Parallel()

	t.Run("default name", func(t *testing.T) {
		t.Parallel()

		expect := i18n.NewFallbackLocalizer().MustLocalize(idName)

		id := SimpleNumericID

		actual := id.GetName(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("custom name", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		id := NumericID{
			CustomName: i18n.NewStaticConfig(expect),
		}

		actual := id.GetName(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})
}

func TestNumericID_Description(t *testing.T) {
	t.Parallel()

	t.Run("default description", func(t *testing.T) {
		t.Parallel()

		expect := i18n.NewFallbackLocalizer().MustLocalize(idDescription)

		id := SimpleNumericID

		actual := id.GetDescription(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("custom description", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		id := NumericID{
			CustomDescription: i18n.NewStaticConfig(expect),
		}

		actual := id.GetDescription(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})
}

func TestNumericID_Parse(t *testing.T) {
	t.Parallel()

	sucessCases := []struct {
		name string
		id   plugin.ArgType

		raw string

		expect uint64
	}{
		{
			name:   "simple id",
			id:     SimpleNumericID,
			raw:    "123",
			expect: 123,
		},
		{
			name:   "min length",
			id:     NumericID{MinLength: 3},
			raw:    "123",
			expect: 123,
		},
		{
			name:   "max length",
			id:     NumericID{MaxLength: 3},
			raw:    "123",
			expect: 123,
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range sucessCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{Raw: c.raw}

				actual, err := c.id.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name string
		id   plugin.ArgType

		raw string

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
	}{
		{
			name:         "below min",
			id:           NumericID{MinLength: 3},
			raw:          "12",
			expectArg:    idBelowMinLengthErrorArg,
			expectFlag:   idBelowMinLengthErrorFlag,
			placeholders: map[string]interface{}{"min": uint(3)},
		},
		{
			name:         "above max",
			id:           NumericID{MaxLength: 3},
			raw:          "1234",
			expectArg:    idAboveMaxLengthErrorArg,
			expectFlag:   idAboveMaxLengthErrorFlag,
			placeholders: map[string]interface{}{"max": uint(3)},
		},
		{
			name:       "not a number",
			id:         SimpleNumericID,
			raw:        "abc",
			expectArg:  idInvalidErrorArg,
			expectFlag: idInvalidErrorFlag,
		},
	}

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		for _, c := range failureCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{
					Raw:  c.raw,
					Kind: plugin.KindArg,
				}

				expect := newArgumentError(c.expectArg, ctx, c.placeholders)

				_, actual := c.id.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = plugin.KindFlag
				expect = newArgumentError(c.expectFlag, ctx, c.placeholders)

				_, actual = c.id.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}
