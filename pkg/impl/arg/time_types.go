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

// =============================================================================
// Time
// =====================================================================================

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
	// Min is the inclusive minimum time.
	Min time.Time
	// Max is the inclusive maximum time.
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
			return nil, newArgParsingErr2(timeRequireUTCOffsetErrorArg, timeRequireUTCOffsetErrorFlag, ctx, nil)
		}

		parsed, err = time.ParseInLocation(timeFormat, ctx.Raw, loc)
	} else if len(ctx.Raw) == len(timeFormatWithTZ) {
		parsed, err = time.Parse(timeFormatWithTZ, ctx.Raw)
	}

	if err != nil || parsed.IsZero() {
		return nil, newArgParsingErr2(timeInvalidErrorArg, timeInvalidErrorFlag, ctx, nil)
	}

	if !t.Min.IsZero() && parsed.Before(t.Min) {
		return nil, newArgParsingErr2(timeBeforeMinErrorArg, timeBeforeMinErrorFlag, ctx, map[string]interface{}{
			"min": t.Min.In(parsed.Location()).Format(timeFormat),
		})
	} else if !t.Max.IsZero() && parsed.After(t.Max) {
		return nil, newArgParsingErr2(timeAfterMaxErrorArg, timeAfterMaxErrorFlag, ctx, map[string]interface{}{
			"max": t.Max.In(parsed.Location()).Format(timeFormat),
		})
	}

	return parsed, nil
}

func (t Time) Default() interface{} {
	return time.Time{}
}

// =============================================================================
// Date
// =====================================================================================

// Date is the type used for dates.
//
// A Date can either be specified without a UTC offset following the format of
// '2006-01-02', or with a UTC offset: '2006-01-02 -0700'.
// However, timezones can be disabled, in which case UTC offsets will be
// ignored.
//
// If the first format is used, DefaultLocation will be assumed as time zone,
// unless the context has a variable called "location" that is of type
// *time.Location.
// If both are nil, UTC offsets will be required.
//
// Go type: time.Time
type Date struct {
	// Min is the inclusive minimum date.
	Min time.Time
	// Max is the inclusive maximum time.
	Max time.Time
	// RequireTimezone specifies, whether timezone information is required.
	// If the Date contains a UTC offset, it will be ignored.
	RequireTimezone bool
}

var (
	// SimpleDate is a Date with no bounds that doesn't require timezone
	// information.
	SimpleDate Type = new(Date)
	// DateWithTZ is a Date with no bounds that requires timezone information.
	DateWithTZ Type = &Date{RequireTimezone: true}
)

func (t Date) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(dateName) // we have a fallback
	return name
}

func (t Date) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(dateDescription) // we have a fallback
	return desc
}

var (
	dateFormat       = "2006-01-02"
	dateFormatWithTZ = "2006-01-02 -0700"
)

func (t Date) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	var (
		parsed time.Time
		err    error
	)

	if len(ctx.Raw) == len(dateFormat) {
		loc := DefaultLocation

		if len(LocationKey) > 0 {
			if val := ctx.Get(LocationKey); val != nil {
				if customLoc, ok := val.(*time.Location); ok && customLoc != nil {
					loc = customLoc
				}
			}
		}

		if loc == nil {
			if t.RequireTimezone {
				return nil, newArgParsingErr2(dateRequireUTCOffsetErrorArg, dateRequireUTCOffsetErrorFlag, ctx, nil)
			}

			loc = time.UTC
		}

		parsed, err = time.ParseInLocation(dateFormat, ctx.Raw, loc)
	} else if len(ctx.Raw) == len(dateFormatWithTZ) {
		parsed, err = time.Parse(dateFormatWithTZ, ctx.Raw)
	}

	if err != nil || parsed.IsZero() {
		return nil, newArgParsingErr2(dateInvalidErrorArg, dateInvalidErrorFlag, ctx, nil)
	}

	if !t.Min.IsZero() && parsed.Before(t.Min) {
		return nil, newArgParsingErr2(dateBeforeMinErrorArg, dateBeforeMinErrorFlag, ctx, map[string]interface{}{
			"min": t.Min.In(parsed.Location()).Format(dateFormat),
		})
	} else if !t.Max.IsZero() && parsed.After(t.Max) {
		return nil, newArgParsingErr2(dateAfterMaxErrorArg, dateAfterMaxErrorFlag, ctx, map[string]interface{}{
			"max": t.Max.In(parsed.Location()).Format(dateFormat),
		})
	}

	return parsed, nil
}

func (t Date) Default() interface{} {
	return time.Time{}
}

// =============================================================================
// DateTime
// =====================================================================================

// DateTime is the type used for combinations of a date and a time.
//
// A DateTime can either be specified without a UTC offset following the format
// of '2006-01-02 15:04', or with a UTC offset: '2006-01-02 15:04 -0700'.
// In the first case, DefaultLocation will be assumed as time zone, unless
// the context has a variable called "location" that is of type *time.Location.
// If both are nil, UTC offsets will be required.
//
// Go type: time.Time
type DateTime struct {
	// Min is the inclusive minimum date.
	Min time.Time
	// Max is the inclusive maximum time.
	Max time.Time
}

var (
	// SimpleDateTime is a DateTime with no bounds
	SimpleDateTime Type = new(DateTime)
)

func (t DateTime) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(dateTimeName) // we have a fallback
	return name
}

func (t DateTime) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(dateTimeDescription) // we have a fallback
	return desc
}

var (
	dateTimeFormat       = "2006-01-02 15:04"
	dateTimeFormatWithTZ = "2006-01-02 15:04 -0700"
)

func (t DateTime) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	var (
		parsed time.Time
		err    error
	)

	if len(ctx.Raw) == len(dateTimeFormat) {
		loc := DefaultLocation

		if len(LocationKey) > 0 {
			if val := ctx.Get(LocationKey); val != nil {
				if customLoc, ok := val.(*time.Location); ok && customLoc != nil {
					loc = customLoc
				}
			}
		}

		if loc == nil {
			return nil, newArgParsingErr2(timeRequireUTCOffsetErrorArg, timeRequireUTCOffsetErrorFlag, ctx, nil)
		}

		parsed, err = time.ParseInLocation(dateTimeFormat, ctx.Raw, loc)
	} else if len(ctx.Raw) == len(dateTimeFormatWithTZ) {
		parsed, err = time.Parse(dateTimeFormatWithTZ, ctx.Raw)
	}

	if err != nil || parsed.IsZero() {
		return nil, newArgParsingErr2(dateTimeInvalidErrorArg, dateTimeInvalidErrorFlag, ctx, nil)
	}

	if !t.Min.IsZero() && parsed.Before(t.Min) {
		return nil, newArgParsingErr2(dateBeforeMinErrorArg, dateBeforeMinErrorFlag, ctx, map[string]interface{}{
			"min": t.Min.In(parsed.Location()).Format(dateTimeFormat),
		})
	} else if !t.Max.IsZero() && parsed.After(t.Max) {
		return nil, newArgParsingErr2(dateAfterMaxErrorArg, dateAfterMaxErrorFlag, ctx, map[string]interface{}{
			"max": t.Max.In(parsed.Location()).Format(dateTimeFormat),
		})
	}

	return parsed, nil
}

func (t DateTime) Default() interface{} {
	return time.Time{}
}
