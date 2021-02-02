package help

import (
	"strings"

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

func (h *Help) genAliases(b *strings.Builder, cmd *plugin.RegisteredCommand, l *i18n.Localizer) *discord.EmbedField {
	if len(cmd.Aliases) == 0 {
		return nil
	}

	b.Reset()

	for i, alias := range cmd.Aliases {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteRune('`')
		b.WriteString(alias)
		b.WriteRune('`')
	}

	return &discord.EmbedField{
		Name:  l.MustLocalize(aliasesFieldName),
		Value: b.String(),
	}
}

func (h *Help) genUsages(b *strings.Builder, ctx *plugin.Context, cmd *plugin.RegisteredCommand) []discord.EmbedField {
	if cmd.Args == nil { // special case, command accepts no arguments
		return []discord.EmbedField{
			{
				Name:  ctx.MustLocalize(usageFieldNameSingle),
				Value: "```" + cmd.Identifier.AsInvoke() + "```",
			},
		}
	}

	infoer, ok := cmd.Args.(plugin.ArgsInfoer)
	if !ok {
		return nil
	}

	infos := infoer.Info(ctx.Localizer)
	if len(infos) == 0 {
		return nil
	}

	fields := make([]discord.EmbedField, 0, 3*len(infos)+1)

	for i, info := range infos {
		var n int
		if len(infos) > 1 {
			n = i + 1
		}

		fields = append(fields, h.genUsage(b, cmd, info, n, ctx.Localizer))

		if args := h.genArguments(b, info, ctx.Localizer); args != nil {
			fields = append(fields, *args)
		}

		if flags := h.genFlags(b, info, ctx.Localizer); flags != nil {
			fields = append(fields, *flags)
		}
	}

	return fields
}

func (h *Help) genUsage(
	b *strings.Builder, cmd *plugin.RegisteredCommand, info plugin.ArgsInfo, n int, l *i18n.Localizer,
) (usage discord.EmbedField) {
	if n == 0 {
		usage.Name = l.MustLocalize(usageFieldNameSingle)
	} else {
		usage.Name = l.MustLocalize(usageFieldNameMulti.
			WithPlaceholders(usageFieldNameMultiPlaceholders{
				Num: n,
			}))
	}

	b.Reset()

	b.WriteString("```")
	b.WriteString(cmd.Identifier.AsInvoke())

	if len(info.Prefix) > 0 {
		b.WriteRune(' ')
		b.WriteString(info.Prefix)
	}

	if info.ArgsFormatter != nil {
		if argUsage := info.ArgsFormatter(h.ArgFormatter); len(argUsage) > 0 {
			b.WriteRune(' ')
			b.WriteString(argUsage)
		}
	}

	b.WriteString("```")

	usage.Value = b.String()
	return
}

func (h *Help) genArguments(b *strings.Builder, info plugin.ArgsInfo, l *i18n.Localizer) *discord.EmbedField {
	if len(info.Required) == 0 && len(info.Optional) == 0 {
		return nil
	}

	b.Reset()

	for i, arg := range info.Required {
		if len(arg.Description) == 0 {
			continue
		}

		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteRune('`')
		b.WriteString(arg.Name)

		if !strings.EqualFold(arg.Name, arg.Type.Name) ||
			(info.Variadic && len(info.Optional) == 0 && i == len(info.Required)-1) {
			b.WriteString(" (")
			b.WriteString(arg.Type.Name)

			if info.Variadic && len(info.Optional) == 0 && i == len(info.Required)-1 {
				b.WriteRune('+')
			}

			b.WriteRune(')')
		}

		b.WriteString("` - ")
		b.WriteString(arg.Description)
	}

	for i, arg := range info.Optional {
		if len(arg.Description) == 0 {
			continue
		}

		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteRune('`')
		b.WriteString(arg.Name)

		if !strings.EqualFold(arg.Name, arg.Type.Name) ||
			(info.Variadic && i == len(info.Optional)-1) {
			b.WriteString(" (")
			b.WriteString(arg.Type.Name)

			if info.Variadic && i == len(info.Optional)-1 {
				b.WriteRune('+')
			}

			b.WriteRune(')')
		}

		b.WriteString("` - ")
		b.WriteString(arg.Description)
	}

	if len(b.String()) == 0 {
		return nil
	}

	return &discord.EmbedField{
		Name:  l.MustLocalize(argumentsFieldName),
		Value: b.String(),
	}
}

func (h *Help) genFlags(b *strings.Builder, info plugin.ArgsInfo, l *i18n.Localizer) *discord.EmbedField {
	if len(info.Flags) == 0 {
		return nil
	}

	b.Reset()

	for i, flag := range info.Flags {
		if i > 0 {
			b.WriteRune('\n')
		}

		b.WriteRune('`')
		b.WriteString(info.FlagFormatter(flag.Name))

		for _, alias := range flag.Aliases {
			b.WriteString(", ")
			b.WriteString(info.FlagFormatter(alias))
		}

		b.WriteString(" (")
		b.WriteString(flag.Type.Name)

		if flag.Multi {
			b.WriteRune('+')
		}

		b.WriteString(")`")

		if len(flag.Description) > 0 {
			b.WriteString(" - ")
			b.WriteString(flag.Description)
		}
	}

	return &discord.EmbedField{
		Name:  l.MustLocalize(flagsFieldName),
		Value: b.String(),
	}
}

func (h *Help) genExamples(b *strings.Builder, cmd *plugin.RegisteredCommand, l *i18n.Localizer) *discord.EmbedField {
	examples := cmd.Examples(l)
	if len(examples) == 0 {
		return nil
	}

	b.Reset()

	for _, e := range examples {
		b.WriteString("```")
		b.WriteString(e)
		b.WriteString("```")
	}

	return &discord.EmbedField{
		Name:  l.MustLocalize(examplesFieldName),
		Value: b.String(),
	}
}

// =============================================================================
// Formatters
// =====================================================================================

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

func (h *Help) formatCommand(b *capbuilder.CappedBuilder, cmd *plugin.RegisteredCommand, l *i18n.Localizer) {
	b.WriteRune('`')
	b.WriteString(cmd.Identifier.AsInvoke())
	b.WriteRune('`')

	if desc := cmd.ShortDescription(l); len(desc) > 0 {
		b.WriteString(" - ")
		b.WriteString(desc)
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
