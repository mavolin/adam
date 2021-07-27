package arg

import (
	"time"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/duration"
)

// LocationKey is the key used to retrieve timezone information through the
// context.
// The type of the value must be *time.Location.
// If LocationKey is nil, no location is available, or the location is nil,
// DefaultLocation will be used.
// If both LocationKey and DefaultLocation are nil, UTC offsets will be
// enforced.
var LocationKey interface{}

// DefaultLocation is the time.Location used if no other timezone information
// is available.
// If this is set to nil, timezone information must be provided either
// through the UTC offset in the argument or through LocationKey's
// corresponding value.
// If both LocationKey and DefaultLocation are nil, UTC offsets will be
// enforced.
var DefaultLocation *time.Location = nil

// =============================================================================
// Duration
// =====================================================================================

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
	// If Max is 0, there won't be an upper bound.
	Max time.Duration
}

var SimpleDuration plugin.ArgType = new(Duration)

func (d Duration) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(durationName) // we have a fallback
	return name
}

func (d Duration) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(durationDescription) // we have a fallback
	return desc
}

func (d Duration) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	parsed, err := duration.Parse(ctx.Raw)

	var perr *duration.ParseError
	if errors.As(err, &perr) {
		switch perr.Code {
		case duration.ErrSize:
			return nil, newArgumentError2(durationSizeErrorArg, durationSizeErrorFlag, ctx, nil)
		case duration.ErrMissingUnit:
			return nil, newArgumentError2(durationMissingUnitErrorArg, durationMissingUnitErrorFlag, ctx, nil)
		case duration.ErrInvalidUnit:
			return nil, newArgumentError(durationInvalidUnitError, ctx, map[string]interface{}{
				"unit": perr.Val,
			})
		case duration.ErrSyntax:
			fallthrough
		default:
			return nil, newArgumentError(durationInvalidError, ctx, nil)
		}
	} else if err != nil {
		return nil, newArgumentError(durationInvalidError, ctx, nil)
	}

	if d.Min > 0 && parsed < d.Min {
		return nil, newArgumentError2(durationBelowMinErrorArg, durationBelowMinErrorFlag, ctx, map[string]interface{}{
			"min": duration.Format(d.Min),
		})
	} else if d.Max > 0 && parsed > d.Max {
		return nil, newArgumentError2(durationAboveMaxErrorArg, durationAboveMaxErrorFlag, ctx, map[string]interface{}{
			"max": duration.Format(d.Max),
		})
	}

	return parsed, nil
}

func (d Duration) GetDefault() interface{} {
	return time.Duration(0)
}

// =============================================================================
// Time
// =====================================================================================

// Time is the type used for points in time.
//
// A time can either be specified without a UTC offset following the format of
// '15:04', or with a UTC offset: '15:04 -0700'.
// In the first case, DefaultLocation will be assumed as time zone, unless
// the context has an element stored under the key LocationKey that is of type
// *time.Location.
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
var SimpleTime plugin.ArgType = new(Time)

func (t Time) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(timeName) // we have a fallback
	return name
}

func (t Time) GetDescription(l *i18n.Localizer) (desc string) {
	if LocationKey == nil && DefaultLocation == nil {
		desc, _ = l.Localize(timeDescriptionMustUTC) // we have a fallback
	} else {
		desc, _ = l.Localize(timeDescriptionOptionalUTC) // we have a fallback
	}

	return
}

var (
	timeFormat       = "15:04"
	timeFormatWithTZ = "15:04 -0700"
)

