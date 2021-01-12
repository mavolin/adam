package arg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestOptions_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		options := Options{
			{
				Prefix: "1",
				Config: CommaConfig{
					Required: []RequiredArg{{Type: mockTypeString}},
				},
			},
		}

		expect := plugin.Args{"abc"}

		ctx := &plugin.Context{Localizer: i18n.FallbackLocalizer}

		actual, flags, err := options.Parse("1 abc", nil, ctx)
		require.NoError(t, err)
		assert.Empty(t, flags)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name    string
		args    string
		options Options

		expect error
	}{
		{
			name:   "empty args",
			args:   "",
			expect: plugin.NewArgumentErrorl(notEnoughArgsError),
		},
		{
			name: "nil config with args",
			args: "1 abc",
			options: Options{
				{
					Prefix: "1",
					Config: nil,
				},
			},
			expect: plugin.NewArgumentErrorl(tooManyArgsError),
		},
		{
			name:    "unknown prefix",
			args:    "unknown_prefix",
			options: Options{{Prefix: "1"}},
			expect: plugin.NewArgumentErrorl(unknownPrefixError.
				WithPlaceholders(unknownPrefixErrorPlaceholders{
					Name: "unknown_prefix",
				})),
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &plugin.Context{
					Localizer: i18n.FallbackLocalizer,
				}

				_, _, actual := c.options.Parse(c.args, nil, ctx)
				assert.Equal(t, c.expect, actual)
			})
		}
	})
}

