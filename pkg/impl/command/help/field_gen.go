package help

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/capbuilder"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// Field Generation
// =====================================================================================

func (h *Help) genPrefixesField(b *capbuilder.CappedBuilder, ctx *plugin.Context) (discord.EmbedField, error) {
	b.Reset(1024)

	f := discord.EmbedField{Name: ctx.MustLocalize(allPrefixesFieldName)}

	self, err := ctx.Self()
	if err != nil {
		return discord.EmbedField{}, err
	}

	name := self.Nick
	if len(name) == 0 {
		name = self.User.Username
	}

	b.WriteString("`@")
	b.WriteString(name)
	b.WriteRune('`')

	for _, prefix := range ctx.Prefixes {
		b.WriteString(", `")
		b.WriteString(prefix)
		b.WriteRune('`')
	}

	f.Value = b.String()

	return f, nil
}

func (h *Help) genCommandsField(
	b *capbuilder.CappedBuilder, s *state.State, ctx *plugin.Context, cmds []*plugin.RegisteredCommand,
) (f discord.EmbedField) {
	b.Reset(1024)

	h.formatCommands(b, cmds, s, ctx, Show)
	if b.ChunkLen() == 0 {
		return
	}

	f.Name = ctx.MustLocalize(commandsFieldName)
	b.Use(len(f.Name))

	f.Value = b.String()
	return
}

func (h *Help) genModuleFields(
	b *capbuilder.CappedBuilder, s *state.State, ctx *plugin.Context, mods []*plugin.RegisteredModule, max int,
) []discord.EmbedField {
	fields := make([]discord.EmbedField, 0, max)

	for i := 0; i < len(mods) && i < max; i++ {
		mod := mods[i]
		b.Reset(1024)

		var f discord.EmbedField

		f.Name = ctx.MustLocalize(moduleTitle.
			WithPlaceholders(moduleTitlePlaceholders{Module: mod.Identifier.AsInvoke()}))
		b.Use(len(f.Name))

		if b.Rem() < 10+len(f.Name) {
			return fields
		}

		h.formatModule(b, mod, s, ctx, Show)

		if b.ChunkLen() == 0 { // hidden, skip
			b.Use(-len(f.Name))
			continue
		}

		f.Value = b.String()

		fields = append(fields, f)
	}

	return fields
}

// =============================================================================
// Formatters
// =====================================================================================

func (h *Help) formatCommand(b *capbuilder.CappedBuilder, cmd *plugin.RegisteredCommand, l *i18n.Localizer) {
	b.WriteRune('`')
	b.WriteString(cmd.Identifier.AsInvoke())
	b.WriteRune('`')

	if desc := cmd.ShortDescription(l); len(desc) > 0 {
		b.WriteString(" - ")
		b.WriteString(desc)
	}
}

// formatCommands Writes a list of all commands that suffice the passed
// HiddenLevel to the passed *capbuilder.CappedBuilder, until there are no more
// commands, or the capbuilder.CappedBuilder's chunk size is reached.
func (h *Help) formatCommands(
	b *capbuilder.CappedBuilder, cmds []*plugin.RegisteredCommand, s *state.State, ctx *plugin.Context, lvl HiddenLevel,
) {
	cmds = filterCommands(cmds, s, ctx, lvl, h.HideFuncs...)

	for i, cmd := range cmds {
		if i > 0 {
			b.WriteRune('\n')
		}

		h.formatCommand(b, cmd, ctx.Localizer)

		if b.Rem() == 0 {
			return
		}
	}
}

func (h *Help) formatModule(
	b *capbuilder.CappedBuilder, mod *plugin.RegisteredModule, s *state.State, ctx *plugin.Context, lvl HiddenLevel,
) {
	if b.ChunkLen() > 0 {
		b.WriteRune('\n')
	}

	h.formatCommands(b, mod.Commands, s, ctx, lvl)

	for _, smod := range mod.Modules {
		h.formatModule(b, smod, s, ctx, lvl)
	}
}
