package help

import (
	"strings"
	"testing"

	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

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

func Test_formatCommand(t *testing.T) {
	t.Run("with short description", func(t *testing.T) {
		expect := "`mod kick` - kicks someone"

		cmd := &plugin.RegisteredCommand{
			Source: mock.Command{
				CommandMeta: mock.CommandMeta{ShortDescription: "kicks someone"},
			},
			Identifier: ".mod.kick",
		}

		b := newCappedBuilder(100, 100)

		new(Help).formatCommand(b, cmd, nil)

		assert.Equal(t, expect, b.string())
	})

	t.Run("no short description", func(t *testing.T) {
		expect := "`mod kick`"

		cmd := &plugin.RegisteredCommand{
			Source:     mock.Command{CommandMeta: new(mock.CommandMeta)},
			Identifier: ".mod.kick",
		}

		b := newCappedBuilder(100, 100)

		new(Help).formatCommand(b, cmd, nil)

		assert.Equal(t, expect, b.string())
	})
}

func Test_formatCommands(t *testing.T) {
	expect := "`mod kick` - kicks someone\n" +
		"`mod abc`\n" +
		"`mod ban` - bans someone"

	cmds := []*plugin.RegisteredCommand{
		{
			Source: mock.Command{
				CommandMeta: mock.CommandMeta{ShortDescription: "kicks someone"},
			},
			Identifier: ".mod.kick",
		},
		{
			Source:     mock.Command{CommandMeta: new(mock.CommandMeta)},
			Identifier: ".mod.abc",
		},
		{
			Source: mock.Command{
				CommandMeta: mock.CommandMeta{ShortDescription: "bans someone"},
			},
			Identifier: ".mod.ban",
		},
	}

	b := newCappedBuilder(100, 100)

	new(Help).formatCommands(b, cmds, nil, new(plugin.Context), Show)

	assert.Equal(t, expect, b.string())
}

func Test_formatModule(t *testing.T) {
	expect := "`mod kick` - kicks someone\n" +
		"`mod abc`\n" +
		"`mod ban` - bans someone\n" +
		"`mod infr list` - lists all infractions\n" +
		"`mod infr edit` - edits an infraction\n" +
		"`mod infr rm`\n" +
		"`mod invite toggle` - turns the invite module on or off"

	mod := &plugin.RegisteredModule{
		Commands: []*plugin.RegisteredCommand{
			{
				Source: mock.Command{
					CommandMeta: mock.CommandMeta{ShortDescription: "kicks someone"},
				},
				Identifier: ".mod.kick",
			},
			{
				Source:     mock.Command{CommandMeta: new(mock.CommandMeta)},
				Identifier: ".mod.abc",
			},
			{
				Source: mock.Command{
					CommandMeta: mock.CommandMeta{ShortDescription: "bans someone"},
				},
				Identifier: ".mod.ban",
			},
		},
		Modules: []*plugin.RegisteredModule{
			{
				Commands: []*plugin.RegisteredCommand{
					{
						Source: mock.Command{
							CommandMeta: mock.CommandMeta{ShortDescription: "lists all infractions"},
						},
						Identifier: ".mod.infr.list",
					},
					{
						Source: mock.Command{
							CommandMeta: mock.CommandMeta{ShortDescription: "edits an infraction"},
						},
						Identifier: ".mod.infr.edit",
					},
					{
						Source:     mock.Command{CommandMeta: new(mock.CommandMeta)},
						Identifier: ".mod.infr.rm",
					},
				},
			},
			{
				Commands: []*plugin.RegisteredCommand{
					{
						Source: mock.Command{
							CommandMeta: mock.CommandMeta{ShortDescription: "turns the invite module on or off"},
						},
						Identifier: ".mod.invite.toggle",
					},
				},
			},
		},
	}

	b := newCappedBuilder(500, 5000)

	new(Help).formatModule(b, mod, nil, new(plugin.Context), Show)

	assert.Equal(t, expect, b.string())
}

func Test_cappedBuilder_writeRune(t *testing.T) {
	t.Run("chunk limit", func(t *testing.T) {
		totalCap := 20
		chunkCap := 5

		b := newCappedBuilder(totalCap, chunkCap)

		for i := 0; i < totalCap; i++ {
			b.writeRune('a')
		}

		expect := strings.Repeat("a", chunkCap)

		assert.Equal(t, b.used, chunkCap)
		assert.Equal(t, expect, b.string())
	})

	t.Run("global limit", func(t *testing.T) {
		totalCap := 10
		chunkCap := 6

		b := newCappedBuilder(totalCap, chunkCap)

		for i := 0; i < chunkCap; i++ {
			b.writeRune('a')
		}

		b.reset(chunkCap)

		for i := 0; i < chunkCap; i++ {
			b.writeRune('a')
		}

		expect := strings.Repeat("a", totalCap-chunkCap)

		assert.Equal(t, b.used, b.cap)
		assert.Equal(t, expect, b.string())
	})
}

func Test_cappedBuilder_writeString(t *testing.T) {
	t.Run("chunk limit", func(t *testing.T) {
		t.Run("full write", func(t *testing.T) {
			totalCap := 20
			chunkCap := 6

			b := newCappedBuilder(totalCap, chunkCap)

			for i := 0; i < totalCap; i++ {
				b.writeString("ab")
			}

			expect := strings.Repeat("ab", chunkCap/2)

			assert.Equal(t, b.used, chunkCap)
			assert.Equal(t, expect, b.string())
		})

		t.Run("partial write", func(t *testing.T) {
			totalCap := 20
			chunkCap := 5

			b := newCappedBuilder(totalCap, chunkCap)

			for i := 0; i < totalCap; i++ {
				b.writeString("ab")
			}

			expect := "ababa"

			assert.Equal(t, b.used, chunkCap)
			assert.Equal(t, expect, b.string())
		})
	})

	t.Run("global limit", func(t *testing.T) {
		t.Run("full write", func(t *testing.T) {
			totalCap := 10
			chunkCap := 6

			b := newCappedBuilder(totalCap, chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.writeString("ab")
			}

			b.reset(chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.writeString("ab")
			}

			expect := strings.Repeat("ab", (totalCap-chunkCap)/2)

			assert.Equal(t, b.used, b.cap)
			assert.Equal(t, expect, b.string())
		})

		t.Run("partial write", func(t *testing.T) {
			totalCap := 9
			chunkCap := 6

			b := newCappedBuilder(totalCap, chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.writeString("ab")
			}

			b.reset(chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.writeString("ab")
			}

			expect := "aba"

			assert.Equal(t, b.used, b.cap)
			assert.Equal(t, expect, b.string())
		})
	})
}

func Test_cappedBuilder_reset(t *testing.T) {
	chunkCap := 3

	b := newCappedBuilder(10, chunkCap)
	assert.Equal(t, chunkCap, b.b.Cap())
	assert.Equal(t, 0, b.b.Len())

	b.writeRune('a')
	assert.Equal(t, 1, b.b.Len())

	b.reset(chunkCap + 1)
	assert.Equal(t, chunkCap+1, b.b.Cap())
	assert.Equal(t, 0, b.b.Len())

	b.writeRune('a')
	assert.Equal(t, 1, b.b.Len())
}

func Test_cappedBuilder_rem(t *testing.T) {
	t.Run("chunk", func(t *testing.T) {
		chunkCap := 5

		b := newCappedBuilder(10, chunkCap)
		b.writeRune('a')

		assert.Equal(t, chunkCap-1, b.rem())
	})

	t.Run("total", func(t *testing.T) {
		totalCap := 10
		chunkCap := 7

		b := newCappedBuilder(totalCap, chunkCap)
		b.writeString(strings.Repeat("a", chunkCap))
		b.reset(chunkCap)

		assert.Equal(t, totalCap-chunkCap, b.rem())
	})
}
