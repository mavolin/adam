package arg

import (
	"time"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/duration"
)

// Duration is the Type used for spans of time.
// Although time.Duration permits negative durations, Duration will return an
// error if it receives a negative duration, seeing as they are rarely desired.
//
// Go type: time.Duration
type Duration struct {
	// Min is the inclusive minimum of the duration.
	// All time.Durations below 0 will be replaced with 0.
	//
	// Defaults to 0.
	Min time.Duration
	// Max is the inclusive maximum of the duration.
	// If Max is 0, there is no maximum.
	Max time.Duration
}

var SimpleDuration Type = new(Duration)

func (d Duration) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(durationName) // we have a fallback
	return name
}

func (d Duration) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(durationDescription) // we have a fallback
	return desc
}

func (d Duration) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	parsed, err := duration.Parse(ctx.Raw)
	if perr, ok := err.(*duration.ParseError); ok {
		switch perr.Code {
		case duration.ErrSize:
			return nil, newArgParsingErr2(durationSizeErrorArg, durationSizeErrorFlag, ctx, nil)
		case duration.ErrMissingUnit:
			return nil, newArgParsingErr2(durationMissingUnitErrorArg, durationMissingUnitErrorFlag, ctx, nil)
		case duration.ErrInvalidUnit:
			return nil, newArgParsingErr(durationInvalidUnitError, ctx, map[string]interface{}{
				"unit": perr.Val,
			})
		default: // also case duration.ErrSyntax
			return nil, newArgParsingErr(durationInvalidError, ctx, nil)
		}
	} else if err != nil {
		return nil, newArgParsingErr(durationInvalidError, ctx, nil)
	}

	if d.Min > 0 && parsed < d.Min {
		return nil, newArgParsingErr2(durationBelowMinErrorArg, durationBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": duration.Format(d.Min),
		})
	} else if d.Max > 0 && parsed > d.Max {
		return nil, newArgParsingErr2(durationAboveMaxErrorArg, durationAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": duration.Format(d.Max),
		})
	}

	return parsed, nil
}

func (d Duration) Default() interface{} {
	return time.Duration(0)
}
