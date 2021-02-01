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

	maxMods := 25 - len(e.Fields)

	if ctx.GuildID > 0 && !h.NoPrefix {
		prefixes, err := h.allPrefixes(b, ctx)
		if err != nil {
			return discord.Embed{}, err
		}

		e.Fields = append([]discord.EmbedField{prefixes}, e.Fields...)
		maxMods--
	}

	if f := h.commands(b, s, ctx, ctx.Commands()); len(f.Name) > 0 {
		e.Fields = append(e.Fields, f)
		maxMods--
	}

	e.Fields = append(e.Fields, h.modules(b, s, ctx, ctx.Modules(), maxMods)...)
	return e, nil
}

func newAllEmbed(ctx *plugin.Context) (discord.Embed, error) {
	e := BaseEmbed.Clone().
		WithSimpleTitlel(allTitle)

	if ctx.GuildID == 0 {
		e.WithDescriptionl(allDescriptionDM)
	} else {
		e.WithDescriptionl(allDescriptionGuild)
	}

	return e.Build(ctx.Localizer)
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
