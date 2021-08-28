package arg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
)

func TestCode_Parse(t *testing.T) {
	t.Parallel()

	successCases := []struct {
		name string

		raw string

		expect *CodeBlock
	}{
		{
			name: "single backtick",
			raw:  "`abc`",
			expect: &CodeBlock{
				Code:         "abc",
				QtyBackticks: 1,
			},
		},
		{
			name: "double backticks",
			raw:  "``abc``",
			expect: &CodeBlock{
				Code:         "abc",
				QtyBackticks: 2,
			},
		},
		{
			name: "double backticks - allow inner backtick",
			raw:  "``abc ` def``",
			expect: &CodeBlock{
				Code:         "abc ` def",
				QtyBackticks: 2,
			},
		},
		{
			name: "triple backticks",
			raw:  "```\nabc\n```",
			expect: &CodeBlock{
				Code:         "abc",
				QtyBackticks: 3,
			},
		},
		{
			name: "triple backticks with lang",
			raw:  "```abc\ndef\n```",
			expect: &CodeBlock{
				Language:     "abc",
				Code:         "def",
				QtyBackticks: 3,
			},
		},
		{
			name: "triple backticks - allow single inner backtick",
			raw:  "```\nabc ` def\n```",
			expect: &CodeBlock{
				Code:         "abc ` def",
				QtyBackticks: 3,
			},
		},
		{
			name: "triple backticks - allow double inner backtick",
			raw:  "```\nabc `` def\n```",
			expect: &CodeBlock{
				Code:         "abc `` def",
				QtyBackticks: 3,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range successCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				ctx := &plugin.ParseContext{Raw: c.raw}

				actual, err := Code.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		ctx := &plugin.ParseContext{
			Raw:  "def not code",
			Kind: plugin.KindArg,
		}

		expect := newArgumentError(codeInvalidErrorArg, ctx, nil)

		_, actual := Code.Parse(nil, ctx)
		assert.Equal(t, expect, actual)

		ctx.Kind = plugin.KindFlag
		expect = newArgumentError(codeInvalidErrorFlag, ctx, nil)

		_, actual = Code.Parse(nil, ctx)
		assert.Equal(t, expect, actual)
	})
}
