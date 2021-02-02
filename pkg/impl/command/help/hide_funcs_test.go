package help

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestCheckHidden(t *testing.T) {
	testCases := []struct {
		name      string
		hiddenLvl HiddenLevel
		cmd       *plugin.RegisteredCommand

		expect HiddenLevel
	}{
		{
			name:      "hidden",
			hiddenLvl: HideList,
			cmd:       &plugin.RegisteredCommand{Hidden: true},
			expect:    HideList,
		},
		{
			name:      "not hidden",
			hiddenLvl: HideList,
			cmd:       &plugin.RegisteredCommand{Hidden: false},
			expect:    Show,
		},
		{
			name:      "level",
			hiddenLvl: HideList,
			cmd:       &plugin.RegisteredCommand{Hidden: true},
			expect:    HideList,
		},
		{
			name:      "level",
			hiddenLvl: Hide,
			cmd:       &plugin.RegisteredCommand{Hidden: true},
			expect:    Hide,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := CheckHidden(c.hiddenLvl)(c.cmd, nil, nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestCheckChannelTypes(t *testing.T) {
	successCases := []struct {
		name      string
		hiddenLvl HiddenLevel
		cmd       *plugin.RegisteredCommand
		ctx       *plugin.Context

		expect HiddenLevel
	}{
		{
			name:      "matching",
			hiddenLvl: HideList,
			cmd:       &plugin.RegisteredCommand{ChannelTypes: plugin.DirectMessages},
			ctx:       &plugin.Context{Message: discord.Message{GuildID: 0}},
			expect:    Show,
		},
		{
			name:      "not matching",
			hiddenLvl: HideList,
			cmd:       &plugin.RegisteredCommand{ChannelTypes: plugin.DirectMessages},
			ctx:       &plugin.Context{Message: discord.Message{GuildID: 123}},
			expect:    HideList,
		},
		{
			name:      "level",
			hiddenLvl: HideList,
			cmd:       &plugin.RegisteredCommand{ChannelTypes: plugin.DirectMessages},
			ctx:       &plugin.Context{Message: discord.Message{GuildID: 123}},
			expect:    HideList,
		},
		{
			name:      "level",
			hiddenLvl: Hide,
			cmd:       &plugin.RegisteredCommand{ChannelTypes: plugin.DirectMessages},
			ctx:       &plugin.Context{Message: discord.Message{GuildID: 123}},
			expect:    Hide,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				actual := CheckChannelTypes(c.hiddenLvl)(c.cmd, nil, c.ctx)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		expect := Show

		cmd := &plugin.RegisteredCommand{ChannelTypes: plugin.GuildTextChannels}

		channelError := errors.New("abc")

		errHandler := mock.NewErrorHandler(t).
			ExpectSilentError(channelError)
		defer errHandler.Eval()

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
	testCases := []struct {
		name      string
		hiddenLvl HiddenLevel
		cmd       *plugin.RegisteredCommand

		expect HiddenLevel
	}{
		{
			name:      "restricted",
			hiddenLvl: HideList,
			cmd: plugin.NewRegisteredCommandWithParent(nil,
				mock.RestrictionFunc(plugin.DefaultRestrictionError)),
			expect: HideList,
		},
		{
			name:      "not restricted",
			hiddenLvl: HideList,
			cmd:       new(plugin.RegisteredCommand),
			expect:    Show,
		},
		{
			name:      "level",
			hiddenLvl: HideList,
			cmd: plugin.NewRegisteredCommandWithParent(nil,
				mock.RestrictionFunc(plugin.DefaultRestrictionError)),
			expect: HideList,
		},
		{
			name:      "level",
			hiddenLvl: Hide,
			cmd: plugin.NewRegisteredCommandWithParent(nil,
				mock.RestrictionFunc(plugin.DefaultRestrictionError)),
			expect: Hide,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := CheckRestrictions(c.hiddenLvl)(c.cmd, nil, nil)
			assert.Equal(t, c.expect, actual)
		})
	}
}

// =============================================================================
// Utilities
// =====================================================================================

func Test_checkHideFuncs(t *testing.T) {
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
		t.Run(c.name, func(t *testing.T) {
			actual := checkHideFuncs(nil, nil, nil, c.funcs...)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_filterCommands(t *testing.T) {
	testCases := []struct {
		name  string
		cmds  []*plugin.RegisteredCommand
		lvl   HiddenLevel
		funcs []HideFunc

		expect []*plugin.RegisteredCommand
	}{
		{
			name: "level",
			cmds: []*plugin.RegisteredCommand{
				{Name: "abc", Hidden: true},
				{Name: "def", Hidden: false},
				{Name: "ghi", Hidden: false},
			},
			lvl: HideList,
			funcs: []HideFunc{
				CheckHidden(HideList),
				func(cmd *plugin.RegisteredCommand, _ *state.State, _ *plugin.Context) HiddenLevel {
					if cmd.Name == "def" {
						return Hide
					}

					return Show
				},
			},
			expect: []*plugin.RegisteredCommand{
				{Name: "abc", Hidden: true},
				{Name: "ghi", Hidden: false},
			},
		},
		{
			name: "level",
			cmds: []*plugin.RegisteredCommand{
				{Name: "abc", Hidden: true},
				{Name: "def", Hidden: false},
				{Name: "ghi", Hidden: false},
			},
			lvl: Show,
			funcs: []HideFunc{
				CheckHidden(HideList),
				func(cmd *plugin.RegisteredCommand, _ *state.State, _ *plugin.Context) HiddenLevel {
					if cmd.Name == "def" {
						return Hide
					}

					return Show
				},
			},
			expect: []*plugin.RegisteredCommand{
				{Name: "ghi", Hidden: false},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := filterCommands(c.cmds, nil, nil, c.lvl, c.funcs...)
			assert.Equal(t, c.expect, actual)
		})
	}
}