//nolint:dupl
func (t Time) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	var (
		parsed time.Time
		err    error
	)

	if len(ctx.Raw) == len(timeFormat) {
		loc := location(ctx)
		if loc == nil {
			return nil, newArgumentError2(timeRequireUTCOffsetErrorArg, timeRequireUTCOffsetErrorFlag, ctx, nil)
		}

		parsed, err = time.ParseInLocation(timeFormat, ctx.Raw, loc)
	} else if len(ctx.Raw) == len(timeFormatWithTZ) {
		parsed, err = time.Parse(timeFormatWithTZ, ctx.Raw)
	}

	if err != nil || parsed.IsZero() {
		if location(ctx) == nil { // no location available, must use utc offset
			return nil, newArgumentError2(timeInvalidErrorMustUTCArg, timeInvalidErrorMustUTCFlag, ctx, nil)
		}

		// there is location information, utc offset is optional
		return nil, newArgumentError2(timeInvalidErrorOptionalUTCArg, timeInvalidErrorOptionalUTCFlag, ctx, nil)
	}

	if !t.Min.IsZero() && parsed.Before(t.Min) {
		return nil, newArgumentError2(timeBeforeMinErrorArg, timeBeforeMinErrorFlag, ctx, map[string]interface{}{
			"min": t.Min.In(parsed.Location()).Format(timeFormat),
		})
	} else if !t.Max.IsZero() && parsed.After(t.Max) {
		return nil, newArgumentError2(timeAfterMaxErrorArg, timeAfterMaxErrorFlag, ctx, map[string]interface{}{
			"max": t.Max.In(parsed.Location()).Format(timeFormat),
		})
	}

	return parsed, nil
}

func (t Time) GetDefault() interface{} {
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
	// NoIgnoreTimeZone specifies, whether timezone information should not be
	// ignored.
	// If set to false, all time.Times generated by this Date will be in the
	// UTC time zone, regardless of whether locations through DefaultLocation
	// or LocationKey are available, or the user specified a UTC offset.
	// Furthermore, even if global settings require UTC offsets, they won't be
	// required if NoIgnoreTimeZone is set to false.
	NoIgnoreTimeZone bool
}

var (
	// SimpleDate is a Date with no bounds that doesn't require timezone
	// information.
	SimpleDate plugin.ArgType = new(Date)
	// DateWithTZ is a Date with no bounds that requires timezone information.
	DateWithTZ plugin.ArgType = &Date{NoIgnoreTimeZone: true}
)

func (d Date) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(dateName) // we have a fallback
	return name
}

func (d Date) GetDescription(l *i18n.Localizer) (desc string) {
	if LocationKey == nil && DefaultLocation == nil && d.NoIgnoreTimeZone {
		desc, _ = l.Localize(dateDescriptionMustUTC) // we have a fallback
	} else {
		desc, _ = l.Localize(dateDescriptionOptionalUTC) // we have a fallback
	}

	return
}

var (
	dateFormat       = "2006-01-02"
	dateFormatWithTZ = "2006-01-02 -0700"
)

func (d Date) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	var (
		parsed time.Time
		err    error
	)

	if len(ctx.Raw) == len(dateFormat) {
		var loc *time.Location
		if !d.NoIgnoreTimeZone {
			loc = time.UTC
		} else {
			loc = location(ctx)
			if loc == nil {
				return nil, newArgumentError2(dateRequireUTCOffsetErrorArg, dateRequireUTCOffsetErrorFlag, ctx, nil)
			}
		}

		parsed, err = time.ParseInLocation(dateFormat, ctx.Raw, loc)
	} else if len(ctx.Raw) == len(dateFormatWithTZ) {
		parsed, err = time.Parse(dateFormatWithTZ, ctx.Raw)
		if err == nil && !d.NoIgnoreTimeZone { // remove time zone, if ignored
			parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
		}
	}

	if err != nil || parsed.IsZero() {
		if !d.NoIgnoreTimeZone { // we don't need a time zone
			return nil, newArgumentError2(dateInvalidErrorNoUTCArg, dateInvalidErrorNoUTCFlag, ctx, nil)
			// no location provided but required, must use utc offset
		} else if location(ctx) == nil {
			return nil, newArgumentError2(dateInvalidErrorMustUTCArg, dateInvalidErrorMustUTCFlag, ctx, nil)
		}

		// there's location information, utc offset is optional; the location requirement is satisfied
		return nil, newArgumentError2(dateInvalidErrorOptionalUTCArg, dateInvalidErrorOptionalUTCFlag, ctx, nil)
	}

	if !d.Min.IsZero() && parsed.Before(d.Min) {
		return nil, newArgumentError2(dateBeforeMinErrorArg, dateBeforeMinErrorFlag, ctx, map[string]interface{}{
			"min": d.Min.In(parsed.Location()).Format(dateFormat),
		})
	} else if !d.Max.IsZero() && parsed.After(d.Max) {
		return nil, newArgumentError2(dateAfterMaxErrorArg, dateAfterMaxErrorFlag, ctx, map[string]interface{}{
			"max": d.Max.In(parsed.Location()).Format(dateFormat),
		})
	}

	return parsed, nil
}

