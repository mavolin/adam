package arg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestChoice_Parse(t *testing.T) {
	successCases := []struct {
		name   string
		choice Choice

		raw string

		expect interface{}
	}{
		{
			name:   "name",
			choice: Choice{{Name: "abc", Value: "def"}},
			raw:    "abc",
			expect: "def",
		},
		{
			name: "alias",
			choice: Choice{
				{Name: "abc", Aliases: []string{"def"}, Value: "ghi"},
			},
			raw:    "def",
			expect: "ghi",
		},
		{
			name:   "value fallback",
			choice: Choice{{Name: "abc"}},
			raw:    "abc",
			expect: "abc",
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{Raw: c.raw}

				actual, err := c.choice.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		choice := Choice{{Name: "abc"}}

		ctx := &Context{
			Raw:  "def",
			Kind: KindArg,
		}

		expect := choiceInvalidErrorArg
		expect.Placeholders = attachDefaultPlaceholders(expect, ctx)

		_, actual := choice.Parse(nil, ctx)
		assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

		ctx = &Context{
			Raw:  "def",
			Kind: KindFlag,
		}

		expect = choiceInvalidErrorFlag
		expect.Placeholders = attachDefaultPlaceholders(expect, ctx)

		_, actual = choice.Parse(nil, ctx)
		assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
	})
}

func TestChoice_Default(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		choice := Choice{{Name: "abc", Value: 123}}

		expect := 0

		actual := choice.Default()
		assert.Equal(t, expect, actual)
	})

	t.Run("string fallback", func(t *testing.T) {
		choice := Choice{{Name: "abc"}}

		expect := ""

		actual := choice.Default()
		assert.Equal(t, expect, actual)
	})

	t.Run("empty", func(t *testing.T) {
		choice := Choice{}

		var expect interface{} = nil

		actual := choice.Default()
		assert.Equal(t, expect, actual)
	})
}

func TestLocalizedChoice_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := "ghi"

		choice := LocalizedChoice{
			{
				Names: []i18n.Config{
					i18n.NewFallbackConfig("abc", "def"),
				},
				Value: expect,
			},
		}

		ctx := &Context{
			Context: &plugin.Context{
				Localizer: mock.NoOpLocalizer,
			},
			Raw: "def",
		}

		actual, err := choice.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		choice := LocalizedChoice{
			{
				Names: []i18n.Config{
					i18n.NewFallbackConfig("abc", "def"),
				},
				Value: "ghi",
			},
		}

		ctx := &Context{
			Context: &plugin.Context{
				Localizer: mock.NoOpLocalizer,
			},
			Raw:  "jkl",
			Kind: KindArg,
		}

		expect := choiceInvalidErrorArg
		expect.Placeholders = attachDefaultPlaceholders(expect, ctx)

		_, actual := choice.Parse(nil, ctx)
		assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

		ctx = &Context{
			Context: &plugin.Context{
				Localizer: mock.NoOpLocalizer,
			},
			Raw:  "jkl",
			Kind: KindFlag,
		}

		expect = choiceInvalidErrorFlag
		expect.Placeholders = attachDefaultPlaceholders(expect, ctx)

		_, actual = choice.Parse(nil, ctx)
		assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
	})
}

func TestLocalizedChoice_Default(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		choice := LocalizedChoice{
			{
				Names: []i18n.Config{
					i18n.NewFallbackConfig("abc", "def"),
				},
				Value: 123,
			},
		}

		expect := 0

		actual := choice.Default()
		assert.Equal(t, expect, actual)
	})

	t.Run("empty", func(t *testing.T) {
		choice := Choice{}

		var expect interface{} = nil

		actual := choice.Default()
		assert.Equal(t, expect, actual)
	})
}
