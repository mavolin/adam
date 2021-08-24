package arg

import (
	"math"
	"strconv"
	"strings"

	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
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
	SimpleInteger plugin.ArgType = new(Integer)
	// PositiveInteger is an Integer with inclusive minimum 0.
	PositiveInteger plugin.ArgType = IntegerWithMin(0)
	// NegativeInteger is an Integer with inclusive maximum -1.
	NegativeInteger plugin.ArgType = IntegerWithMax(-1)
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

func (i Integer) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(integerName) // we have a fallback
	return name
}

func (i Integer) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(integerDescription) // we have a fallback
	return desc
}

func (i Integer) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	parsed, err := strconv.Atoi(ctx.Raw)
	if err != nil {
		var nerr *strconv.NumError
		if errors.As(err, &nerr) && nerr.Err == strconv.ErrRange { //nolint:errorlint
			if strings.HasPrefix(ctx.Raw, "-") {
				return nil, newArgumentError(numberBelowRangeError, ctx, nil)
			}

			return nil, newArgumentError(numberOverRangeError, ctx, nil)
		}

		return nil, newArgumentError(integerSyntaxError, ctx, nil)
	}

	if i.Min != nil && parsed < *i.Min {
		return nil, newArgumentError2(numberBelowMinErrorArg, numberBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": *i.Min,
		})
	}

	if i.Max != nil && parsed > *i.Max {
		return nil, newArgumentError2(numberAboveMaxErrorArg, numberAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": *i.Max,
		})
	}

	return parsed, nil
}

func (i Integer) GetDefault() interface{} {
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
	// SimpleDecimal is a decimal with no bounds.
	SimpleDecimal plugin.ArgType = new(Decimal)
	// PositiveDecimal is an Decimal with inclusive minimum 0.
	PositiveDecimal plugin.ArgType = DecimalWithMin(0)
	// NegativeDecimal is an Decimal with inclusive maximum -1.
	NegativeDecimal plugin.ArgType = DecimalWithMax(-1)
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

func (i Decimal) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(decimalName) // we have a fallback
	return name
}

func (i Decimal) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(decimalDescription) // we have a fallback
	return desc
}

func (i Decimal) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	parsed, err := strconv.ParseFloat(ctx.Raw, 64)
	if err != nil || math.IsInf(parsed, 0) || math.IsNaN(parsed) {
		var nerr *strconv.NumError
		if errors.As(err, &nerr) && nerr.Err == strconv.ErrRange { //nolint:errorlint
			if strings.HasPrefix(ctx.Raw, "-") {
				return nil, newArgumentError(numberBelowRangeError, ctx, nil)
			}

			return nil, newArgumentError(numberOverRangeError, ctx, nil)
		}

		return nil, newArgumentError(decimalSyntaxError, ctx, nil)
	}

	if i.Min != nil && parsed < *i.Min {
		return nil, newArgumentError2(numberBelowMinErrorArg, numberBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": *i.Min,
		})
	}

	if i.Max != nil && parsed > *i.Max {
		return nil, newArgumentError2(numberAboveMaxErrorArg, numberAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": *i.Max,
		})
	}

	return parsed, nil
}

func (i Decimal) GetDefault() interface{} {
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
	CustomName *i18n.Config
	// CustomDescription allows you to set a custom description for the id.
	// If not set, the default description will be used.
	CustomDescription *i18n.Config

	// MinLength is the inclusive minimum length the id may have.
	MinLength uint
	// MaxLength is the inclusive maximum length the id may have.
	// If MaxLength is 0, the id won't have a maximum.
	MaxLength uint
}

var (
	// SimpleNumericID is a NumericID with no length boundaries and no custom name
	// or description.
	SimpleNumericID plugin.ArgType = new(NumericID)
	_               plugin.ArgType = NumericID{}
)

func (id NumericID) GetName(l *i18n.Localizer) string {
	if id.CustomName != nil {
		name, err := l.Localize(id.CustomName)
		if err == nil {
			return name
		}
	}

	name, _ := l.Localize(idName) // we have a fallback
	return name
}

func (id NumericID) GetDescription(l *i18n.Localizer) string {
	if id.CustomDescription != nil {
		desc, err := l.Localize(id.CustomDescription)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(idDescription) // we have a fallback
	return desc
}

func (id NumericID) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	parsed, err := strconv.ParseUint(ctx.Raw, 10, 64)
	if err != nil {
		return nil, newArgumentError2(idInvalidErrorArg, idInvalidErrorFlag, ctx, nil)
	}

	if uint(len(ctx.Raw)) < id.MinLength {
		return nil, newArgumentError2(
			idBelowMinLengthErrorArg, idBelowMinLengthErrorFlag, ctx, map[string]interface{}{
				"min": id.MinLength,
			})
	} else if id.MaxLength > 0 && uint(len(ctx.Raw)) > id.MaxLength {
		return nil, newArgumentError2(
			idAboveMaxLengthErrorArg, idAboveMaxLengthErrorFlag, ctx, map[string]interface{}{
				"max": id.MaxLength,
			})
	}

	return parsed, nil
}

func (id NumericID) GetDefault() interface{} {
	return uint64(0)
}
