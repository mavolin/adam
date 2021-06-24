package arg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func Test_shellwordParser_Parse(t *testing.T) {
	successCases := []struct {
		name   string
		config plugin.ArgConfig

		rawArgs string

		expectArgs  plugin.Args
		expectFlags plugin.Flags
	}{
		{
			name: "flags",
			config: &Config{
				Flags: []Flag{
					{
						Name: "test",
						Type: mockTypeInt,
					},
					{
						Name: "test2",
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "-test 123 -test2 abc",
			expectFlags: plugin.Flags{
				"test":  123,
				"test2": "abc",
			},
		},
		{
			name: "flag alias",
			config: &Config{
				Flags: []Flag{
					{
						Name:    "test",
						Aliases: []string{"t"},
						Type:    mockTypeInt,
					},
				},
			},
			rawArgs: "-t 123",
			expectFlags: plugin.Flags{
				"test": 123,
			},
		},
		{
			name: "default flag",
			config: &Config{
				Flags: []Flag{
					{
						Name:    "default",
						Type:    mockTypeInt,
						Default: 123,
					},
				},
			},
			rawArgs: "",
			expectFlags: plugin.Flags{
				"default": 123,
			},
		},
		{
			name: "default multi flag",
			config: &Config{
				Flags: []Flag{
					{
						Name:    "default",
						Type:    mockTypeInt,
						Default: []int{123},
						Multi:   true,
					},
				},
			},
			rawArgs: "",
			expectFlags: plugin.Flags{
				"default": []int{123},
			},
		},
		{
			name: "multi flag - single use",
			config: &Config{
				Flags: []Flag{
					{
						Name:  "multi",
						Type:  mockTypeInt,
						Multi: true,
					},
				},
			},
			rawArgs: "-multi 123",
			expectFlags: plugin.Flags{
				"multi": []int{123},
			},
		},
		{
			name: "multi flag - multi use",
			config: &Config{
				Flags: []Flag{
					{
						Name:  "multi",
						Type:  mockTypeInt,
						Multi: true,
					},
				},
			},
			rawArgs: "-multi 123 -multi 456",
			expectFlags: plugin.Flags{
				"multi": []int{123, 456},
			},
		},
		{
			name: "switch flag",
			config: &Config{
				Flags: []Flag{
					{
						Name: "switch",
						Type: Switch,
					},
				},
			},
			rawArgs: "-switch",
			expectFlags: plugin.Flags{
				"switch": true,
			},
		},
		{
			name: "required args",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
					{
						Name: "arg2",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "123 'abc def'",
			expectArgs: plugin.Args{123, "abc def"},
		},
		{
			name: "optional args",
			config: &Config{
				OptionalArgs: []OptionalArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
					{
						Name: "arg2",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "123 abc",
			expectArgs: plugin.Args{123, "abc"},
		},
		{
			name: "optional args default",
			config: &Config{
				OptionalArgs: []OptionalArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
					{
						Name:    "arg2",
						Type:    mockTypeString,
						Default: "abc",
					},
				},
			},
			rawArgs:    "123",
			expectArgs: plugin.Args{123, "abc"},
		},
		{
			name: "optional arg variadic default",
			config: &Config{
				OptionalArgs: []OptionalArg{
					{
						Name:    "arg1",
						Type:    mockTypeInt,
						Default: 123,
					},
					{
						Name:    "arg2",
						Type:    mockTypeString,
						Default: []string{"abc"},
					},
				},
				Variadic: true,
			},
			rawArgs:    "",
			expectArgs: plugin.Args{123, []string{"abc"}},
		},
		{
			name: "type default",
			config: &Config{
				OptionalArgs: []OptionalArg{
					{
						Name: "arg2",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "",
			expectArgs: plugin.Args{""},
		},
		{
			name: "single variadic required arg",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
				},
				Variadic: true,
			},
			rawArgs:    "123",
			expectArgs: plugin.Args{[]int{123}},
		},
		{
			name: "multiple variadic required args",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
				},
				Variadic: true,
			},
			rawArgs:    "123 456",
			expectArgs: plugin.Args{[]int{123, 456}},
		},
		{
			name: "single variadic optional arg",
			config: &Config{
				OptionalArgs: []OptionalArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
				},
				Variadic: true,
			},
			rawArgs:    "123",
			expectArgs: plugin.Args{[]int{123}},
		},
		{
			name: "multiple variadic optional args",
			config: &Config{
				OptionalArgs: []OptionalArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
				},
				Variadic: true,
			},
			rawArgs:    "123 456",
			expectArgs: plugin.Args{[]int{123, 456}},
		},
		{
			name: "variadic optional arg default",
			config: &Config{
				OptionalArgs: []OptionalArg{
					{
						Name:    "arg1",
						Type:    mockTypeInt,
						Default: []int{123},
					},
				},
				Variadic: true,
			},
			rawArgs:    "",
			expectArgs: plugin.Args{[]int{123}},
		},
		{
			name: "flags and args",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
					{
						Name: "arg2",
						Type: mockTypeString,
					},
				},
				OptionalArgs: []OptionalArg{
					{
						Name: "arg3",
						Type: mockTypeInt,
					},
					{
						Name:    "arg4",
						Type:    mockTypeString,
						Default: "ghi",
					},
				},
				Flags: []Flag{
					{
						Name: "test",
						Type: mockTypeInt,
					},
					{
						Name: "test2",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "-test2 abc 123 def 456 -test 789",
			expectArgs: plugin.Args{123, "def", 456, "ghi"},
			expectFlags: plugin.Flags{
				"test":  789,
				"test2": "abc",
			},
		},
		{
			name: "single quotes",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "'abc def'",
			expectArgs: plugin.Args{"abc def"},
		},
		{
			name: "double quotes",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    `"abc def"`,
			expectArgs: plugin.Args{"abc def"},
		},
		{
			name: "single backtick",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "`abc def`",
			expectArgs: plugin.Args{"abc def"},
		},
		{
			name: "double backtick",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "``abc ` def``",
			expectArgs: plugin.Args{"abc ` def"},
		},
		{
			name: "triple backtick",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "```abc def```",
			expectArgs: plugin.Args{"abc def"},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := new(plugin.Context)

				err := ShellwordParser.Parse(c.rawArgs, c.config, nil, ctx)
				if aerr, ok := err.(*plugin.ArgumentError); ok && aerr != nil {
					desc, err := aerr.Description(i18n.NewFallbackLocalizer())
					if err != nil {
						require.Fail(t, "Received unexpected error:\nargument parsing error")
					}

					require.Fail(t, "Received unexpected error:\n"+desc)
				}
				require.NoError(t, err)

				if len(c.expectArgs) == 0 {
					assert.Len(t, ctx.Args, 0)
				} else {
					assert.Equal(t, c.expectArgs, ctx.Args)
				}

				if len(c.expectFlags) == 0 {
					assert.Len(t, ctx.Flags, 0)
				} else {
					assert.Equal(t, c.expectFlags, ctx.Flags)
				}
			})
		}
	})

	failureCases := []struct {
		name   string
		config plugin.ArgConfig

		rawArgs string

		expect error
	}{
		{
			name: "not enough args",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "",
			expect:  plugin.NewArgumentErrorl(notEnoughArgsError),
		},
		{
			name: "too many args",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "abc, def",
			expect:  plugin.NewArgumentErrorl(tooManyArgsError),
		},
		{
			name:    "commandType accepts no args",
			config:  &Config{},
			rawArgs: "abc",
			expect:  plugin.NewArgumentErrorl(noArgsError),
		},
		{
			name: "unknown flag",
			config: &Config{
				Flags: []Flag{
					{
						Name: "known",
						Type: mockTypeInt,
					},
				},
			},
			rawArgs: "-known 123 -unknown flag",
			expect: plugin.NewArgumentErrorl(unknownFlagError.
				WithPlaceholders(unknownFlagErrorPlaceholders{
					Name: "unknown",
				})),
		},
		{
			name: "multi flag violation",
			config: &Config{
				Flags: []Flag{
					{
						Name:  "abc",
						Type:  mockTypeInt,
						Multi: false,
					},
				},
			},
			rawArgs: "-abc 123 -abc 456",
			expect: plugin.NewArgumentErrorl(flagUsedMultipleTimesError.
				WithPlaceholders(flagUsedMultipleTimesErrorPlaceholders{
					Name: "abc",
				})),
		},
		{
			name: "group not closed",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "'abc def",
			expect: plugin.NewArgumentErrorl(groupNotClosedError.
				WithPlaceholders(groupNotClosedErrorPlaceholders{
					Quote: "'",
				})),
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				actual := ShellwordParser.Parse(c.rawArgs, c.config, nil, new(plugin.Context))
				assert.Equal(t, c.expect, actual)
			})
		}
	})
}

