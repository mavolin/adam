package arg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestDelimiterParser_Parse(t *testing.T) {
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
			rawArgs: "-test 123, -test2 abc",
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
			rawArgs: "-multi 123, -multi 456",
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
			rawArgs:    "123, abc",
			expectArgs: plugin.Args{123, "abc"},
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
			rawArgs:    "123, abc",
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
			rawArgs:    "123, 456",
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
			rawArgs:    "123, 456",
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
			rawArgs:    "-test2 abc, -test 789, 123, def, 456",
			expectArgs: plugin.Args{123, "def", 456, "ghi"},
			expectFlags: plugin.Flags{
				"test":  789,
				"test2": "abc",
			},
		},
		{
			name: "arg comma escape",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "abc,, def",
			expectArgs: plugin.Args{"abc, def"},
		},
		{
			name: "flag comma escape",
			config: &Config{
				Flags: []Flag{
					{
						Name: "abc",
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "-abc de,,f",
			expectFlags: plugin.Flags{
				"abc": "de,f",
			},
		},
		{
			name: "minus escape",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
				Flags: []Flag{
					{
						Name: "test",
						Type: mockTypeInt,
					},
				},
			},
			rawArgs:    "--test 123",
			expectArgs: plugin.Args{"-test 123"},
			expectFlags: plugin.Flags{
				"test": 0,
			},
		},
		{
			name: "no minus escape in second arNameg",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
					{
						Name: "arg2",
						Type: mockTypeString,
					},
				},
				Flags: []Flag{
					{
						Name: "test",
						Type: mockTypeInt,
					},
				},
			},
			rawArgs:    "abc, -test 123",
			expectArgs: plugin.Args{"abc", "-test 123"},
			expectFlags: plugin.Flags{
				"test": 0,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := new(plugin.Context)

				parser := &DelimiterParser{Delimiter: ','}
				err := parser.Parse(c.rawArgs, c.config, nil, ctx)
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
			name:    "command accepts no args",
			config:  &Config{},
			rawArgs: "abc",
			expect:  plugin.NewArgumentErrorl(noArgsError),
		},
		{
			name: "empty arg",
			config: &Config{
				RequiredArgs: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
				OptionalArgs: []OptionalArg{
					{
						Name: "arg2",
						Type: mockTypeString,
					},
					{
						Name: "arg3",
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "abc, , def",
			expect: plugin.NewArgumentErrorl(emptyArgError.
				WithPlaceholders(emptyArgErrorPlaceholders{
					Position: 2,
				})),
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
			rawArgs: "-known 123, -unknown flag",
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
			rawArgs: "-abc 123, -abc 456",
			expect: plugin.NewArgumentErrorl(flagUsedMultipleTimesError.
				WithPlaceholders(flagUsedMultipleTimesErrorPlaceholders{
					Name: "abc",
				})),
		},
		{
			name: "switch with content",
			config: &Config{
				Flags: []Flag{
					{
						Name: "abc",
						Type: Switch,
					},
				},
			},
			rawArgs: "-abc 123",
			expect: plugin.NewArgumentErrorl(switchWithContentError.
				WithPlaceholders(&switchWithContentErrorPlaceholders{
					Name: "abc",
				})),
		},
		{
			name: "empty normal flag",
			config: &Config{
				Flags: []Flag{
					{
						Name: "abc",
						Type: mockTypeInt,
					},
				},
			},
			rawArgs: "-abc",
			expect: plugin.NewArgumentErrorl(emptyFlagError.
				WithPlaceholders(emptyFlagErrorPlaceholders{
					Name: "abc",
				})),
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				parser := &DelimiterParser{Delimiter: ','}

				actual := parser.Parse(c.rawArgs, c.config, nil, new(plugin.Context))
				assert.Equal(t, c.expect, actual)
			})
		}
	})
}

func TestDelimiterParser_FormatArgs(t *testing.T) {
	t.Run("no minus escape", func(t *testing.T) {
		parser := &DelimiterParser{Delimiter: ','}

		args := []string{"Foo", "Bar", "Foo, Bar"}
		flags := map[string]string{
			"foo": "bar",
			"bar": "foo",
		}

		expect := ", Foo, Bar, Foo,, Bar"
		actual := parser.FormatArgs(nil, args, flags)
		assert.True(t, strings.HasSuffix(actual, expect))
		assert.Contains(t, actual[:len(actual)-len(expect)+2], "-foo bar, ")
		assert.Contains(t, actual[:len(actual)-len(expect)+2], "-bar foo, ")
	})

	t.Run("minus escape", func(t *testing.T) {
		parser := &DelimiterParser{Delimiter: ','}

		args := []string{"-Foo", "Bar", "Foo, Bar"}

		expectArgs := "--Foo, Bar, Foo,, Bar"
		actual := parser.FormatArgs(nil, args, nil)
		assert.Equal(t, expectArgs, actual)
	})
}

func TestDelimiterParser_FormatUsage(t *testing.T) {
	parser := &DelimiterParser{Delimiter: ','}

	args := []string{"<Foo>", "<Bar>", "[FooBar]"}

	expect := "<Foo>, <Bar>, [FooBar]"
	actual := parser.FormatUsage(nil, args)
	assert.Equal(t, expect, actual)
}

func TestDelimiterParser_FormatFlag(t *testing.T) {
	expect := "-foo"
	actual := new(DelimiterParser).FormatFlag("foo")
	assert.Equal(t, expect, actual)
}
