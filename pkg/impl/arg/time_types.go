package arg

import (
	"time"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// LocationKey is the key used to retrieve timezone information through the
// context.
// The type of the value must be *time.Location.
// If no LocationKey is specified, no location is available, or the location is
// nil, DefaultLocation will be used.
var LocationKey = "location"

// DefaultLocation is the time.Location used if no other timezone information
// is available.
// If this is set to nil, timezone information must be provided either
// through the UTC offset in the argument or through LocationKey's
// corresponding value.
var DefaultLocation = time.UTC

// Time is the type used for points in time.
//
// A time can either be specified without a UTC offset following the format of
// '15:04', or with a UTC offset: '15:04 -0700'.
// In the first case, DefaultLocation will be assumed as time zone, unless
// the context has a variable called "location" that is of type *time.Location.
// If both are nil, UTC offsets will be required.
//
// Go type: time.Time
type Time struct {
	// Min is the minimum time, the used time may be.
	Min time.Time
	// Max is the maximum time, the used time may be.
	Max time.Time
}

// SimpleTime is a Time with no bounds.
var SimpleTime Type = new(Time)

func (t Time) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(timeName) // we have a fallback
	return name
}

func (t Time) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(timeDescription) // we have a fallback
	return desc
}

var (
	timeFormat       = "15:04"
	timeFormatWithTZ = "15:04 -0700"
)

func (t Time) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	var (
		parsed time.Time
		err    error
	)

	if len(ctx.Raw) == len(timeFormat) {
		loc := DefaultLocation

		if len(LocationKey) > 0 {
			if val := ctx.Get(LocationKey); val != nil {
				if customLoc, ok := val.(*time.Location); ok && customLoc != nil {
					loc = customLoc
				}
			}
		}

		if loc == nil {
			return nil, newArgParsingErr(timeRequireUTCOffsetError, ctx, nil)
		}

		parsed, err = time.ParseInLocation(timeFormat, ctx.Raw, loc)
	} else if len(ctx.Raw) == len(timeFormatWithTZ) {
		parsed, err = time.Parse(timeFormatWithTZ, ctx.Raw)
	}

	if err != nil || parsed.IsZero() {
		return nil, newArgParsingErr2(timeInvalidErrorArg, timeInvalidErrorFlag, ctx, nil)
	}

	if !t.Min.IsZero() && parsed.Before(t.Min) {
		return nil, newArgParsingErr(timeBeforeMinError, ctx, map[string]interface{}{
			"min": t.Min.In(parsed.Location()).Format(timeFormat),
		})
	} else if !t.Max.IsZero() && parsed.After(t.Max) {
		return nil, newArgParsingErr(timeAfterMaxError, ctx, map[string]interface{}{
			"max": t.Max.In(parsed.Location()).Format(timeFormat),
		})
	}

	return parsed, nil
}

func (t Time) Default() interface{} {
	return time.Time{}
}
