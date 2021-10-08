package help

import (
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/module"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestCheckHidden(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		hiddenLvl HiddenLevel
		cmd       plugin.ResolvedCommand

		expect HiddenLevel
	}{
		{
			name:      "hidden",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Hidden: true,
			}),
			expect: HideList,
		},
		{
			name:      "not hidden",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Hidden: false,
			}),
			expect: Show,
		},
		{
			name:      "level",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Hidden: true,
			}),
			expect: HideList,
		},
		{
			name:      "level",
			hiddenLvl: Hide,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Hidden: true,
			}),
			expect: Hide,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := CheckHidden(c.hiddenLvl)(c.cmd, nil, nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestCheckChannelTypes(t *testing.T) {
	t.Parallel()

	successCases := []struct {
		name      string
		hiddenLvl HiddenLevel
		cmd       plugin.ResolvedCommand
		ctx       *plugin.Context

		expect HiddenLevel
	}{
		{
			name:      "matching",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				ChannelTypes: plugin.DirectMessages,
			}),
			ctx:    &plugin.Context{Message: discord.Message{GuildID: 0}},
			expect: Show,
		},
		{
			name:      "not matching",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				ChannelTypes: plugin.DirectMessages,
			}),
			ctx:    &plugin.Context{Message: discord.Message{GuildID: 123}},
			expect: HideList,
		},
		{
			name:      "level",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				ChannelTypes: plugin.DirectMessages,
			}),
			ctx:    &plugin.Context{Message: discord.Message{GuildID: 123}},
			expect: HideList,
		},
		{
			name:      "level",
			hiddenLvl: Hide,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				ChannelTypes: plugin.DirectMessages,
			}),
			ctx:    &plugin.Context{Message: discord.Message{GuildID: 123}},
			expect: Hide,
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range successCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				actual := CheckChannelTypes(c.hiddenLvl)(c.cmd, nil, c.ctx)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		expect := Show

		cmd := mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
			ChannelTypes: plugin.GuildTextChannels,
		})

		channelError := errors.New("abc")

		errHandler := mock.NewErrorHandler(t).
			ExpectSilentError(channelError)

		ctx := &plugin.Context{
			Message: discord.Message{GuildID: 123},
			DiscordDataProvider: mock.DiscordDataProvider{
				ChannelError: channelError,
			},
			ErrorHandler: errHandler,
		}

		actual := CheckChannelTypes(expect)(cmd, nil, ctx)
		assert.Equal(t, expect, actual)
	})
}

func TestCheckRestrictions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		hiddenLvl HiddenLevel
		cmd       plugin.ResolvedCommand

		expect HiddenLevel
	}{
		{
			name:      "restricted",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Restrictions: mock.RestrictionFunc(plugin.DefaultFatalRestrictionError),
			}),
			expect: HideList,
		},
		{
			name:      "not restricted",
			hiddenLvl: HideList,
			cmd:       mock.ResolveCommand(plugin.BuiltInSource, mock.Command{}),
			expect:    Show,
		},
		{
			name:      "not fatal",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Restrictions: mock.RestrictionFunc(plugin.DefaultRestrictionError),
			}),
			expect: Show,
		},
		{
			name:      "level",
			hiddenLvl: HideList,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Restrictions: mock.RestrictionFunc(plugin.DefaultFatalRestrictionError),
			}),
			expect: HideList,
		},
		{
			name:      "level",
			hiddenLvl: Hide,
			cmd: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
				Restrictions: mock.RestrictionFunc(plugin.DefaultFatalRestrictionError),
			}),
			expect: Hide,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := CheckRestrictions(c.hiddenLvl)(c.cmd, nil, nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

// =============================================================================
// Utilities
// =====================================================================================

func TestHelp_calcCommandHiddenLevel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		funcs []HideFunc

		expect HiddenLevel
	}{
		{
			name:   "success",
			funcs:  []HideFunc{mockHideFunc(Show), mockHideFunc(HideList), mockHideFunc(Show)},
			expect: HideList,
		},
		{
			name:   "hide is max",
			funcs:  []HideFunc{mockHideFunc(Show), mockHideFunc(HideList), mockHideFunc(Hide + 1)},
			expect: Hide,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			h := &Help{Options: Options{HideFuncs: c.funcs}}

			actual := h.calcCommandHiddenLevel(nil, nil, nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestHelp_calcModuleHiddenLevel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		lvls []HiddenLevel

		expect HiddenLevel
	}{
		{
			name:   "success",
			lvls:   []HiddenLevel{Hide, HideList, Hide},
			expect: HideList,
		},
		{
			name:   "hide is max",
			lvls:   []HiddenLevel{Hide + 1},
			expect: Hide,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			h := &Help{Options: Options{HideFuncs: make([]HideFunc, len(c.lvls))}}

			mod := module.New(module.Meta{})

			for i, lvl := range c.lvls {
				lvl := lvl
				istr := strconv.Itoa(i)

				h.HideFuncs[i] = func(cmd plugin.ResolvedCommand, _ *state.State, _ *plugin.Context) HiddenLevel {
					if cmd.Name() == istr {
						return lvl
					}

					return Show
				}

				mod.AddCommand(mock.Command{Name: istr})
			}

			actual := h.calcModuleHiddenLevel(nil, nil, mock.ResolveModule(plugin.BuiltInSource, mod))
			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_filterCommands(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		cmds  []plugin.ResolvedCommand
		lvl   HiddenLevel
		funcs []HideFunc

		expect []plugin.ResolvedCommand
	}{
		{
			name: "HiddenList",
			cmds: []plugin.ResolvedCommand{
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "abc", Hidden: true,
				}),
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "def", Hidden: false,
				}),
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "ghi", Hidden: false,
				}),
			},
			lvl: HideList,
			funcs: []HideFunc{
				CheckHidden(HideList),
				func(cmd plugin.ResolvedCommand, _ *state.State, _ *plugin.Context) HiddenLevel {
					if cmd.Name() == "def" {
						return Hide
					}

					return Show
				},
			},
			expect: []plugin.ResolvedCommand{
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "abc", Hidden: true,
				}),
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "ghi", Hidden: false,
				}),
			},
		},
		{
			name: "Show",
			cmds: []plugin.ResolvedCommand{
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "abc", Hidden: true,
				}),
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "def", Hidden: false,
				}),
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "ghi", Hidden: false,
				}),
			},
			lvl: Show,
			funcs: []HideFunc{
				CheckHidden(HideList),
				func(cmd plugin.ResolvedCommand, _ *state.State, _ *plugin.Context) HiddenLevel {
					if cmd.Name() == "def" {
						return Hide
					}

					return Show
				},
			},
			expect: []plugin.ResolvedCommand{
				mock.ResolveCommand(plugin.BuiltInSource, mock.Command{
					Name: "ghi", Hidden: false,
				}),
			},
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			h := &Help{Options: Options{HideFuncs: c.funcs}}

			actual := h.filterCommands(nil, nil, c.lvl, c.cmds...)
			assert.Len(t, actual, len(c.expect))

			for i, expect := range c.expect {
				assert.Equal(t, expect.Name(), actual[i].Name())
			}
		})
	}
}
