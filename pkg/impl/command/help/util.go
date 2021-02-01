package help

import (
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// checkHideFuncs checks the passed HideFuncs and returns the highest
// HiddenLevel found.
func checkHideFuncs(cmd *plugin.RegisteredCommand, s *state.State, ctx *plugin.Context, f ...HideFunc) HiddenLevel {
	var lvl HiddenLevel

	for _, f := range f {
		lvl2 := f(cmd, s, ctx)
		if lvl2 >= Hide {
			return Hide
		} else if lvl2 > lvl {
			lvl = lvl2
		}
	}

	return lvl
}

func filterCommands(
	cmds []*plugin.RegisteredCommand, s *state.State, ctx *plugin.Context, lvl HiddenLevel, f ...HideFunc,
) []*plugin.RegisteredCommand {
	filtered := make([]*plugin.RegisteredCommand, 0, len(cmds))

	for _, cmd := range cmds {
		if checkHideFuncs(cmd, s, ctx, f...) <= lvl {
			filtered = append(filtered, cmd)
		}
	}

	return filtered
}

// =============================================================================
// Formatters
// =====================================================================================

func (h *Help) formatCommand(b *cappedBuilder, cmd *plugin.RegisteredCommand, l *i18n.Localizer) {
	b.writeRune('`')
	b.writeString(cmd.Identifier.AsInvoke())
	b.writeRune('`')

	if desc := cmd.ShortDescription(l); len(desc) > 0 {
		b.writeString(" - ")
		b.writeString(desc)
	}
}

// formatCommands writes a list of all commands that suffice the passed
// HiddenLevel to the passed *cappedBuilder, until there are no more commands,
// or the cappedBuilder's chunk size is reached.
func (h *Help) formatCommands(
	b *cappedBuilder, cmds []*plugin.RegisteredCommand, s *state.State, ctx *plugin.Context, lvl HiddenLevel,
) {
	cmds = filterCommands(cmds, s, ctx, lvl, h.HideFuncs...)

	for i, cmd := range cmds {
		if i > 0 {
			b.writeRune('\n')
		}

		h.formatCommand(b, cmd, ctx.Localizer)

		if b.rem() == 0 {
			return
		}
	}
}

func (h *Help) formatModule(
	b *cappedBuilder, mod *plugin.RegisteredModule, s *state.State, ctx *plugin.Context, lvl HiddenLevel,
) {
	if b.b.Len() > 0 {
		b.writeRune('\n')
	}

	h.formatCommands(b, mod.Commands, s, ctx, lvl)

	for _, smod := range mod.Modules {
		h.formatModule(b, smod, s, ctx, lvl)
	}
}

// =============================================================================
// cappedBuilder
// =====================================================================================

type cappedBuilder struct {
	used, cap int

	b *strings.Builder
}

func newCappedBuilder(totalCap, chunkCap int) *cappedBuilder {
	if totalCap < chunkCap {
		chunkCap = totalCap
	}

	b := new(strings.Builder)
	b.Grow(chunkCap)

	return &cappedBuilder{cap: totalCap, b: b}
}

func (b *cappedBuilder) writeRune(r rune) {
	if b.used < b.cap && b.b.Len() < b.b.Cap() {
		b.b.WriteRune(r)
		b.used++
	}
}

func (b *cappedBuilder) writeString(s string) {
	if b.used+len(s) < b.cap && b.b.Len()+len(s) < b.b.Cap() {
		b.b.WriteString(s)
		b.used += len(s)
	} else if b.used <= b.cap || b.b.Len() <= b.b.Cap() {
		end := b.cap - b.used
		if b.b.Cap()-b.b.Len() < end {
			end = b.b.Cap() - b.b.Len()
		}

		b.b.WriteString(s[:end])
		b.used += end
	}
}

func (b *cappedBuilder) use(n int) {
	b.used += n
}

func (b *cappedBuilder) string() string {
	return b.b.String()
}

func (b *cappedBuilder) reset(chunkCap int) {
	b.b.Reset()
	b.b.Grow(chunkCap)
}

func (b *cappedBuilder) rem() int {
	rem := b.cap - b.used
	if chunkRem := b.b.Cap() - b.b.Len(); chunkRem < rem {
		rem = chunkRem
	}

	return rem
}
