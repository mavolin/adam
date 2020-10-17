package arg

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// =============================================================================
// Integer
// =====================================================================================

// Integer is the type used for whole numbers.
//
// Go type: int
type Integer struct {
	// Min is the inclusive minimum of the integer.
	// If Min is nil, there is no minimum.
	Min *int
	// Max is the inclusive maximum of the integer.
	// If Max is nil, there is no maximum.
	Max *int
}

var (
	// SimpleInteger is an Integer with no bounds.
	SimpleInteger = Integer{}
	// PositiveInteger is an Integer with inclusive minimum 0.
	PositiveInteger = IntegerWithMin(0)
	// NegativeInteger is an Integer with inclusive maximum -1.
	NegativeInteger = IntegerWithMax(-1)
)

// IntegerWithMin creates a new Integer with the passed inclusive minimum.
func IntegerWithMin(min int) Integer {
	return Integer{Min: &min}
}

// IntegerWithMax creates a new Integer with the passed inclusive maximum.
func IntegerWithMax(max int) Integer {
	return Integer{Max: &max}
}

// IntegerWithBounds creates a new Integer with the passed inclusive minimum
// and maximum.
func IntegerWithBounds(min, max int) Integer {
	return Integer{
		Min: &min,
		Max: &max,
	}
}

func (i Integer) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(integerName) // we have a fallback
	return name
}

func (i Integer) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(integerDescription) // we have a fallback
	return desc
}

func (i Integer) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	parsed, err := strconv.Atoi(ctx.Raw)
	if err != nil {
		if nerr, ok := err.(*strconv.NumError); ok && nerr.Err == strconv.ErrRange {
			if strings.HasPrefix(ctx.Raw, "-") {
				return nil, newArgParsingErr(numberUnderRangeErrorArg, numberUnderRangeErrorFlag, ctx, nil)
			}

			return nil, newArgParsingErr(numberOverRangeErrorArg, numberOverRangeErrorFlag, ctx, nil)
		}

		return nil, newArgParsingErr(integerSyntaxError, integerSyntaxError, ctx, nil)
	}

	if i.Min != nil && parsed < *i.Min {
		return nil, newArgParsingErr(numberBelowMinErrorArg, numberBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": *i.Min,
		})
	}

	if i.Max != nil && parsed > *i.Max {
		return nil, newArgParsingErr(numberAboveMaxErrorArg, numberAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": *i.Max,
		})
	}

	return parsed, nil
}

func (i Integer) Default() interface{} {
	return 0
}

// =============================================================================
// Decimal
// =====================================================================================

// Decimal is the Type used for decimal numbers.
//
// Go type: float64
type Decimal struct {
	Min *float64
	Max *float64
}

var (
	// SimpleDecimal is a decimal with no bounds
	SimpleDecimal = Decimal{}
	// PositiveDecimal is an Decimal with inclusive minimum 0.
	PositiveDecimal = DecimalWithMin(0)
	// NegativeDecimal is an Decimal with inclusive maximum -1.
	NegativeDecimal = DecimalWithMax(-1)
)

// DecimalWithMin creates a new Decimal with the passed inclusive minimum.
func DecimalWithMin(min float64) Decimal {
	return Decimal{Min: &min}
}

// DecimalWithMax creates a new Decimal with the passed inclusive maximum.
func DecimalWithMax(max float64) Decimal {
	return Decimal{Max: &max}
}

// DecimalWithBounds creates a new Decimal with the passed inclusive minimum
// and maximum.
func DecimalWithBounds(min, max float64) Decimal {
	return Decimal{
		Min: &min,
		Max: &max,
	}
}

func (i Decimal) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(decimalName) // we have a fallback
	return name
}

func (i Decimal) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(decimalDescription) // we have a fallback
	return desc
}

func (i Decimal) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	parsed, err := strconv.ParseFloat(ctx.Raw, 64)
	if err != nil || math.IsInf(parsed, 0) || math.IsNaN(parsed) {
		if nerr, ok := err.(*strconv.NumError); ok && nerr.Err == strconv.ErrRange {
			if strings.HasPrefix(ctx.Raw, "-") {
				return nil, newArgParsingErr(numberUnderRangeErrorArg, numberUnderRangeErrorFlag, ctx, nil)
			}

			return nil, newArgParsingErr(numberOverRangeErrorArg, numberOverRangeErrorFlag, ctx, nil)
		}

		return nil, newArgParsingErr(decimalSyntaxError, decimalSyntaxError, ctx, nil)
	}

	fmt.Println(parsed)

	if i.Min != nil && parsed < *i.Min {
		return nil, newArgParsingErr(numberBelowMinErrorArg, numberBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": *i.Min,
		})
	}

	if i.Max != nil && parsed > *i.Max {
		return nil, newArgParsingErr(numberAboveMaxErrorArg, numberAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": *i.Max,
		})
	}

	return parsed, nil
}

func (i Decimal) Default() interface{} {
	return float64(0)
}
