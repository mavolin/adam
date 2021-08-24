package help

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/internal/capbuilder"
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
	b *capbuilder.CappedBuilder, s *state.State, ctx *plugin.Context, cmds []plugin.ResolvedCommand,
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
	b *capbuilder.CappedBuilder, s *state.State, ctx *plugin.Context, mods []plugin.ResolvedModule, max int,
) []discord.EmbedField {
	fields := make([]discord.EmbedField, 0, max)

	for i := 0; i < len(mods) && i < max; i++ {
		mod := mods[i]
		b.Reset(1024)

		var f discord.EmbedField

		f.Name = ctx.MustLocalize(moduleTitle.
			WithPlaceholders(moduleTitlePlaceholders{Module: mod.ID().AsInvoke()}))
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

func (h *Help) genAliasesField(
	b *strings.Builder, ctx *plugin.Context, cmd plugin.ResolvedCommand,
) *discord.EmbedField {
	if len(cmd.Aliases()) == 0 {
		return nil
	}

	b.Reset()

	for i, alias := range cmd.Aliases() {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteRune('`')
		b.WriteString(alias)
		b.WriteRune('`')
	}

	return &discord.EmbedField{
		Name:  ctx.MustLocalize(aliasesFieldName),
		Value: b.String(),
	}
}

func (h *Help) genUsage(
	b *strings.Builder, ctx *plugin.Context, cmd plugin.ResolvedCommand,
) (usage discord.EmbedField) {
	usage.Name = ctx.MustLocalize(usageFieldNameSingle)

	b.Reset()

	b.WriteString("```")
	b.WriteString(cmd.ID().AsInvoke())

	if cmd.Args() == nil {
		b.WriteString("```")
		usage.Value = b.String()
		return usage
	}

	b.WriteRune(' ')

	var (
		requiredArgs = cmd.Args().GetRequiredArgs()
		optionalArgs = cmd.Args().GetOptionalArgs()

		usageArgs = make([]string, 0, len(requiredArgs)+len(optionalArgs))
	)

	for i, arg := range requiredArgs {
		name := arg.GetName(ctx.Localizer)
		typeName := arg.GetType().GetName(ctx.Localizer)
		variadic := cmd.Args().IsVariadic() && i == len(requiredArgs)-1 && len(optionalArgs) == 0

		usageArgs = append(usageArgs, h.ArgFormatter(name, typeName, false, variadic))
	}

	for i, arg := range optionalArgs {
		name := arg.GetName(ctx.Localizer)
		typeName := arg.GetType().GetName(ctx.Localizer)
		variadic := cmd.Args().IsVariadic() && i == len(optionalArgs)-1

		usageArgs = append(usageArgs, h.ArgFormatter(name, typeName, true, variadic))
	}

	b.WriteString(cmd.ArgParser().FormatUsage(cmd.Args(), usageArgs))
	b.WriteString("```")

	usage.Value = b.String()
	return usage
}

//nolint:funlen,gocognit
func (h *Help) genArguments(
	b *strings.Builder, ctx *plugin.Context, cmd plugin.ResolvedCommand,
) *discord.EmbedField {
	if cmd.Args() == nil {
		return nil
	}

	requiredArgs := cmd.Args().GetRequiredArgs()
	optionalArgs := cmd.Args().GetOptionalArgs()
	if len(requiredArgs) == 0 && len(optionalArgs) == 0 {
		return nil
	}

	b.Reset()

	for i, arg := range requiredArgs {
		desc := arg.GetDescription(ctx.Localizer)

		if len(desc) == 0 {
			continue
		}

		name := arg.GetName(ctx.Localizer)
		typeName := arg.GetType().GetName(ctx.Localizer)

		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteRune('`')

		b.WriteString(name)

		variadic := cmd.Args().IsVariadic() && i == len(requiredArgs)-1 && len(optionalArgs) == 0

		if (!strings.EqualFold(name, typeName) || variadic) && len(typeName) > 0 {
			b.WriteString(" (")
			b.WriteString(typeName)

			if variadic {
				b.WriteRune('+')
			}

			b.WriteRune(')')
		}

		b.WriteString("` - ")
		b.WriteString(desc)
	}

	for i, arg := range optionalArgs {
		desc := arg.GetDescription(ctx.Localizer)

		if len(desc) == 0 {
			continue
		}

		name := arg.GetName(ctx.Localizer)
		typeName := arg.GetType().GetName(ctx.Localizer)

		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteRune('`')
		b.WriteString(name)

		variadic := cmd.Args().IsVariadic() && i == len(optionalArgs)-1

		if (!strings.EqualFold(name, typeName) || variadic) && len(typeName) > 0 {
			b.WriteString(" (")
			b.WriteString(typeName)

			if variadic {
				b.WriteRune('+')
			}

			b.WriteRune(')')
		}

		b.WriteString("` - ")
		b.WriteString(desc)
	}

	if len(b.String()) == 0 {
		return nil
	}

	return &discord.EmbedField{
		Name:  ctx.MustLocalize(argumentsFieldName),
		Value: b.String(),
	}
}

func (h *Help) genFlags(b *strings.Builder, ctx *plugin.Context, cmd plugin.ResolvedCommand) *discord.EmbedField {
	if cmd.Args() == nil {
		return nil
	}

	flags := cmd.Args().GetFlags()
	if len(flags) == 0 {
		return nil
	}

	b.Reset()

	for i, flag := range flags {
		if i > 0 {
			b.WriteRune('\n')
		}

		b.WriteRune('`')
		b.WriteString(cmd.ArgParser().FormatFlag(flag.GetName()))

		for _, alias := range flag.GetAliases() {
			b.WriteString(", ")
			b.WriteString(cmd.ArgParser().FormatFlag(alias))
		}

		typeName := flag.GetType().GetName(ctx.Localizer)
		if len(typeName) > 0 {
			b.WriteString(" (")
			b.WriteString(typeName)

			if flag.IsMulti() {
				b.WriteRune('+')
			}

			b.WriteString(")`")
		}

		desc := flag.GetDescription(ctx.Localizer)
		if len(desc) > 0 {
			b.WriteString(" - ")
			b.WriteString(desc)
		}
	}

	return &discord.EmbedField{
		Name:  ctx.MustLocalize(flagsFieldName),
		Value: b.String(),
	}
}

func (h *Help) genExamples(b *strings.Builder, ctx *plugin.Context, cmd plugin.ResolvedCommand) *discord.EmbedField {
	examples := cmd.Examples(ctx.Localizer)
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
		Name:  ctx.MustLocalize(examplesFieldName),
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
	b *capbuilder.CappedBuilder, cmds []plugin.ResolvedCommand, s *state.State, ctx *plugin.Context,
	maxLvl HiddenLevel,
) {
	cmds = h.filterCommands(s, ctx, maxLvl, cmds...)

	for i, cmd := range cmds {
		if i > 0 {
			b.WriteRune('\n')
		}

		b.WriteRune('`')
		b.WriteString(cmd.ID().AsInvoke())
		b.WriteRune('`')

		if desc := cmd.ShortDescription(ctx.Localizer); len(desc) > 0 {
			b.WriteString(" - ")
			b.WriteString(desc)
		}

		if b.Rem() == 0 {
			return
		}
	}
}

func (h *Help) formatModule(
	b *capbuilder.CappedBuilder, mod plugin.ResolvedModule, s *state.State, ctx *plugin.Context, lvl HiddenLevel,
) {
	if b.ChunkLen() > 0 {
		b.WriteRune('\n')
	}

	h.formatCommands(b, mod.Commands(), s, ctx, lvl)

	for _, smod := range mod.Modules() {
		h.formatModule(b, smod, s, ctx, lvl)
	}
}
