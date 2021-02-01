package help

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestHelp_module(t *testing.T) {
	ctx := &plugin.Context{
		Localizer: i18n.NewFallbackLocalizer(),
		Args: plugin.Args{
			mock.GenerateRegisteredModule(plugin.BuiltInProvider, mock.Module{
				ModuleMeta: mock.ModuleMeta{
					Name:            "abc",
					LongDescription: "abc desc",
					Hidden:          true,
				},
				CommandsReturn: []plugin.Command{
					mock.Command{
						CommandMeta: mock.CommandMeta{
							Name:             "def",
							ShortDescription: "def desc",
							Hidden:           true,
						},
					},
					mock.Command{
						CommandMeta: mock.CommandMeta{
							Name:             "ghi",
							ShortDescription: "ghi desc",
						},
					},
					mock.Command{CommandMeta: mock.CommandMeta{Name: "jkl"}},
				},
				ModulesReturn: []plugin.Module{
					mock.Module{
						ModuleMeta: mock.ModuleMeta{Name: "mno"},
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
								ModuleMeta: mock.ModuleMeta{Name: "ghi"},
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
						ModuleMeta: mock.ModuleMeta{Name: "pqr"},
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
						ModuleMeta: mock.ModuleMeta{Name: "stu"},
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
			}),
		},
	}

	expect := BaseEmbed.Clone().
		WithSimpleTitlel(moduleTitle.
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
}
