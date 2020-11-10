package arg

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/duration"
)

func TestDuration_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := 1*duration.Week + 3*duration.Day

		ctx := &Context{Raw: "1w 3d"}

		actual, err := SimpleDuration.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name string

		duration Type
		raw      string

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
	}{
		{
			name:       "size",
			duration:   SimpleDuration,
			raw:        fmt.Sprintf("%dh", int64(math.MaxInt64)),
			expectArg:  durationSizeErrorArg,
			expectFlag: durationSizeErrorFlag,
		},
		{
			name:       "syntax",
			duration:   SimpleDuration,
			raw:        "abc",
			expectArg:  durationInvalidError,
			expectFlag: durationInvalidError,
		},
		{
			name:       "missing unit",
			duration:   SimpleDuration,
			raw:        "123 456h",
			expectArg:  durationMissingUnitErrorArg,
			expectFlag: durationMissingUnitErrorFlag,
		},
		{
			name:         "invalid unit",
			duration:     SimpleDuration,
			raw:          "123abc",
			expectArg:    durationInvalidUnitError,
			expectFlag:   durationInvalidUnitError,
			placeholders: map[string]interface{}{"unit": "abc"},
		},
		{
			name:         "below min",
			duration:     Duration{Min: 5 * time.Second},
			raw:          "4s",
			expectArg:    durationBelowMinErrorArg,
			expectFlag:   durationBelowMinErrorFlag,
			placeholders: map[string]interface{}{"min": "5s"},
		},
		{
			name:         "above max",
			duration:     Duration{Max: 5 * time.Second},
			raw:          "6s",
			expectArg:    durationAboveMaxErrorArg,
			expectFlag:   durationAboveMaxErrorFlag,
			placeholders: map[string]interface{}{"max": "5s"},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{
					Raw:  c.raw,
					Kind: KindArg,
				}

				expect := newArgParsingErr(c.expectArg, ctx, c.placeholders)

				_, actual := c.duration.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = KindFlag

				expect = newArgParsingErr(c.expectFlag, ctx, c.placeholders)

				_, actual = c.duration.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func TestTime_Parse(t *testing.T) {
	successCases := []struct {
		name string

		raw             string
		location        *time.Location
		defaultLocation *time.Location

		expect time.Time
	}{
		{
			name:            "default timezone",
			raw:             "13:01",
			location:        nil,
			defaultLocation: time.UTC,
			expect:          time.Date(0, 1, 1, 13, 1, 0, 0, time.UTC),
		},
		{
			name:            "context timezone",
			raw:             "13:01",
			location:        time.FixedZone("CET", 200),
			defaultLocation: nil,
			expect:          time.Date(0, 1, 1, 13, 1, 0, 0, time.FixedZone("CET", 200)),
		},
		{
			name:            "utc offset",
			raw:             "13:01 +0200",
			location:        nil,
			defaultLocation: nil,
			expect:          time.Date(0, 1, 1, 13, 1, 0, 0, time.FixedZone("", 7200)),
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				DefaultLocation = c.defaultLocation

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							Base: state.NewBase(),
						},
					},
					Raw: c.raw,
				}

				ctx.Set(LocationKey, c.location)

				actual, err := SimpleTime.Parse(nil, ctx)
				require.NoError(t, err)
				require.IsType(t, time.Time{}, actual)
				if !assert.True(t, c.expect.Equal(actual.(time.Time))) { // produce a diff
					assert.Equal(t, c.expect, actual)
				}
			})
		}
	})

	failureCases := []struct {
		name string

		raw             string
		min, max        time.Time
		location        *time.Location
		defaultLocation *time.Location

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
	}{
		{
			name:            "require offset",
			raw:             "13:01",
			location:        nil,
			defaultLocation: nil,
			expectArg:       timeRequireUTCOffsetErrorArg,
			expectFlag:      timeRequireUTCOffsetErrorFlag,
		},
		{
			name:            "invalid",
			raw:             "abc",
			defaultLocation: time.UTC,
			expectArg:       timeInvalidErrorArg,
			expectFlag:      timeInvalidErrorFlag,
		},
		{
			name:       "before min",
			raw:        "13:01",
			min:        time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
			location:   time.UTC,
			expectArg:  timeBeforeMinErrorArg,
			expectFlag: timeBeforeMinErrorFlag,
			placeholders: map[string]interface{}{
				"min": time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC).Format(timeFormat),
			},
		},
		{
			name:            "after max",
			raw:             "13:01",
			min:             time.Time{},
			max:             time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			location:        time.UTC,
			defaultLocation: nil,
			expectArg:       timeAfterMaxErrorArg,
			expectFlag:      timeAfterMaxErrorFlag,
			placeholders: map[string]interface{}{
				"max": time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC).Format(timeFormat),
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				DefaultLocation = c.defaultLocation

				ti := &Time{
					Min: c.min,
					Max: c.max,
				}

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							Base: state.NewBase(),
						},
					},
					Raw:  c.raw,
					Kind: KindArg,
				}

				ctx.Set(LocationKey, c.location)

				expect := newArgParsingErr(c.expectArg, ctx, c.placeholders)

				_, actual := ti.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = KindFlag
				expect = newArgParsingErr(c.expectFlag, ctx, c.placeholders)

				_, actual = ti.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func TestDate_Parse(t *testing.T) {
	successCases := []struct {
		name string

		raw             string
		requireTimezone bool
		location        *time.Location
		defaultLocation *time.Location

		expect time.Time
	}{
		{
			name:            "default timezone",
			raw:             "2020-10-31",
			location:        nil,
			defaultLocation: time.UTC,
			requireTimezone: true,
			expect:          time.Date(2020, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "context timezone",
			raw:             "2020-10-31",
			location:        time.FixedZone("CET", 200),
			defaultLocation: nil,
			requireTimezone: true,
			expect:          time.Date(2020, 10, 31, 0, 0, 0, 0, time.FixedZone("CET", 200)),
		},
		{
			name:            "utc offset",
			raw:             "2020-10-31 +0200",
			location:        nil,
			defaultLocation: nil,
			requireTimezone: true,
			expect:          time.Date(2020, 10, 31, 0, 0, 0, 0, time.FixedZone("", 7200)),
		},
		{
			name:            "require no timezone",
			raw:             "2020-10-31",
			location:        nil,
			defaultLocation: nil,
			requireTimezone: false,
			expect:          time.Date(2020, 10, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				date := &Date{RequireTimezone: c.requireTimezone}

				DefaultLocation = c.defaultLocation

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							Base: state.NewBase(),
						},
					},
					Raw: c.raw,
				}

				ctx.Set(LocationKey, c.location)

				actual, err := date.Parse(nil, ctx)
				require.NoError(t, err)
				require.IsType(t, time.Time{}, actual)
				if !assert.True(t, c.expect.Equal(actual.(time.Time))) { // produce a diff
					assert.Equal(t, c.expect, actual)
				}
			})
		}
	})

	failureCases := []struct {
		name string

		raw             string
		requireTimezone bool
		min, max        time.Time
		location        *time.Location
		defaultLocation *time.Location

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
	}{
		{
			name:            "require offset",
			raw:             "2020-10-31",
			requireTimezone: true,
			location:        nil,
			defaultLocation: nil,
			expectArg:       dateRequireUTCOffsetErrorArg,
			expectFlag:      dateRequireUTCOffsetErrorFlag,
		},
		{
			name:            "invalid",
			raw:             "abc",
			defaultLocation: time.UTC,
			expectArg:       dateInvalidErrorArg,
			expectFlag:      dateInvalidErrorFlag,
		},
		{
			name:       "before min",
			raw:        "2020-10-31",
			min:        time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
			location:   time.UTC,
			expectArg:  dateBeforeMinErrorArg,
			expectFlag: dateBeforeMinErrorFlag,
			placeholders: map[string]interface{}{
				"min": time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Format(dateFormat),
			},
		},
		{
			name:       "after max",
			raw:        "2020-10-31",
			max:        time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC),
			location:   time.UTC,
			expectArg:  dateAfterMaxErrorArg,
			expectFlag: dateAfterMaxErrorFlag,
			placeholders: map[string]interface{}{
				"max": time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC).Format(dateFormat),
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				DefaultLocation = c.defaultLocation

				ti := &Date{
					RequireTimezone: c.requireTimezone,
					Min:             c.min,
					Max:             c.max,
				}

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							Base: state.NewBase(),
						},
					},
					Raw:  c.raw,
					Kind: KindArg,
				}

				ctx.Set(LocationKey, c.location)

				expect := newArgParsingErr(c.expectArg, ctx, c.placeholders)

				_, actual := ti.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = KindFlag
				expect = newArgParsingErr(c.expectFlag, ctx, c.placeholders)

				_, actual = ti.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func TestDateTime_Parse(t *testing.T) {
	successCases := []struct {
		name string

		raw             string
		location        *time.Location
		defaultLocation *time.Location

		expect time.Time
	}{
		{
			name:            "default timezone",
			raw:             "2020-10-31 13:01",
			location:        nil,
			defaultLocation: time.UTC,
			expect:          time.Date(2020, 10, 31, 13, 1, 0, 0, time.UTC),
		},
		{
			name:            "context timezone",
			raw:             "2020-10-31 13:01",
			location:        time.FixedZone("CET", 200),
			defaultLocation: nil,
			expect:          time.Date(2020, 10, 31, 13, 1, 0, 0, time.FixedZone("CET", 200)),
		},
		{
			name:            "utc offset",
			raw:             "2020-10-31 13:01 +0200",
			location:        nil,
			defaultLocation: nil,
			expect:          time.Date(2020, 10, 31, 13, 1, 0, 0, time.FixedZone("", 7200)),
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				DefaultLocation = c.defaultLocation

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							Base: state.NewBase(),
						},
					},
					Raw: c.raw,
				}

				ctx.Set(LocationKey, c.location)

				actual, err := SimpleDateTime.Parse(nil, ctx)
				require.NoError(t, err)
				require.IsType(t, time.Time{}, actual)
				if !assert.True(t, c.expect.Equal(actual.(time.Time))) { // produce a diff
					assert.Equal(t, c.expect, actual)
				}
			})
		}
	})

	failureCases := []struct {
		name string

		raw             string
		min, max        time.Time
		location        *time.Location
		defaultLocation *time.Location

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
	}{
		{
			name:            "require offset",
			raw:             "2020-10-31 13:01",
			location:        nil,
			defaultLocation: nil,
			expectArg:       timeRequireUTCOffsetErrorArg,
			expectFlag:      timeRequireUTCOffsetErrorFlag,
		},
		{
			name:            "invalid",
			raw:             "abc",
			defaultLocation: time.UTC,
			expectArg:       dateTimeInvalidErrorArg,
			expectFlag:      dateTimeInvalidErrorFlag,
		},
		{
			name:       "before min",
			raw:        "2020-10-31 13:01",
			min:        time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
			location:   time.UTC,
			expectArg:  dateBeforeMinErrorArg,
			expectFlag: dateBeforeMinErrorFlag,
			placeholders: map[string]interface{}{
				"min": time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Format(dateTimeFormat),
			},
		},
		{
			name:       "after max",
			raw:        "2020-10-31 13:01",
			max:        time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC),
			location:   time.UTC,
			expectArg:  dateAfterMaxErrorArg,
			expectFlag: dateAfterMaxErrorFlag,
			placeholders: map[string]interface{}{
				"max": time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC).Format(dateTimeFormat),
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				DefaultLocation = c.defaultLocation

				ti := &DateTime{
					Min: c.min,
					Max: c.max,
				}

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							Base: state.NewBase(),
						},
					},
					Raw:  c.raw,
					Kind: KindArg,
				}

				ctx.Set(LocationKey, c.location)

				expect := newArgParsingErr(c.expectArg, ctx, c.placeholders)

				_, actual := ti.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = KindFlag
				expect = newArgParsingErr(c.expectFlag, ctx, c.placeholders)

				_, actual = ti.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func TestTimeZone_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := &Context{Raw: "America/New_York"}

		expect, err := time.LoadLocation("America/New_York")
		if err != nil {
			fmt.Println("aborting TestTimeZone_Parse: no timezone data available")
			return // test os doesn't have timezone data
		}

		actual, err := TimeZone.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		ctx := &Context{Raw: "not a timezone"}

		expect := newArgParsingErr(timeZoneInvalidError, ctx, nil)

		_, actual := TimeZone.Parse(nil, ctx)
		assert.Equal(t, expect, actual)
	})
}
