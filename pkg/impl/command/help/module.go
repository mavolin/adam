package help

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

func (h *Help) module(s *state.State, ctx *plugin.Context, mod *plugin.RegisteredModule) (discord.Embed, error) {
	e, err := newModuleEmbed(mod, ctx.Localizer)
	if err != nil {
		return discord.Embed{}, err
	}

	maxMods := 25 - len(e.Fields)

	b := newCappedBuilder(embedutil.MaxChars-embedutil.CountChars(e), 1024)

	if f := h.commands(b, s, ctx, mod.Commands); len(f.Name) > 0 {
		e.Fields = append(e.Fields, f)
		maxMods--
	}

	e.Fields = append(e.Fields, h.modules(b, s, ctx, mod.Modules, maxMods)...)

	if len(e.Fields) == 0 {
		return discord.Embed{}, plugin.NewArgumentErrorl(pluginNotFoundError)
	}

	return e, nil
}

func newModuleEmbed(mod *plugin.RegisteredModule, l *i18n.Localizer) (discord.Embed, error) {
	eb := BaseEmbed.Clone().
		WithSimpleTitlel(moduleTitle.
			WithPlaceholders(moduleTitlePlaceholders{
				Module: mod.Identifier.AsInvoke(),
			}))

	if desc := mod.LongDescription(l); len(desc) > 0 {
		eb.WithDescription(desc)
	}

	return eb.Build(l)
}
