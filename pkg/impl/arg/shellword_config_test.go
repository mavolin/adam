package arg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestShellwordConfig_Parse(t *testing.T) {
	successCases := []struct {
		name   string
		config ShellwordConfig

		rawArgs string

		expectArgs  plugin.Args
		expectFlags plugin.Flags
	}{
		{
			name: "flags",
			config: ShellwordConfig{
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
			config: ShellwordConfig{
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
			config: ShellwordConfig{
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
			config: ShellwordConfig{
				Flags: []Flag{
					{
						Name:    "default",
						Type:    mockTypeInt,
						Default: 123,
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
			config: ShellwordConfig{
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
			config: ShellwordConfig{
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
			config: ShellwordConfig{
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
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeInt,
					},
					{
						Name: i18nutil.NewText("arg2"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "123 'abc def'",
			expectArgs: plugin.Args{123, "abc def"},
		},
		{
			name: "optional args",
			config: ShellwordConfig{
				Optional: []OptionalArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeInt,
					},
					{
						Name: i18nutil.NewText("arg2"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "123 abc",
			expectArgs: plugin.Args{123, "abc"},
		},
		{
			name: "optional args default",
			config: ShellwordConfig{
				Optional: []OptionalArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeInt,
					},
					{
						Name:    i18nutil.NewText("arg2"),
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
			config: ShellwordConfig{
				Optional: []OptionalArg{
					{
						Name:    i18nutil.NewText("arg1"),
						Type:    mockTypeInt,
						Default: 123,
					},
					{
						Name:    i18nutil.NewText("arg2"),
						Type:    mockTypeString,
						Default: "abc",
					},
				},
				Variadic: true,
			},
			rawArgs:    "",
			expectArgs: plugin.Args{123, []string{"abc"}},
		},
		{
			name: "type default",
			config: ShellwordConfig{
				Optional: []OptionalArg{
					{
						Name: i18nutil.NewText("arg2"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "",
			expectArgs: plugin.Args{""},
		},
		{
			name: "single variadic required arg",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
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
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
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
			config: ShellwordConfig{
				Optional: []OptionalArg{
					{
						Name: i18nutil.NewText("arg1"),
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
			config: ShellwordConfig{
				Optional: []OptionalArg{
					{
						Name: i18nutil.NewText("arg1"),
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
			config: ShellwordConfig{
				Optional: []OptionalArg{
					{
						Name:    i18nutil.NewText("arg1"),
						Type:    mockTypeInt,
						Default: 123,
					},
				},
				Variadic: true,
			},
			rawArgs:    "",
			expectArgs: plugin.Args{[]int{123}},
		},
		{
			name: "flags and args",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeInt,
					},
					{
						Name: i18nutil.NewText("arg2"),
						Type: mockTypeString,
					},
				},
				Optional: []OptionalArg{
					{
						Name: i18nutil.NewText("arg3"),
						Type: mockTypeInt,
					},
					{
						Name:    i18nutil.NewText("arg4"),
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
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "'abc def'",
			expectArgs: plugin.Args{"abc def"},
		},
		{
			name: "double quotes",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    `"abc def"`,
			expectArgs: plugin.Args{"abc def"},
		},
		{
			name: "single backtick",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "`abc def`",
			expectArgs: plugin.Args{"abc def"},
		},
		{
			name: "double backtick",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "``abc ` def``",
			expectArgs: plugin.Args{"abc ` def"},
		},
		{
			name: "triple backtick",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
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
				actualArgs, actualFlags, err := c.config.Parse(c.rawArgs, nil, new(plugin.Context))
				if ape, ok := err.(*errors.ArgumentParsingError); ok && ape != nil {
					desc, err := ape.Description(mock.NoOpLocalizer)
					if err != nil {
						require.Fail(t, "Received unexpected error:\nargument parsing error")
					}

					require.Fail(t, "Received unexpected error:\n"+desc)
				}
				require.NoError(t, err)

				if len(c.expectArgs) == 0 {
					assert.Len(t, actualArgs, 0)
				} else {
					assert.Equal(t, c.expectArgs, actualArgs)
				}

				if len(c.expectFlags) == 0 {
					assert.Len(t, actualFlags, 0)
				} else {
					assert.Equal(t, c.expectFlags, actualFlags)
				}
			})
		}
	})

	failureCases := []struct {
		name   string
		config ShellwordConfig

		rawArgs string

		expect error
	}{
		{
			name: "not enough args",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "",
			expect:  errors.NewArgumentParsingErrorl(notEnoughArgsError),
		},
		{
			name: "too many args",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "abc, def",
			expect:  errors.NewArgumentParsingErrorl(tooManyArgsError),
		},
		{
			name:    "command accepts no args",
			config:  ShellwordConfig{},
			rawArgs: "abc",
			expect:  errors.NewArgumentParsingErrorl(noArgsError),
		},
		{
			name: "unknown flag",
			config: ShellwordConfig{
				Flags: []Flag{
					{
						Name: "known",
						Type: mockTypeInt,
					},
				},
			},
			rawArgs: "-known 123 -unknown flag",
			expect: errors.NewArgumentParsingErrorl(unknownFlagError.
				WithPlaceholders(unknownFlagErrorPlaceholders{
					Name: "unknown",
				})),
		},
		{
			name: "multi flag violation",
			config: ShellwordConfig{
				Flags: []Flag{
					{
						Name:  "abc",
						Type:  mockTypeInt,
						Multi: false,
					},
				},
			},
			rawArgs: "-abc 123 -abc 456",
			expect: errors.NewArgumentParsingErrorl(flagUsedMultipleTimesError.
				WithPlaceholders(flagUsedMultipleTimesErrorPlaceholders{
					Name: "abc",
				})),
		},
		{
			name: "group not closed",
			config: ShellwordConfig{
				Required: []RequiredArg{
					{
						Name: i18nutil.NewText("arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs: "'abc def",
			expect: errors.NewArgumentParsingErrorl(groupNotClosedError.
				WithPlaceholders(groupNotClosedErrorPlaceholders{
					Quote: "'",
				})),
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				_, _, actual := c.config.Parse(c.rawArgs, nil, new(plugin.Context))
				assert.Equal(t, c.expect, actual)
			})
		}
	})
}

func TestShellwordConfig_Info(t *testing.T) {
	cfg := ShellwordConfig{
		Required: []RequiredArg{
			{
				Name: i18nutil.NewText("r1.name"),
				Type: mockType{
					name: "r1.type.name",
					desc: "r1.type.desc",
				},
				Description: i18nutil.NewText("r1.desc"),
			},
			{
				Name: i18nutil.NewText("r2.name"),
				Type: mockType{
					name: "r2.type.name",
					desc: "r2.type.desc",
				},
				Description: i18nutil.NewText("r2.desc"),
			},
		},
		Optional: []OptionalArg{
			{
				Name: i18nutil.NewText("o1.name"),
				Type: mockType{
					name: "o1.type.name",
					desc: "o1.type.desc",
				},
				Description: i18nutil.NewText("o1.desc"),
			},
			{
				Name: i18nutil.NewText("o2.name"),
				Type: mockType{
					name: "o2.type.name",
					desc: "o2.type.desc",
				},
				Description: i18nutil.NewText("o2.desc"),
			},
		},
		Variadic: true,
		Flags: []Flag{
			{
				Name:    "f1.name",
				Aliases: []string{"f1.alias.1"},
				Type: mockType{
					name: "f1.type.name",
					desc: "f1.type.desc",
				},
				Description: i18nutil.NewText("f1.desc"),
			},
			{
				Name: "f2.name",
				Type: mockType{
					name: "f2.type.name",
					desc: "f2.type.desc",
				},
				Description: i18nutil.NewText("f2.desc"),
				Multi:       true,
			},
		},
	}

	expect := []plugin.ArgsInfo{
		{
			Prefix: "",
			Required: []plugin.ArgInfo{
				{
					Name: "r1.name",
					Type: plugin.TypeInfo{
						Name:        "r1.type.name",
						Description: "r1.type.desc",
					},
					Description: "r1.desc",
				},
				{
					Name: "r2.name",
					Type: plugin.TypeInfo{
						Name:        "r2.type.name",
						Description: "r2.type.desc",
					},
					Description: "r2.desc",
				},
			},
			Optional: []plugin.ArgInfo{
				{
					Name: "o1.name",
					Type: plugin.TypeInfo{
						Name:        "o1.type.name",
						Description: "o1.type.desc",
					},
					Description: "o1.desc",
				},
				{
					Name: "o2.name",
					Type: plugin.TypeInfo{
						Name:        "o2.type.name",
						Description: "o2.type.desc",
					},
					Description: "o2.desc",
				},
			},
			Variadic: true,
			Flags: []plugin.FlagInfo{
				{
					Name:    "f1.name",
					Aliases: []string{"f1.alias.1"},
					Type: plugin.TypeInfo{
						Name:        "f1.type.name",
						Description: "f1.type.desc",
					},
					Description: "f1.desc",
					Multi:       false,
				},
				{
					Name:    "f2.name",
					Aliases: nil,
					Type: plugin.TypeInfo{
						Name:        "f2.type.name",
						Description: "f2.type.desc",
					},
					Description: "f2.desc",
					Multi:       true,
				},
			},
		},
	}

	actual := cfg.Info(nil)
	assert.Equal(t, expect, actual)
}