func (d Date) GetDefault() interface{} {
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

// SimpleDateTime is a DateTime with no bounds
var SimpleDateTime plugin.ArgType = new(DateTime)

func (t DateTime) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(dateTimeName) // we have a fallback
	return name
}

func (t DateTime) GetDescription(l *i18n.Localizer) (desc string) {
	if LocationKey == nil && DefaultLocation == nil {
		desc, _ = l.Localize(dateTimeDescriptionMustUTC) // we have a fallback
	} else {
		desc, _ = l.Localize(dateTimeDescriptionOptionalUTC) // we have a fallback
	}

	return
}

var (
	dateTimeFormat       = "2006-01-02 15:04"
	dateTimeFormatWithTZ = "2006-01-02 15:04 -0700"
)

//nolint:dupl
func (t DateTime) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	var (
		parsed time.Time
		err    error
	)

	if len(ctx.Raw) == len(dateTimeFormat) {
		loc := location(ctx)
		if loc == nil {
			return nil, newArgumentError2(timeRequireUTCOffsetErrorArg, timeRequireUTCOffsetErrorFlag, ctx, nil)
		}

		parsed, err = time.ParseInLocation(dateTimeFormat, ctx.Raw, loc)
	} else if len(ctx.Raw) == len(dateTimeFormatWithTZ) {
		parsed, err = time.Parse(dateTimeFormatWithTZ, ctx.Raw)
	}

	if err != nil || parsed.IsZero() {
		if location(ctx) == nil { // no location provided, must use utc offset
			return nil, newArgumentError2(dateTimeInvalidErrorMustUTCArg, dateTimeInvalidErrorMustUTCFlag, ctx, nil)
		}

		// there is location information, utc offsets are optional
		return nil, newArgumentError2(dateTimeInvalidErrorOptionalUTCArg, dateTimeInvalidErrorOptionalUTCFlag, ctx, nil)
	}

	if !t.Min.IsZero() && parsed.Before(t.Min) {
		return nil, newArgumentError2(dateBeforeMinErrorArg, dateBeforeMinErrorFlag, ctx, map[string]interface{}{
			"min": t.Min.In(parsed.Location()).Format(dateTimeFormat),
		})
	} else if !t.Max.IsZero() && parsed.After(t.Max) {
		return nil, newArgumentError2(dateAfterMaxErrorArg, dateAfterMaxErrorFlag, ctx, map[string]interface{}{
			"max": t.Max.In(parsed.Location()).Format(dateTimeFormat),
		})
	}

	return parsed, nil
}

func (t DateTime) GetDefault() interface{} {
	return time.Time{}
}

// =============================================================================
// TimeZone
// =====================================================================================

// TimeZone is the Type used for time zones.
// A time zone is the name of a time zone in the IANA time zone database.
//
// You must ensure that time zone information is available on your system,
// refer to time.LoadLocation for more information.
// Alternatively, you can import time/tzdata to add timezone data to the
// executable.
//
// Go type: *time.Location
var TimeZone plugin.ArgType = new(timeZone)

type timeZone struct{}

func (z timeZone) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(timeZoneName) // we have a fallback
	return name
}

func (z timeZone) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(timeZoneDescription) // we have a fallback
	return desc
}

func (z timeZone) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	parsed, err := time.LoadLocation(ctx.Raw)
	if err != nil {
		return nil, newArgumentError(timeZoneInvalidError, ctx, nil)
	}

	return parsed, nil
}

func (z timeZone) GetDefault() interface{} {
	return (*time.Location)(nil)
}

// =============================================================================
// Helpers
// =====================================================================================

func location(ctx *plugin.ParseContext) *time.Location {
	l := DefaultLocation

	if LocationKey != nil {
		if val := ctx.Get(LocationKey); val != nil {
			if customLoc, ok := val.(*time.Location); ok && customLoc != nil {
				l = customLoc
			}
		}
	}

	return l
}
