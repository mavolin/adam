package arg

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestText_Parse(t *testing.T) {
	sucessCases := []struct {
		name string
		text Text

		raw string

		expect string
	}{
		{
			name:   "simple text",
			text:   SimpleText,
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "min length",
			text:   Text{MinLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "max length",
			text:   Text{MaxLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "regexp",
			text:   Text{Regexp: regexp.MustCompile("abc")},
			raw:    "abc",
			expect: "abc",
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range sucessCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{Raw: c.raw}

				actual, err := c.text.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name string
		text Text

		raw string

		expectArg, expectFlag *i18n.Config
	}{
		{
			name: "below min",
			text: Text{MinLength: 3},
			raw:  "ab",
			expectArg: textBelowMinLengthErrorArg.
				WithPlaceholders(map[string]interface{}{
					"min": uint(3),
				}),
			expectFlag: textBelowMinLengthErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"min": uint(3),
				}),
		},
		{
			name: "above max",
			text: Text{MaxLength: 3},
			raw:  "abcd",
			expectArg: textAboveMaxLengthErrorArg.
				WithPlaceholders(map[string]interface{}{
					"max": uint(3),
				}),
			expectFlag: textAboveMaxLengthErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"max": uint(3),
				}),
		},
		{
			name: "regexp not matching",
			text: Text{Regexp: regexp.MustCompile("abc")},
			raw:  "def",
			expectArg: regexpNotMatchingErrorArg.
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
			expectFlag: regexpNotMatchingErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
		},
		{
			name: "regexp not matching - custom error",
			text: Text{
				Regexp:          regexp.MustCompile("abc"),
				RegexpErrorArg:  i18n.NewFallbackConfig("arg", "arg"),
				RegexpErrorFlag: i18n.NewFallbackConfig("flag", "flag"),
			},
			raw: "def",
			expectArg: i18n.NewFallbackConfig("arg", "arg").
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
			expectFlag: i18n.NewFallbackConfig("flag", "flag").
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{
					Raw:  c.raw,
					Kind: KindArg,
				}

				c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, ctx)

				_, actual := c.text.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

				ctx = &Context{
					Raw:  c.raw,
					Kind: KindFlag,
				}

				c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, ctx)

				_, actual = c.text.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
			})
		}
	})
}

func TestAlphanumericID_Name(t *testing.T) {
	t.Run("default name", func(t *testing.T) {
		expect := mock.NoOpLocalizer.MustLocalize(idName)

		id := SimpleAlphanumericID

		actual := id.Name(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})

	t.Run("custom name", func(t *testing.T) {
		expect := "abc"

		id := AlphanumericID{
			CustomName: i18nutil.NewText(expect),
		}

		actual := id.Name(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})
}

func TestAlphanumericID_Description(t *testing.T) {
	t.Run("default description", func(t *testing.T) {
		expect := mock.NoOpLocalizer.MustLocalize(idDescription)

		id := SimpleAlphanumericID

		actual := id.Description(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})

	t.Run("custom description", func(t *testing.T) {
		expect := "abc"

		id := AlphanumericID{
			CustomDescription: i18nutil.NewText(expect),
		}

		actual := id.Description(mock.NoOpLocalizer)
		assert.Equal(t, expect, actual)
	})
}

func TestAlphanumericID_Parse(t *testing.T) {
	sucessCases := []struct {
		name string
		text AlphanumericID

		raw string

		expect string
	}{
		{
			name:   "simple text",
			text:   SimpleAlphanumericID,
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "min length",
			text:   AlphanumericID{MinLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "max length",
			text:   AlphanumericID{MaxLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "regexp",
			text:   AlphanumericID{Regexp: regexp.MustCompile("abc")},
			raw:    "abc",
			expect: "abc",
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range sucessCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{Raw: c.raw}

				actual, err := c.text.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name string
		text AlphanumericID

		raw string

		expectArg, expectFlag *i18n.Config
	}{
		{
			name: "below min",
			text: AlphanumericID{MinLength: 3},
			raw:  "ab",
			expectArg: idBelowMinLengthErrorArg.
				WithPlaceholders(map[string]interface{}{
					"min": uint(3),
				}),
			expectFlag: idBelowMinLengthErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"min": uint(3),
				}),
		},
		{
			name: "above max",
			text: AlphanumericID{MaxLength: 3},
			raw:  "abcd",
			expectArg: idAboveMaxLengthErrorArg.
				WithPlaceholders(map[string]interface{}{
					"max": uint(3),
				}),
			expectFlag: idAboveMaxLengthErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"max": uint(3),
				}),
		},
		{
			name: "regexp not matching",
			text: AlphanumericID{Regexp: regexp.MustCompile("abc")},
			raw:  "def",
			expectArg: regexpNotMatchingErrorArg.
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
			expectFlag: regexpNotMatchingErrorFlag.
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
		},
		{
			name: "regexp not matching - custom error",
			text: AlphanumericID{
				Regexp:          regexp.MustCompile("abc"),
				RegexpErrorArg:  i18n.NewFallbackConfig("arg", "arg"),
				RegexpErrorFlag: i18n.NewFallbackConfig("flag", "flag"),
			},
			raw: "def",
			expectArg: i18n.NewFallbackConfig("arg", "arg").
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
			expectFlag: i18n.NewFallbackConfig("flag", "flag").
				WithPlaceholders(map[string]interface{}{
					"regexp": "abc",
				}),
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{
					Raw:  c.raw,
					Kind: KindArg,
				}

				c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, ctx)

				_, actual := c.text.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

				ctx = &Context{
					Raw:  c.raw,
					Kind: KindFlag,
				}

				c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, ctx)

				_, actual = c.text.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
			})
		}
	})
}
