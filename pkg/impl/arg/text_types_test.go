package arg

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestText_Name(t *testing.T) {
	t.Parallel()

	t.Run("default name", func(t *testing.T) {
		t.Parallel()

		expect := i18n.NewFallbackLocalizer().MustLocalize(textName)

		txt := SimpleText

		actual := txt.GetName(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("custom name", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		txt := Text{CustomName: i18n.NewStaticConfig(expect)}

		actual := txt.GetName(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})
}

func TestText_Description(t *testing.T) {
	t.Parallel()

	t.Run("default description", func(t *testing.T) {
		t.Parallel()

		expect := i18n.NewFallbackLocalizer().MustLocalize(textDescription)

		txt := SimpleText

		actual := txt.GetDescription(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("custom description", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		txt := Text{CustomDescription: i18n.NewStaticConfig(expect)}

		actual := txt.GetDescription(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})
}

func TestText_Parse(t *testing.T) {
	t.Parallel()

	sucessCases := []struct {
		name string
		text plugin.ArgType

		raw string

		expect string
	}{
		{
			name:   "simple id",
			text:   SimpleText,
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "min length",
			text:   &Text{MinLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "max length",
			text:   &Text{MaxLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "regexp",
			text:   &Text{Regexp: regexp.MustCompile("abc")},
			raw:    "abc",
			expect: "abc",
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range sucessCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{Raw: c.raw}

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
		placeholders          map[string]interface{}
	}{
		{
			name:         "below min",
			text:         Text{MinLength: 3},
			raw:          "ab",
			expectArg:    textBelowMinLengthErrorArg,
			expectFlag:   textBelowMinLengthErrorFlag,
			placeholders: map[string]interface{}{"min": uint(3)},
		},
		{
			name:         "above max",
			text:         Text{MaxLength: 3},
			raw:          "abcd",
			expectArg:    textAboveMaxLengthErrorArg,
			expectFlag:   textAboveMaxLengthErrorFlag,
			placeholders: map[string]interface{}{"max": uint(3)},
		},
		{
			name:         "regexp not matching",
			text:         Text{Regexp: regexp.MustCompile("abc")},
			raw:          "def",
			expectArg:    regexpNotMatchingErrorArg,
			expectFlag:   regexpNotMatchingErrorFlag,
			placeholders: map[string]interface{}{"regexp": "abc"},
		},
		{
			name: "regexp not matching - custom error",
			text: Text{
				Regexp:          regexp.MustCompile("abc"),
				RegexpErrorArg:  i18n.NewFallbackConfig("arg", "arg"),
				RegexpErrorFlag: i18n.NewFallbackConfig("flag", "flag"),
			},
			raw:          "def",
			expectArg:    i18n.NewFallbackConfig("arg", "arg"),
			expectFlag:   i18n.NewFallbackConfig("flag", "flag"),
			placeholders: map[string]interface{}{"regexp": "abc"},
		},
	}

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		for _, c := range failureCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{
					Raw:  c.raw,
					Kind: plugin.KindArg,
				}

				expect := newArgumentError(c.expectArg, ctx, c.placeholders)

				_, actual := c.text.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = plugin.KindFlag
				expect = newArgumentError(c.expectFlag, ctx, c.placeholders)

				_, actual = c.text.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func TestLink_Parse(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expect := "https://github.com/mavolin/adam"

		ctx := &plugin.ParseContext{
			Kind: plugin.KindArg,
			Raw:  expect,
		}

		actual, err := SimpleLink.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name string
		link plugin.ArgType

		raw string

		expectArg, expectFlag *i18n.Config
	}{
		{
			name:       "default validator not matching",
			link:       SimpleLink,
			raw:        "ftps://abc.de",
			expectArg:  linkInvalidErrorArg,
			expectFlag: linkInvalidErrorFlag,
		},
		{
			name: "custom validator failure",
			link: &Link{
				Validator: func(u *url.URL) bool {
					return u.Host == "google"
				},
			},
			raw:        "https://bing.com",
			expectArg:  linkInvalidErrorArg,
			expectFlag: linkInvalidErrorFlag,
		},
		{
			name: "validator failure - custom error",
			link: &Link{
				ErrorArg:  i18n.NewFallbackConfig("abc", "abc"),
				ErrorFlag: i18n.NewFallbackConfig("def", "def"),
			},
			raw:        "ghi",
			expectArg:  i18n.NewFallbackConfig("abc", "abc"),
			expectFlag: i18n.NewFallbackConfig("def", "def"),
		},
	}

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		for _, c := range failureCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{
					Raw:  c.raw,
					Kind: plugin.KindArg,
				}

				expect := newArgumentError(c.expectArg, ctx, nil)

				_, actual := c.link.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = plugin.KindFlag
				expect = newArgumentError(c.expectFlag, ctx, nil)

				_, actual = c.link.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func TestAlphanumericID_Name(t *testing.T) {
	t.Parallel()

	t.Run("default name", func(t *testing.T) {
		t.Parallel()

		expect := i18n.NewFallbackLocalizer().MustLocalize(idName)

		id := SimpleAlphanumericID

		actual := id.GetName(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("custom name", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		id := AlphanumericID{CustomName: i18n.NewStaticConfig(expect)}

		actual := id.GetName(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})
}

func TestAlphanumericID_Description(t *testing.T) {
	t.Parallel()

	t.Run("default description", func(t *testing.T) {
		t.Parallel()

		expect := i18n.NewFallbackLocalizer().MustLocalize(idDescription)

		id := SimpleAlphanumericID

		actual := id.GetDescription(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("custom description", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		id := AlphanumericID{CustomDescription: i18n.NewStaticConfig(expect)}

		actual := id.GetDescription(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})
}

func TestAlphanumericID_Parse(t *testing.T) {
	t.Parallel()

	sucessCases := []struct {
		name string
		id   plugin.ArgType

		raw string

		expect string
	}{
		{
			name:   "simple id",
			id:     SimpleAlphanumericID,
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "min length",
			id:     AlphanumericID{MinLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "max length",
			id:     AlphanumericID{MaxLength: 3},
			raw:    "abc",
			expect: "abc",
		},
		{
			name:   "regexp",
			id:     AlphanumericID{Regexp: regexp.MustCompile("abc")},
			raw:    "abc",
			expect: "abc",
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range sucessCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{Raw: c.raw}

				actual, err := c.id.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name string
		id   *AlphanumericID

		raw string

		expectArg, expectFlag *i18n.Config
		placeholders          map[string]interface{}
	}{
		{
			name:         "below min",
			id:           &AlphanumericID{MinLength: 3},
			raw:          "ab",
			expectArg:    idBelowMinLengthErrorArg,
			expectFlag:   idBelowMinLengthErrorFlag,
			placeholders: map[string]interface{}{"min": uint(3)},
		},
		{
			name:         "above max",
			id:           &AlphanumericID{MaxLength: 3},
			raw:          "abcd",
			expectArg:    idAboveMaxLengthErrorArg,
			expectFlag:   idAboveMaxLengthErrorFlag,
			placeholders: map[string]interface{}{"max": uint(3)},
		},
		{
			name:         "regexp not matching",
			id:           &AlphanumericID{Regexp: regexp.MustCompile("abc")},
			raw:          "def",
			expectArg:    regexpNotMatchingErrorArg,
			expectFlag:   regexpNotMatchingErrorFlag,
			placeholders: map[string]interface{}{"regexp": "abc"},
		},
		{
			name: "regexp not matching - custom error",
			id: &AlphanumericID{
				Regexp:          regexp.MustCompile("abc"),
				RegexpErrorArg:  i18n.NewFallbackConfig("arg", "arg"),
				RegexpErrorFlag: i18n.NewFallbackConfig("flag", "flag"),
			},
			raw:          "def",
			expectArg:    i18n.NewFallbackConfig("arg", "arg"),
			expectFlag:   i18n.NewFallbackConfig("flag", "flag"),
			placeholders: map[string]interface{}{"regexp": "abc"},
		},
	}

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		for _, c := range failureCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{
					Raw:  c.raw,
					Kind: plugin.KindArg,
				}

				expect := newArgumentError(c.expectArg, ctx, c.placeholders)

				_, actual := c.id.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = plugin.KindFlag
				expect = newArgumentError(c.expectFlag, ctx, c.placeholders)

				_, actual = c.id.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}
