package help

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/internal/capbuilder"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

// =============================================================================
// Formatters
// =====================================================================================

func Test_formatCommand(t *testing.T) {
	t.Run("with short description", func(t *testing.T) {
		expect := "`mod kick` - kicks someone"

		cmd := &plugin.RegisteredCommand{
			Source: mock.Command{
				CommandMeta: mock.CommandMeta{ShortDescription: "kicks someone"},
			},
			ID: ".mod.kick",
		}

		b := capbuilder.New(100, 100)

		new(Help).formatCommand(b, cmd, nil)

		assert.Equal(t, expect, b.String())
	})

	t.Run("no short description", func(t *testing.T) {
		expect := "`mod kick`"

		cmd := &plugin.RegisteredCommand{
			Source: mock.Command{CommandMeta: new(mock.CommandMeta)},
			ID:     ".mod.kick",
		}

		b := capbuilder.New(100, 100)

		new(Help).formatCommand(b, cmd, nil)

		assert.Equal(t, expect, b.String())
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
			ID: ".mod.kick",
		},
		{
			Source: mock.Command{CommandMeta: new(mock.CommandMeta)},
			ID:     ".mod.abc",
		},
		{
			Source: mock.Command{
				CommandMeta: mock.CommandMeta{ShortDescription: "bans someone"},
			},
			ID: ".mod.ban",
		},
	}

	b := capbuilder.New(100, 100)

	new(Help).formatCommands(b, cmds, nil, new(plugin.Context), Show)

	assert.Equal(t, expect, b.String())
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
				ID: ".mod.kick",
			},
			{
				Source: mock.Command{CommandMeta: new(mock.CommandMeta)},
				ID:     ".mod.abc",
			},
			{
				Source: mock.Command{
					CommandMeta: mock.CommandMeta{ShortDescription: "bans someone"},
				},
				ID: ".mod.ban",
			},
		},
		Modules: []*plugin.RegisteredModule{
			{
				Commands: []*plugin.RegisteredCommand{
					{
						Source: mock.Command{
							CommandMeta: mock.CommandMeta{ShortDescription: "lists all infractions"},
						},
						ID: ".mod.infr.list",
					},
					{
						Source: mock.Command{
							CommandMeta: mock.CommandMeta{ShortDescription: "edits an infraction"},
						},
						ID: ".mod.infr.edit",
					},
					{
						Source: mock.Command{CommandMeta: new(mock.CommandMeta)},
						ID:     ".mod.infr.rm",
					},
				},
			},
			{
				Commands: []*plugin.RegisteredCommand{
					{
						Source: mock.Command{
							CommandMeta: mock.CommandMeta{ShortDescription: "turns the invite module on or off"},
						},
						ID: ".mod.invite.toggle",
					},
				},
			},
		},
	}

	b := capbuilder.New(500, 5000)

	new(Help).formatModule(b, mod, nil, new(plugin.Context), Show)

	assert.Equal(t, expect, b.String())
}
