package help

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestHelp_all(t *testing.T) {
	t.Run("guild", func(t *testing.T) {
		ctx := &plugin.Context{
			Message:   discord.Message{GuildID: 123},
			Localizer: i18n.NewFallbackLocalizer(),
			Prefixes:  []string{"my_cool_prefix"},
			Provider: mock.PluginProvider{
				PluginRepositoriesReturn: []plugin.Repository{
					{
						ProviderName: plugin.BuiltInProvider,
						Commands: []plugin.Command{
							mock.Command{
								CommandMeta: mock.CommandMeta{
									Name:             "abc",
									ShortDescription: "abc desc",
									Hidden:           true,
								},
							},
							mock.Command{
								CommandMeta: mock.CommandMeta{
									Name:             "def",
									ShortDescription: "def desc",
								},
							},
							mock.Command{
								CommandMeta: mock.CommandMeta{Name: "ghi"},
							},
						},
						Modules: []plugin.Module{
							mock.Module{
								ModuleMeta: mock.ModuleMeta{Name: "jkl"},
								CommandsReturn: []plugin.Command{
									mock.Command{
										CommandMeta: mock.CommandMeta{
											Name:   "abc",
											Hidden: true,
										},
									},
									mock.Command{
										CommandMeta: mock.CommandMeta{
											Name:             "def",
											ShortDescription: "def desc",
										},
									},
								},
								ModulesReturn: []plugin.Module{
									mock.Module{
										ModuleMeta: mock.ModuleMeta{
											Name: "ghi",
										},
										CommandsReturn: []plugin.Command{
											mock.Command{
												CommandMeta: mock.CommandMeta{
													Name:             "abc",
													ShortDescription: "abc desc",
												},
											},
										},
									},
								},
							},
							mock.Module{
								ModuleMeta: mock.ModuleMeta{Name: "mno"},
								CommandsReturn: []plugin.Command{
									mock.Command{
										CommandMeta: mock.CommandMeta{
											Name:             "abc",
											ShortDescription: "abc desc",
										},
									},
									mock.Command{
										CommandMeta: mock.CommandMeta{
											Name:             "def",
											ShortDescription: "def desc",
										},
									},
								},
							},
							mock.Module{
								ModuleMeta: mock.ModuleMeta{Name: "pqr"},
								CommandsReturn: []plugin.Command{
									mock.Command{
										CommandMeta: mock.CommandMeta{
											Name:   "stu,abc",
											Hidden: true,
										},
									},
								},
							},
						},
					},
				},
			},
			DiscordDataProvider: mock.DiscordDataProvider{
				SelfReturn: &discord.Member{
					User: discord.User{Username: "NotTesty"},
					Nick: "Testy",
				},
			},
		}

		expect := BaseEmbed.Clone().
			WithSimpleTitlel(allTitle).
			WithDescriptionl(allDescriptionGuild).
			WithField(ctx.MustLocalize(allPrefixesFieldName), "`@Testy`, `my_cool_prefix`").
			WithField(ctx.MustLocalize(commandsFieldName), "`def` - def desc\n`ghi`").
			WithField(ctx.MustLocalize(allModuleFieldName.
				WithPlaceholders(allModuleFieldNamePlaceholders{"jkl"})),
				"`jkl def` - def desc\n`jkl ghi abc` - abc desc").
			WithField(ctx.MustLocalize(allModuleFieldName.
				WithPlaceholders(allModuleFieldNamePlaceholders{"mno"})),
				"`mno abc` - abc desc\n`mno def` - def desc").
			MustBuild(ctx.Localizer)

		actual, err := New(Options{
			HideFuncs: []HideFunc{CheckHidden(HideList)},
		}).all(nil, ctx)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("dm", func(t *testing.T) {
		ctx := &plugin.Context{
			Localizer: i18n.NewFallbackLocalizer(),
			Provider: mock.PluginProvider{
				PluginRepositoriesReturn: []plugin.Repository{
					{
						ProviderName: plugin.BuiltInProvider,
						Commands: []plugin.Command{
							mock.Command{
								CommandMeta: mock.CommandMeta{Name: "abc"},
							},
							mock.Command{
								CommandMeta: mock.CommandMeta{Name: "def"},
							},
						},
					},
				},
			},
		}

		expect := BaseEmbed.Clone().
			WithSimpleTitlel(allTitle).
			WithDescriptionl(allDescriptionDM).
			WithField(ctx.MustLocalize(commandsFieldName), "`abc`\n`def`").
			MustBuild(ctx.Localizer)

		actual, err := New(Options{HideFuncs: []HideFunc{}}).all(nil, ctx)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})
}
