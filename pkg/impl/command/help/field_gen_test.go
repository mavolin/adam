package help

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/internal/capbuilder"
	"github.com/mavolin/adam/pkg/impl/module"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

// =============================================================================
// Formatters
// =====================================================================================

func Test_formatCommands(t *testing.T) {
	t.Parallel()

	expect := "`mod abc`\n" +
		"`mod ban` - bans someone\n" +
		"`mod kick` - kicks someone"

	mod := module.New(&module.Meta{Name: "mod"})
	mod.AddCommand(mock.Command{Name: "abc"})
	mod.AddCommand(mock.Command{
		Name:             "ban",
		ShortDescription: "bans someone",
	})
	mod.AddCommand(mock.Command{
		Name:             "kick",
		ShortDescription: "kicks someone",
	})

	cmds := mock.ResolveModule(plugin.BuiltInSource, mod).Commands()

	b := capbuilder.New(100, 100)

	new(Help).formatCommands(b, cmds, nil, new(plugin.Context), Show)

	assert.Equal(t, expect, b.String())
}

func Test_formatModule(t *testing.T) {
	t.Parallel()

	expect := "`mod abc`\n" +
		"`mod ban` - bans someone\n" +
		"`mod kick` - kicks someone\n" +
		"`mod infr edit` - edits an infraction\n" +
		"`mod infr list` - lists all infractions\n" +
		"`mod infr rm`\n" +
		"`mod invite toggle` - turns the invite module on or off"

	mod := module.New(&module.Meta{Name: "mod"})
	mod.AddCommand(mock.Command{Name: "abc"})
	mod.AddCommand(mock.Command{
		Name:             "ban",
		ShortDescription: "bans someone",
	})
	mod.AddCommand(mock.Command{
		Name:             "kick",
		ShortDescription: "kicks someone",
	})

	infr := module.New(module.Meta{Name: "infr"})
	infr.AddCommand(mock.Command{
		Name:             "edit",
		ShortDescription: "edits an infraction",
	})
	infr.AddCommand(mock.Command{
		Name:             "list",
		ShortDescription: "lists all infractions",
	})
	infr.AddCommand(mock.Command{Name: "rm"})

	mod.AddModule(infr)

	invite := module.New(module.Meta{Name: "invite"})
	invite.AddCommand(mock.Command{
		Name:             "toggle",
		ShortDescription: "turns the invite module on or off",
	})

	mod.AddModule(invite)

	rmod := mock.ResolveModule(plugin.BuiltInSource, mod)
	b := capbuilder.New(500, 5000)

	new(Help).formatModule(b, rmod, nil, new(plugin.Context), Show)

	assert.Equal(t, expect, b.String())
}
