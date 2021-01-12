package arg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestCommaConfig_Parse(t *testing.T) {
	successCases := []struct {
		name   string
		config CommaConfig

		rawArgs string

		expectArgs  plugin.Args
		expectFlags plugin.Flags
	}{
		{
			name: "flags",
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
				Required: []RequiredArg{
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
			config: CommaConfig{
				Optional: []OptionalArg{
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
			config: CommaConfig{
				Optional: []OptionalArg{
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
			config: CommaConfig{
				Optional: []OptionalArg{
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
			config: CommaConfig{
				Optional: []OptionalArg{
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
			config: CommaConfig{
				Required: []RequiredArg{
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
			config: CommaConfig{
				Required: []RequiredArg{
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
			config: CommaConfig{
				Optional: []OptionalArg{
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
			config: CommaConfig{
				Optional: []OptionalArg{
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
			config: CommaConfig{
				Optional: []OptionalArg{
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
			config: CommaConfig{
				Required: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeInt,
					},
					{
						Name: "arg2",
						Type: mockTypeString,
					},
				},
				Optional: []OptionalArg{
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
			rawArgs:    "-test2 abc, 123, def, 456, -test 789",
			expectArgs: plugin.Args{123, "def", 456, "ghi"},
			expectFlags: plugin.Flags{
				"test":  789,
				"test2": "abc",
			},
		},
		{
			name: "arg comma escape",
			config: CommaConfig{
				Required: []RequiredArg{
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
			config: CommaConfig{
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
			config: CommaConfig{
				Required: []RequiredArg{
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
			name: "no minus escape required arg",
			config: CommaConfig{
				Required: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
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
			rawArgs:    "abc, -test 123",
			expectArgs: plugin.Args{"abc", "-test 123"},
			expectFlags: plugin.Flags{
				"test": 0,
			},
		},
		{
			name: "no minus escape if no flag",
			config: CommaConfig{
				Required: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "-abc",
			expectArgs: plugin.Args{"-abc"},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				actualArgs, actualFlags, err := c.config.Parse(c.rawArgs, nil, new(plugin.Context))
				if ape, ok := err.(*plugin.ArgumentError); ok && ape != nil {
					desc, err := ape.Description(i18n.FallbackLocalizer)
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
		config CommaConfig

		rawArgs string

		expect error
	}{
		{
			name: "not enough args",
			config: CommaConfig{
				Required: []RequiredArg{
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
			config: CommaConfig{
				Required: []RequiredArg{
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
			config:  CommaConfig{},
			rawArgs: "abc",
			expect:  plugin.NewArgumentErrorl(noArgsError),
		},
		{
			name: "empty arg",
			config: CommaConfig{
				Required: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
				Optional: []OptionalArg{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
				_, _, actual := c.config.Parse(c.rawArgs, nil, new(plugin.Context))
				assert.Equal(t, c.expect, actual)
			})
		}
	})
}

func TestCommaConfig_Info(t *testing.T) {
	cfg := CommaConfig{
		Required: []RequiredArg{
			{
				Name: "r1.name",
				Type: mockType{
					name: "r1.type.name",
					desc: "r1.type.desc",
				},
				Description: "r1.desc",
			},
			{
				Name: "r2.name",
				Type: mockType{
					name: "r2.type.name",
					desc: "r2.type.desc",
				},
				Description: "r2.desc",
			},
		},
		Optional: []OptionalArg{
			{
				Name: "o1.name",
				Type: mockType{
					name: "o1.type.name",
					desc: "o1.type.desc",
				},
				Description: "o1.desc",
			},
			{
				Name: "o2.name",
				Type: mockType{
					name: "o2.type.name",
					desc: "o2.type.desc",
				},
				Description: "o2.desc",
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
				Description: "f1.desc",
			},
			{
				Name: "f2.name",
				Type: mockType{
					name: "f2.type.name",
					desc: "f2.type.desc",
				},
				Description: "f2.desc",
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

func TestLocalizedCommaConfig_Parse(t *testing.T) {
	successCases := []struct {
		name   string
		config LocalizedCommaConfig

		rawArgs string

		expectArgs  plugin.Args
		expectFlags plugin.Flags
	}{
		{
			name: "flags",
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeInt,
					},
					{
						Name: i18n.NewFallbackConfig("", "arg2"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "123, abc",
			expectArgs: plugin.Args{123, "abc"},
		},
		{
			name: "optional args",
			config: LocalizedCommaConfig{
				Optional: []LocalizedOptionalArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeInt,
					},
					{
						Name: i18n.NewFallbackConfig("", "arg2"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "123, abc",
			expectArgs: plugin.Args{123, "abc"},
		},
		{
			name: "optional args default",
			config: LocalizedCommaConfig{
				Optional: []LocalizedOptionalArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeInt,
					},
					{
						Name:    i18n.NewFallbackConfig("", "arg2"),
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
			config: LocalizedCommaConfig{
				Optional: []LocalizedOptionalArg{
					{
						Name:    i18n.NewFallbackConfig("", "arg1"),
						Type:    mockTypeInt,
						Default: 123,
					},
					{
						Name:    i18n.NewFallbackConfig("", "arg2"),
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
			config: LocalizedCommaConfig{
				Optional: []LocalizedOptionalArg{
					{
						Name: i18n.NewFallbackConfig("", "arg2"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "",
			expectArgs: plugin.Args{""},
		},
		{
			name: "single variadic required arg",
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
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
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
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
			config: LocalizedCommaConfig{
				Optional: []LocalizedOptionalArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
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
			config: LocalizedCommaConfig{
				Optional: []LocalizedOptionalArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
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
			config: LocalizedCommaConfig{
				Optional: []LocalizedOptionalArg{
					{
						Name:    i18n.NewFallbackConfig("", "arg1"),
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
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeInt,
					},
					{
						Name: i18n.NewFallbackConfig("", "arg2"),
						Type: mockTypeString,
					},
				},
				Optional: []LocalizedOptionalArg{
					{
						Name: i18n.NewFallbackConfig("", "arg3"),
						Type: mockTypeInt,
					},
					{
						Name:    i18n.NewFallbackConfig("", "arg4"),
						Type:    mockTypeString,
						Default: "ghi",
					},
				},
				Flags: []LocalizedFlag{
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
			rawArgs:    "-test2 abc, 123, def, 456, -test 789",
			expectArgs: plugin.Args{123, "def", 456, "ghi"},
			expectFlags: plugin.Flags{
				"test":  789,
				"test2": "abc",
			},
		},
		{
			name: "arg comma escape",
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "abc,, def",
			expectArgs: plugin.Args{"abc, def"},
		},
		{
			name: "flag comma escape",
			config: LocalizedCommaConfig{
				Flags: []LocalizedFlag{
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
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeString,
					},
				},
				Flags: []LocalizedFlag{
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
			name: "no minus escape required arg",
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeString,
					},
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeString,
					},
				},
				Flags: []LocalizedFlag{
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
		{
			name: "no minus escape if no flag",
			config: LocalizedCommaConfig{
				Required: []LocalizedRequiredArg{
					{
						Name: i18n.NewFallbackConfig("", "arg1"),
						Type: mockTypeString,
					},
				},
			},
			rawArgs:    "-abc",
			expectArgs: plugin.Args{"-abc"},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &plugin.Context{Localizer: i18n.FallbackLocalizer}

				actualArgs, actualFlags, err := c.config.Parse(c.rawArgs, nil, ctx)
				if ape, ok := err.(*plugin.ArgumentError); ok && ape != nil {
					desc, err := ape.Description(i18n.FallbackLocalizer)
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
		config CommaConfig

		rawArgs string

		expect error
	}{
		{
			name: "not enough args",
			config: CommaConfig{
				Required: []RequiredArg{
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
			config: CommaConfig{
				Required: []RequiredArg{
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
			config:  CommaConfig{},
			rawArgs: "abc",
			expect:  plugin.NewArgumentErrorl(noArgsError),
		},
		{
			name: "empty arg",
			config: CommaConfig{
				Required: []RequiredArg{
					{
						Name: "arg1",
						Type: mockTypeString,
					},
				},
				Optional: []OptionalArg{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
			config: CommaConfig{
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
				_, _, actual := c.config.Parse(c.rawArgs, nil, new(plugin.Context))
				assert.Equal(t, c.expect, actual)
			})
		}
	})
}

func TestLocalizedCommaConfig_Info(t *testing.T) {
	cfg := LocalizedCommaConfig{
		Required: []LocalizedRequiredArg{
			{
				Name: i18n.NewFallbackConfig("", "r1.name"),
				Type: mockType{
					name: "r1.type.name",
					desc: "r1.type.desc",
				},
				Description: i18n.NewFallbackConfig("", "r1.desc"),
			},
			{
				Name: i18n.NewFallbackConfig("", "r2.name"),
				Type: mockType{
					name: "r2.type.name",
					desc: "r2.type.desc",
				},
				Description: i18n.NewFallbackConfig("", "r2.desc"),
			},
		},
		Optional: []LocalizedOptionalArg{
			{
				Name: i18n.NewFallbackConfig("", "o1.name"),
				Type: mockType{
					name: "o1.type.name",
					desc: "o1.type.desc",
				},
				Description: i18n.NewFallbackConfig("", "o1.desc"),
			},
			{
				Name: i18n.NewFallbackConfig("", "o2.name"),
				Type: mockType{
					name: "o2.type.name",
					desc: "o2.type.desc",
				},
				Description: i18n.NewFallbackConfig("", "o2.desc"),
			},
		},
		Variadic: true,
		Flags: []LocalizedFlag{
			{
				Name:    "f1.name",
				Aliases: []string{"f1.alias.1"},
				Type: mockType{
					name: "f1.type.name",
					desc: "f1.type.desc",
				},
				Description: i18n.NewFallbackConfig("", "f1.desc"),
			},
			{
				Name: "f2.name",
				Type: mockType{
					name: "f2.type.name",
					desc: "f2.type.desc",
				},
				Description: i18n.NewFallbackConfig("", "f2.desc"),
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

	actual := cfg.Info(i18n.FallbackLocalizer)
	assert.Equal(t, expect, actual)
}