func Test_shellwordParser_FormatArgs(t *testing.T) {
	args := []string{"Foo", "-Bar", "Foo, Bar", `Bar\ "Foo"`}
	flags := map[string]string{
		"foo":     "bar",
		"bar":     "-foo",
		"foo-bar": "bar foo",
		"bar-foo": `bar\ "Foo"`,
	}

	expectArgs := `Foo "-Bar" "Foo, Bar" "Bar\\ \"Foo\""`
	actual := ShellwordParser.FormatArgs(nil, args, flags)
	assert.True(t, strings.HasSuffix(actual, expectArgs))
	assert.Contains(t, actual[:len(actual)-len(expectArgs)+1], "-foo bar ")
	assert.Contains(t, actual[:len(actual)-len(expectArgs)+1], `-bar "-foo" `)
	assert.Contains(t, actual[:len(actual)-len(expectArgs)+1], `-foo-bar `)
	assert.Contains(t, actual[:len(actual)-len(expectArgs)+1], `-bar-foo "bar\\ \"Foo\"" `)
}

func Test_shellwordParser_FormatUsage(t *testing.T) {
	args := []string{"<Foo>", "<Bar>", "[FooBar]"}

	expect := "<Foo> <Bar> [FooBar]"
	actual := ShellwordParser.FormatUsage(nil, args)
	assert.Equal(t, expect, actual)
}

func Test_shellwordParser_FormatFlag(t *testing.T) {
	expect := "-foo"
	actual := ShellwordParser.FormatFlag("foo")
	assert.Equal(t, expect, actual)
}
