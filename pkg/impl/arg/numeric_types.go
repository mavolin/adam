package arg

import (
	"strconv"
	"strings"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// =============================================================================
// INTEGER
// =====================================================================================

// Integer is the type used for whole numbers.
// It uses int as underlying type.
type Integer struct {
	// Min is the inclusive minimum of the integer.
	// If Min is nil, there is no minimum.
	Min *int
	// Max is the inclusive maximum of the integer.
	// If Max is nil, there is no maximum.
	Max *int
}

var (
	// BasicInteger is an Integer with no bounds.
	BasicInteger = Integer{}
	// PositiveInteger is an Integer with inclusive minimum 0.
	PositiveInteger = IntegerWithMin(0)
	// NegativeInteger is an Integer with inclusive maximum -1.
	NegativeInteger = IntegerWithMax(0)
)

// IntegerWithMin creates a new integer with the passed inclusive minimum.
func IntegerWithMin(min int) Integer {
	return Integer{Min: &min}
}

// IntegerWithMax creates a new integer with the passed inclusive maximum.
func IntegerWithMax(max int) Integer {
	return Integer{Max: &max}
}

// IntegerWithBounds creates a new integer with the passed inclusive minimum
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
				return nil, newArgParsingErr(integerUnderRangeErrorArg, integerUnderRangeErrorFlag, ctx, nil)
			}

			return nil, newArgParsingErr(integerOverRangeErrorArg, integerOverRangeErrorFlag, ctx, nil)
		}

		return nil, newArgParsingErr(integerSyntaxError, integerSyntaxError, ctx, nil)
	}

	if i.Min != nil && parsed < *i.Min {
		return nil, newArgParsingErr(integerBelowMinErrorArg, integerBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": *i.Min,
		})
	}

	if i.Max != nil && parsed > *i.Max {
		return nil, newArgParsingErr(integerAboveMaxErrorArg, integerAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": *i.Max,
		})
	}

	return parsed, nil
}

func (i Integer) Default() interface{} {
	return 0
}
