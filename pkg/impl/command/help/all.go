package help

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

func (h *Help) all(s *state.State, ctx *plugin.Context) (discord.Embed, error) {
	e, err := newAllEmbed(ctx)
	if err != nil {
		return discord.Embed{}, err
	}

	b := newCappedBuilder(embedutil.MaxChars-embedutil.CountChars(e), 1024)

	if ctx.GuildID > 0 && !h.NoPrefix {
		prefixes, err := h.allPrefixes(b, ctx)
		if err != nil {
			return discord.Embed{}, err
		}

		e.Fields = append([]discord.EmbedField{prefixes}, e.Fields...)
	}

	if f := h.allCommands(b, s, ctx); len(f.Name) > 0 {
		e.Fields = append(e.Fields, f)
	}

	e.Fields = append(e.Fields, h.allModules(b, s, ctx)...)

	return e, nil
}

func newAllEmbed(ctx *plugin.Context) (discord.Embed, error) {
	eb := BaseEmbed.Clone().
		WithSimpleTitlel(allTitle)

	if ctx.GuildID == 0 {
		eb.WithDescriptionl(allDescriptionDM)
	} else {
		eb.WithDescriptionl(allDescriptionGuild)
	}

	return eb.Build(ctx.Localizer)
}

func (h *Help) allPrefixes(b *cappedBuilder, ctx *plugin.Context) (discord.EmbedField, error) {
	b.reset(1024)

	f := discord.EmbedField{Name: ctx.MustLocalize(allPrefixesFieldName)}

	self, err := ctx.Self()
	if err != nil {
		return discord.EmbedField{}, err
	}

	name := self.Nick
	if len(name) == 0 {
		name = self.User.Username
	}

	b.writeString("`@")
	b.writeString(name)
	b.writeRune('`')

	for _, prefix := range ctx.Prefixes {
		b.writeString(", `")
		b.writeString(prefix)
		b.writeRune('`')
	}

	f.Value = b.string()

	return f, nil
}

func (h *Help) allCommands(b *cappedBuilder, s *state.State, ctx *plugin.Context) (f discord.EmbedField) {
	b.reset(1024)

	h.formatCommands(b, ctx.Commands(), s, ctx, Show)
	if b.b.Len() == 0 {
		return
	}

	f.Name = ctx.MustLocalize(commandsFieldName)
	b.use(len(f.Name))

	f.Value = b.string()
	return
}

func (h *Help) allModules(b *cappedBuilder, s *state.State, ctx *plugin.Context) []discord.EmbedField {
	max := 25
	if ctx.GuildID != 0 {
		max--
	}

	if len(ctx.Commands()) > 0 {
		max--
	}

	mods := ctx.Modules()
	if len(mods) > max {
		mods = mods[:max]
	}

	fields := make([]discord.EmbedField, 0, len(mods))

	for _, mod := range mods {
		b.reset(1024)

		var f discord.EmbedField

		f.Name = ctx.MustLocalize(allModuleFieldName.
			WithPlaceholders(allModuleFieldNamePlaceholders{Module: mod.Name}))
		b.use(len(f.Name))

		if b.rem() < 10+len(f.Name) {
			return fields
		}

		h.formatModule(b, mod, s, ctx, Show)

		if b.b.Len() == 0 { // hidden, skip
			b.use(-len(f.Name))
			continue
		}

		f.Value = b.string()

		fields = append(fields, f)
	}

	return fields
}
