package arg

import (
	"testing"
	"time"

	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

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
	}{
		{
			name:            "require offset",
			raw:             "13:01",
			location:        nil,
			defaultLocation: nil,
			expectArg:       timeRequireUTCOffsetError,
			expectFlag:      timeRequireUTCOffsetError,
		},
		{
			name:            "invalid",
			raw:             "abc",
			defaultLocation: time.UTC,
			expectArg:       timeInvalidErrorArg,
			expectFlag:      timeInvalidErrorFlag,
		},
		{
			name:     "before min",
			raw:      "13:01",
			min:      time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
			location: time.UTC,
			expectArg: timeBeforeMinError.
				WithPlaceholders(map[string]interface{}{
					"min": time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
			expectFlag: timeBeforeMinError.
				WithPlaceholders(map[string]interface{}{
					"min": time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
		},
		{
			name:            "after max",
			raw:             "13:01",
			min:             time.Time{},
			max:             time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			location:        time.UTC,
			defaultLocation: nil,
			expectArg: timeAfterMaxError.
				WithPlaceholders(map[string]interface{}{
					"max": time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
			expectFlag: timeAfterMaxError.
				WithPlaceholders(map[string]interface{}{
					"max": time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
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

				expect := c.expectArg
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				_, actual := ti.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

				ctx.Kind = KindFlag

				expect = c.expectFlag
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				_, actual = ti.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
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
			name:            "no timezone",
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
	}{
		{
			name:            "require offset",
			raw:             "2020-10-31",
			requireTimezone: true,
			location:        nil,
			defaultLocation: nil,
			expectArg:       dateRequireUTCOffsetError,
			expectFlag:      dateRequireUTCOffsetError,
		},
		{
			name:            "invalid",
			raw:             "abc",
			defaultLocation: time.UTC,
			expectArg:       dateInvalidErrorArg,
			expectFlag:      dateInvalidErrorFlag,
		},
		{
			name:     "before min",
			raw:      "2020-10-31",
			min:      time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
			location: time.UTC,
			expectArg: dateBeforeMinError.
				WithPlaceholders(map[string]interface{}{
					"min": time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
			expectFlag: dateBeforeMinError.
				WithPlaceholders(map[string]interface{}{
					"min": time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
		},
		{
			name:            "after max",
			raw:             "2020-10-31",
			min:             time.Time{},
			max:             time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC),
			location:        time.UTC,
			defaultLocation: nil,
			expectArg: dateAfterMaxError.
				WithPlaceholders(map[string]interface{}{
					"max": time.Date(2020, 10, 31, 0, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
			expectFlag: dateAfterMaxError.
				WithPlaceholders(map[string]interface{}{
					"max": time.Date(2020, 10, 31, 0, 0, 0, 0, time.UTC).Format(timeFormat),
				}),
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

				expect := c.expectArg
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				_, actual := ti.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

				ctx.Kind = KindFlag

				expect = c.expectFlag
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				_, actual = ti.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
			})
		}
	})
}