func TestOptionsConfig_Info(t *testing.T) {
	options := Options{
		{
			Prefix: "1",
			Config: CommaConfig{
				Required: []RequiredArg{
					{
						Name: "1.r1.name",
						Type: mockType{
							name: "1.r1.type.name",
							desc: "1.r1.type.desc",
						},
						Description: "1.r1.desc",
					},
					{
						Name: "1.r2.name",
						Type: mockType{
							name: "1.r2.type.name",
							desc: "1.r2.type.desc",
						},
						Description: "1.r2.desc",
					},
				},
				Optional: []OptionalArg{
					{
						Name: "1.o1.name",
						Type: mockType{
							name: "1.o1.type.name",
							desc: "1.o1.type.desc",
						},
						Description: "1.o1.desc",
					},
					{
						Name: "1.o2.name",
						Type: mockType{
							name: "1.o2.type.name",
							desc: "1.o2.type.desc",
						},
						Description: "1.o2.desc",
					},
				},
				Variadic: true,
				Flags: []Flag{
					{
						Name:    "1.f1.name",
						Aliases: []string{"1.f1.alias.1"},
						Type: mockType{
							name: "1.f1.type.name",
							desc: "1.f1.type.desc",
						},
						Description: "1.f1.desc",
					},
					{
						Name: "1.f2.name",
						Type: mockType{
							name: "1.f2.type.name",
							desc: "1.f2.type.desc",
						},
						Description: "1.f2.desc",
						Multi:       true,
					},
				},
			},
		},
		{
			Prefix: "2",
			Config: CommaConfig{
				Required: []RequiredArg{
					{
						Name: "2.r1.name",
						Type: mockType{
							name: "2.r1.type.name",
							desc: "2.r1.type.desc",
						},
						Description: "2.r1.desc",
					},
					{
						Name: "2.r2.name",
						Type: mockType{
							name: "2.r2.type.name",
							desc: "2.r2.type.desc",
						},
						Description: "2.r2.desc",
					},
				},
				Optional: []OptionalArg{
					{
						Name: "2.o1.name",
						Type: mockType{
							name: "2.o1.type.name",
							desc: "2.o1.type.desc",
						},
						Description: "2.o1.desc",
					},
					{
						Name: "2.o2.name",
						Type: mockType{
							name: "2.o2.type.name",
							desc: "2.o2.type.desc",
						},
						Description: "2.o2.desc",
					},
				},
				Variadic: true,
				Flags: []Flag{
					{
						Name:    "2.f1.name",
						Aliases: []string{"2.f1.alias.1"},
						Type: mockType{
							name: "2.f1.type.name",
							desc: "2.f1.type.desc",
						},
						Description: "2.f1.desc",
					},
					{
						Name: "2.f2.name",
						Type: mockType{
							name: "2.f2.type.name",
							desc: "2.f2.type.desc",
						},
						Description: "2.f2.desc",
						Multi:       true,
					},
				},
			},
		},
	}

	expect := []plugin.ArgsInfo{
		{
			Prefix: "1",
			Required: []plugin.ArgInfo{
				{
					Name: "1.r1.name",
					Type: plugin.TypeInfo{
						Name:        "1.r1.type.name",
						Description: "1.r1.type.desc",
					},
					Description: "1.r1.desc",
				},
				{
					Name: "1.r2.name",
					Type: plugin.TypeInfo{
						Name:        "1.r2.type.name",
						Description: "1.r2.type.desc",
					},
					Description: "1.r2.desc",
				},
			},
			Optional: []plugin.ArgInfo{
				{
					Name: "1.o1.name",
					Type: plugin.TypeInfo{
						Name:        "1.o1.type.name",
						Description: "1.o1.type.desc",
					},
					Description: "1.o1.desc",
				},
				{
					Name: "1.o2.name",
					Type: plugin.TypeInfo{
						Name:        "1.o2.type.name",
						Description: "1.o2.type.desc",
					},
					Description: "1.o2.desc",
				},
			},
			Variadic: true,
			Flags: []plugin.FlagInfo{
				{
					Name:    "1.f1.name",
					Aliases: []string{"1.f1.alias.1"},
					Type: plugin.TypeInfo{
						Name:        "1.f1.type.name",
						Description: "1.f1.type.desc",
					},
					Description: "1.f1.desc",
					Multi:       false,
				},
				{
					Name:    "1.f2.name",
					Aliases: nil,
					Type: plugin.TypeInfo{
						Name:        "1.f2.type.name",
						Description: "1.f2.type.desc",
					},
					Description: "1.f2.desc",
					Multi:       true,
				},
			},
		},
		{
			Prefix: "2",
			Required: []plugin.ArgInfo{
				{
					Name: "2.r1.name",
					Type: plugin.TypeInfo{
						Name:        "2.r1.type.name",
						Description: "2.r1.type.desc",
					},
					Description: "2.r1.desc",
				},
				{
					Name: "2.r2.name",
					Type: plugin.TypeInfo{
						Name:        "2.r2.type.name",
						Description: "2.r2.type.desc",
					},
					Description: "2.r2.desc",
				},
			},
			Optional: []plugin.ArgInfo{
				{
					Name: "2.o1.name",
					Type: plugin.TypeInfo{
						Name:        "2.o1.type.name",
						Description: "2.o1.type.desc",
					},
					Description: "2.o1.desc",
				},
				{
					Name: "2.o2.name",
					Type: plugin.TypeInfo{
						Name:        "2.o2.type.name",
						Description: "2.o2.type.desc",
					},
					Description: "2.o2.desc",
				},
			},
			Variadic: true,
			Flags: []plugin.FlagInfo{
				{
					Name:    "2.f1.name",
					Aliases: []string{"2.f1.alias.1"},
					Type: plugin.TypeInfo{
						Name:        "2.f1.type.name",
						Description: "2.f1.type.desc",
					},
					Description: "2.f1.desc",
					Multi:       false,
				},
				{
					Name:    "2.f2.name",
					Aliases: nil,
					Type: plugin.TypeInfo{
						Name:        "2.f2.type.name",
						Description: "2.f2.type.desc",
					},
					Description: "2.f2.desc",
					Multi:       true,
				},
			},
		},
	}

	actual := options.Info(nil)
	assert.Equal(t, expect, actual)
}
