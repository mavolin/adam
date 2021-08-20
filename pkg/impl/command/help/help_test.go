package help

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/impl/module"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestHelp_Invoke(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		t.Run("guild", func(t *testing.T) {
			ctx := &plugin.Context{
				Message:   discord.Message{GuildID: 123},
				Localizer: i18n.NewFallbackLocalizer(),
				Args:      plugin.Args{nil},
				Prefixes:  []string{"my_cool_prefix"},
				Provider: mock.NewPluginProvider([]plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Commands: []plugin.Command{
							mock.Command{
								CommandMeta: command.Meta{
									Name:             "abc",
									ShortDescription: "abc desc",
									Hidden:           true,
								},
							},
							mock.Command{
								CommandMeta: command.Meta{
									Name:             "def",
									ShortDescription: "def desc",
								},
							},
							mock.Command{
								CommandMeta: command.Meta{Name: "ghi"},
							},
						},
						Modules: []plugin.Module{
							mock.Module{
								ModuleMeta: module.Meta{Name: "jkl"},
								CommandsReturn: []plugin.Command{
									mock.Command{
										CommandMeta: command.Meta{
											Name:   "abc",
											Hidden: true,
										},
									},
									mock.Command{
										CommandMeta: command.Meta{
											Name:             "def",
											ShortDescription: "def desc",
										},
									},
								},
								ModulesReturn: []plugin.Module{
									mock.Module{
										ModuleMeta: module.Meta{
											Name: "ghi",
										},
										CommandsReturn: []plugin.Command{
											mock.Command{
												CommandMeta: command.Meta{
													Name:             "abc",
													ShortDescription: "abc desc",
												},
											},
										},
									},
								},
							},
							mock.Module{
								ModuleMeta: module.Meta{Name: "mno"},
								CommandsReturn: []plugin.Command{
									mock.Command{
										CommandMeta: command.Meta{
											Name:             "abc",
											ShortDescription: "abc desc",
										},
									},
									mock.Command{
										CommandMeta: command.Meta{
											Name:             "def",
											ShortDescription: "def desc",
										},
									},
								},
							},
							mock.Module{
								ModuleMeta: module.Meta{Name: "pqr"},
								CommandsReturn: []plugin.Command{
									mock.Command{
										CommandMeta: command.Meta{
											Name:   "stu,abc",
											Hidden: true,
										},
									},
								},
							},
						},
					},
				}, nil),
				DiscordDataProvider: mock.DiscordDataProvider{
					SelfReturn: &discord.Member{
						User: discord.User{Username: "NotTesty"},
						Nick: "Testy",
					},
				},
			}

			expect := BaseEmbed.Clone().
				WithTitlel(allTitle).
				WithDescriptionl(allDescriptionGuild).
				WithField(ctx.MustLocalize(allPrefixesFieldName), "`@Testy`, `my_cool_prefix`").
				WithField(ctx.MustLocalize(commandsFieldName), "`def` - def desc\n`ghi`").
				WithField(ctx.MustLocalize(moduleTitle.
					WithPlaceholders(moduleTitlePlaceholders{"jkl"})),
					"`jkl def` - def desc\n`jkl ghi abc` - abc desc").
				WithField(ctx.MustLocalize(moduleTitle.
					WithPlaceholders(moduleTitlePlaceholders{"mno"})),
					"`mno abc` - abc desc\n`mno def` - def desc").
				MustBuild(ctx.Localizer)

			actual, err := New(Options{
				HideFuncs: []HideFunc{CheckHidden(HideList)},
			}).Invoke(nil, ctx)
			require.NoError(t, err)

			assert.Equal(t, expect, actual)
		})

		t.Run("dm", func(t *testing.T) {
			ctx := &plugin.Context{
				Localizer: i18n.NewFallbackLocalizer(),
				Args:      plugin.Args{nil},
				Provider: mock.NewPluginProvider([]plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Commands: []plugin.Command{
							mock.Command{
								CommandMeta: command.Meta{Name: "abc"},
							},
							mock.Command{
								CommandMeta: command.Meta{Name: "def"},
							},
						},
					},
				}, nil),
			}

			expect := BaseEmbed.Clone().
				WithTitlel(allTitle).
				WithDescriptionl(allDescriptionDM).
				WithField(ctx.MustLocalize(commandsFieldName), "`abc`\n`def`").
				MustBuild(ctx.Localizer)

			actual, err := New(Options{HideFuncs: []HideFunc{}}).Invoke(nil, ctx)
			require.NoError(t, err)

			assert.Equal(t, expect, actual)
		})
	})

	t.Run("module", func(t *testing.T) {
		mod := mock.Module{
			ModuleMeta: module.Meta{
				Name:            "abc",
				LongDescription: "abc desc",
			},
			CommandsReturn: []plugin.Command{
				mock.Command{
					CommandMeta: command.Meta{
						Name:             "def",
						ShortDescription: "def desc",
						Hidden:           true,
					},
				},
				mock.Command{
					CommandMeta: command.Meta{
						Name:             "ghi",
						ShortDescription: "ghi desc",
					},
				},
				mock.Command{CommandMeta: command.Meta{Name: "jkl"}},
			},
			ModulesReturn: []plugin.Module{
				mock.Module{
					ModuleMeta: module.Meta{Name: "mno"},
					CommandsReturn: []plugin.Command{
						mock.Command{
							CommandMeta: command.Meta{
								Name:   "abc",
								Hidden: true,
							},
						},
						mock.Command{
							CommandMeta: command.Meta{
								Name:             "def",
								ShortDescription: "def desc",
							},
						},
					},
					ModulesReturn: []plugin.Module{
						mock.Module{
							ModuleMeta: module.Meta{Name: "ghi"},
							CommandsReturn: []plugin.Command{
								mock.Command{
									CommandMeta: command.Meta{
										Name:             "abc",
										ShortDescription: "abc desc",
									},
								},
							},
						},
					},
				},
				mock.Module{
					ModuleMeta: module.Meta{Name: "pqr"},
					CommandsReturn: []plugin.Command{
						mock.Command{
							CommandMeta: command.Meta{
								Name:             "abc",
								ShortDescription: "abc desc",
							},
						},
						mock.Command{
							CommandMeta: command.Meta{
								Name:             "def",
								ShortDescription: "def desc",
							},
						},
					},
				},
				mock.Module{
					ModuleMeta: module.Meta{Name: "stu"},
					CommandsReturn: []plugin.Command{
						mock.Command{
							CommandMeta: command.Meta{
								Name:   "stu,abc",
								Hidden: true,
							},
						},
					},
				},
			},
		}

		ctx := &plugin.Context{
			Localizer: i18n.NewFallbackLocalizer(),
			Args:      plugin.Args{mock.ResolveModule(plugin.BuiltInSource, mod)},
		}

		expect := BaseEmbed.Clone().
			WithTitlel(moduleTitle.
				WithPlaceholders(moduleTitlePlaceholders{
					Module: "abc",
				})).
			WithDescription("abc desc").
			WithField(ctx.MustLocalize(commandsFieldName),
				"`abc ghi` - ghi desc\n`abc jkl`").
			WithField(ctx.MustLocalize(moduleTitle.
				WithPlaceholders(moduleTitlePlaceholders{"abc mno"})),
				"`abc mno def` - def desc\n`abc mno ghi abc` - abc desc").
			WithField(ctx.MustLocalize(moduleTitle.
				WithPlaceholders(moduleTitlePlaceholders{"abc pqr"})),
				"`abc pqr abc` - abc desc\n`abc pqr def` - def desc").
			MustBuild(ctx.Localizer)

		actual, err := New(Options{
			HideFuncs: []HideFunc{CheckHidden(HideList)},
		}).Invoke(nil, ctx)

		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("command", func(t *testing.T) {
		t.Run("single option", func(t *testing.T) {
			cmd := mock.Command{
				CommandMeta: command.Meta{
					Name:            "abc",
					Aliases:         []string{"cba"},
					LongDescription: "long description",
					ArgParser:       &arg.DelimiterParser{Delimiter: ','},
					Args: &arg.Config{
						RequiredArgs: []arg.RequiredArg{
							{
								Name:        "my arg",
								Type:        arg.SimpleText,
								Description: "my arg description",
							},
							{
								Name:        "date",
								Type:        arg.SimpleDate,
								Description: "date description",
							},
						},
						OptionalArgs: []arg.OptionalArg{
							{
								Name: "optional arg",
								Type: arg.SimpleInteger,
							},
							{
								Name:        "decimal",
								Type:        arg.SimpleDecimal,
								Description: "decimal description",
							},
						},
						Flags: []arg.Flag{
							{
								Name:        "flag",
								Type:        arg.User,
								Description: "flag description",
								Multi:       false,
							},
							{
								Name:    "multi",
								Aliases: []string{"m"},
								Type:    arg.SimpleAlphanumericID,
								Multi:   true,
							},
						},
						Variadic: true,
					},
					ExampleArgs: plugin.ExampleArgs{
						{Args: []string{"example one", "2021-06-24"}},
						{Args: []string{"example two", "2003-05-09"}},
					},
				},
			}

			ctx := &plugin.Context{
				Localizer: i18n.NewFallbackLocalizer(),
				Args:      plugin.Args{mock.ResolveCommand(plugin.BuiltInSource, cmd)},
			}

			expect := BaseEmbed.Clone().
				WithTitlel(commandTitle.
					WithPlaceholders(commandTitlePlaceholders{
						Command: "abc",
					})).
				WithDescription("long description").
				WithField(ctx.MustLocalize(aliasesFieldName), "`cba`").
				WithField(ctx.MustLocalize(usageFieldNameSingle), "```abc <my arg>, <date>, [optional arg], [decimal+]```").
				WithField(ctx.MustLocalize(argumentsFieldName),
					"`my arg (Text)` - my arg description\n`date` - date description\n"+
						"`decimal (Decimal+)` - decimal description").
				WithField(ctx.MustLocalize(flagsFieldName),
					"`-flag (User)` - flag description\n`-multi, -m (ID+)`").
				WithField(ctx.MustLocalize(examplesFieldName),
					"```abc example one, 2021-06-24``````abc example two, 2003-05-09```").
				MustBuild(ctx.Localizer)

			actual, err := New(Options{HideFuncs: []HideFunc{}}).Invoke(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})

		t.Run("no args", func(t *testing.T) {
			cmd := mock.Command{
				CommandMeta: command.Meta{
					Name:            "abc",
					Aliases:         []string{"cba"},
					LongDescription: "long description",
					Args:            nil,
				},
			}

			ctx := &plugin.Context{
				Localizer: i18n.NewFallbackLocalizer(),
				Args: plugin.Args{
					mock.ResolveCommand(plugin.BuiltInSource, cmd),
				},
			}

			expect := BaseEmbed.Clone().
				WithTitlel(commandTitle.
					WithPlaceholders(commandTitlePlaceholders{
						Command: "abc",
					})).
				WithDescription("long description").
				WithField(ctx.MustLocalize(aliasesFieldName), "`cba`").
				WithField(ctx.MustLocalize(usageFieldNameSingle), "```abc```").
				MustBuild(ctx.Localizer)

			actual, err := New(Options{HideFuncs: []HideFunc{}}).Invoke(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})
}
