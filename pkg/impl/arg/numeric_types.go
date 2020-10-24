package arg

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
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
	SimpleInteger Type = new(Integer)
	// PositiveInteger is an Integer with inclusive minimum 0.
	PositiveInteger Type = IntegerWithMin(0)
	// NegativeInteger is an Integer with inclusive maximum -1.
	NegativeInteger Type = IntegerWithMax(-1)
)

// IntegerWithMin creates a new Integer with the passed inclusive minimum.
func IntegerWithMin(min int) *Integer {
	return &Integer{Min: &min}
}

// IntegerWithMax creates a new Integer with the passed inclusive maximum.
func IntegerWithMax(max int) *Integer {
	return &Integer{Max: &max}
}

// IntegerWithBounds creates a new Integer with the passed inclusive minimum
// and maximum.
func IntegerWithBounds(min, max int) *Integer {
	return &Integer{
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
				return nil, newArgParsingErr(numberBelowRangeError, ctx, nil)
			}

			return nil, newArgParsingErr(numberOverRangeError, ctx, nil)
		}

		return nil, newArgParsingErr(integerSyntaxError, ctx, nil)
	}

	if i.Min != nil && parsed < *i.Min {
		return nil, newArgParsingErr2(numberBelowMinErrorArg, numberBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": *i.Min,
		})
	}

	if i.Max != nil && parsed > *i.Max {
		return nil, newArgParsingErr2(numberAboveMaxErrorArg, numberAboveMaxErrorFlag, ctx, map[string]interface{}{
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
	SimpleDecimal Type = new(Decimal)
	// PositiveDecimal is an Decimal with inclusive minimum 0.
	PositiveDecimal Type = DecimalWithMin(0)
	// NegativeDecimal is an Decimal with inclusive maximum -1.
	NegativeDecimal Type = DecimalWithMax(-1)
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
				return nil, newArgParsingErr(numberBelowRangeError, ctx, nil)
			}

			return nil, newArgParsingErr(numberOverRangeError, ctx, nil)
		}

		return nil, newArgParsingErr(decimalSyntaxError, ctx, nil)
	}

	fmt.Println(parsed)

	if i.Min != nil && parsed < *i.Min {
		return nil, newArgParsingErr2(numberBelowMinErrorArg, numberBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": *i.Min,
		})
	}

	if i.Max != nil && parsed > *i.Max {
		return nil, newArgParsingErr2(numberAboveMaxErrorArg, numberAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": *i.Max,
		})
	}

	return parsed, nil
}

func (i Decimal) Default() interface{} {
	return float64(0)
}

// =============================================================================
// NumericID
// =====================================================================================

// NumericID is the Type used for ids consisting only of numbers.
// Additionally, ids must be positive.
// By default, NumericIDs share the same name and description as
// AlphanumericIDs, simply their definition differs.
//
// In contrast to an AlphaNumericID, a NumericID returns uint64s.
// This also means, it is capped at 64 bit positive integers.
// If your IDs exceed that limit, use a AlphanumericID with a regular
// expression.
//
// Go type: uint64
type NumericID struct {
	// CustomName allows you to set a custom name for the id.
	// If not set, the default name will be used.
	CustomName *i18nutil.Text
	// CustomDescription allows you to set a custom description for the id.
	// If not set, the default description will be used.
	CustomDescription *i18nutil.Text

	// MinLength is the inclusive minimum length the id may have.
	MinLength uint
	// MaxLength is the inclusive maximum length the text may have.
	// If MaxLength is 0, the text won't have a maximum.
	MaxLength uint
}

// SimpleNumericID is a NumericID with no length boundaries and no custom name
// or description.
var SimpleNumericID Type = new(NumericID)
var _ Type = NumericID{}

func (id NumericID) Name(l *i18n.Localizer) string {
	if id.CustomName != nil {
		name, err := id.CustomName.Get(l)
		if err == nil {
			return name
		}
	}

	name, _ := l.Localize(idName) // we have id fallback
	return name
}

func (id NumericID) Description(l *i18n.Localizer) string {
	if id.CustomDescription != nil {
		desc, err := id.CustomDescription.Get(l)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(idDescription) // we have id fallback
	return desc
}

func (id NumericID) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	parsed, err := strconv.ParseUint(ctx.Raw, 10, 64)
	if err != nil {
		return nil, newArgParsingErr2(idNotANumberErrorArg, idNotANumberErrorFlag, ctx, nil)
	}

	if uint(len(ctx.Raw)) < id.MinLength {
		return nil, newArgParsingErr2(
			idBelowMinLengthErrorArg, idBelowMinLengthErrorFlag, ctx, map[string]interface{}{
				"min": id.MinLength,
			})
	} else if id.MaxLength > 0 && uint(len(ctx.Raw)) > id.MaxLength {
		return nil, newArgParsingErr2(
			idAboveMaxLengthErrorArg, idAboveMaxLengthErrorFlag, ctx, map[string]interface{}{
				"max": id.MaxLength,
			})
	}

	return parsed, nil
}

func (id NumericID) Default() interface{} {
	return uint64(0)
}
