package help

import (
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

type HiddenLevel uint8

const (
	// Show is the HiddenLevel used if the command should be shown.
	Show HiddenLevel = iota
	// HideList is the HiddenLevel used if the command should be hidden from
	// command lists, i.e. from the general help page, and it's parent
	// module's help pages, if it has any.
	//
	// It will still be shown, if directly requesting it.
	HideList
	// Hide is the HiddenLevel used if the command should not be shown at all,
	// i.e. neither through command lists nor if asking directly for it.
	Hide
)

// Allows checks if this HiddenLevel would allow a plugin with the given
// HiddenLevel to be shown.
func (a HiddenLevel) Allows(b HiddenLevel) bool {
	return a >= b
}

type HideFunc func(plugin.ResolvedCommand, *state.State, *plugin.Context) HiddenLevel

// CheckHidden returns a HideFunc that returns the passed HiddenLevel, if the
// checked command is marked as Hidden.
func CheckHidden(lvl HiddenLevel) HideFunc {
	return func(cmd plugin.ResolvedCommand, s *state.State, ctx *plugin.Context) HiddenLevel {
		if cmd.IsHidden() {
			return lvl
		}

		return Show
	}
}

// CheckChannelTypes returns a HideFunc that checks if the commands
// plugin.ChannelTypes match those of the invoking channel, and returns the
// passed HiddenLevel if not.
//
// If an error occurs, it will be handled silently and Show will be returned.
func CheckChannelTypes(lvl HiddenLevel) HideFunc {
	return func(cmd plugin.ResolvedCommand, _ *state.State, ctx *plugin.Context) HiddenLevel {
		ok, err := cmd.ChannelTypes().Check(ctx)
		if err != nil {
			ctx.HandleErrorSilently(err)
			return Show
		}

		if !ok {
			return lvl
		}

		return Show
	}
}

// CheckRestrictions returns a HideFunc that returns the passed HiddenLevel, if
// the checked command is restricted.
// If plugin.ResolvedCommand.IsRestricted returns an error that fulfills
// errors.As(err, **plugin.RestrictionError) and that is fatal, lvl is
// returned.
// Otherwise, Show is returned.
func CheckRestrictions(lvl HiddenLevel) HideFunc {
	return func(cmd plugin.ResolvedCommand, s *state.State, ctx *plugin.Context) HiddenLevel {
		err := cmd.IsRestricted(s, ctx)
		if err != nil {
			var rerr *plugin.RestrictionError
			if errors.As(err, &rerr) {
				if rerr.Fatal {
					return lvl
				}

				return Show
			}

			ctx.HandleErrorSilently(err)
		}

		return Show
	}
}

// =============================================================================
// Utilities
// =====================================================================================

// calcCommandHiddenLevel checks the Help's HideFuncs and returns the highest
// HiddenLevel found.
func (h *Help) calcCommandHiddenLevel(s *state.State, ctx *plugin.Context, cmd plugin.ResolvedCommand) HiddenLevel {
	var lvl HiddenLevel

	for _, f := range h.HideFuncs {
		fLvl := f(cmd, s, ctx)
		if fLvl >= Hide {
			return Hide
		} else if fLvl > lvl && fLvl <= Hide {
			lvl = fLvl
		}
	}

	return lvl
}

// calcModuleHiddenLevel checks the Help's HideFuncs and returns the lowest
// HiddenLevel for all of mod's commands and subcommands.
func (h *Help) calcModuleHiddenLevel(s *state.State, ctx *plugin.Context, mod plugin.ResolvedModule) HiddenLevel {
	lvl := Hide

	for _, cmd := range mod.Commands() {
		cmdLvl := h.calcCommandHiddenLevel(s, ctx, cmd)
		if cmdLvl == Show {
			return Show
		} else if cmdLvl < lvl {
			lvl = cmdLvl
		}
	}

	for _, mod := range mod.Modules() {
		modLvl := h.calcModuleHiddenLevel(s, ctx, mod)
		if modLvl == Show {
			return Show
		} else if modLvl < lvl {
			lvl = modLvl
		}
	}

	return lvl
}

// filterCommands filters the passed commands and returns only those that have
// HiddenLevel smaller than max.
func (h *Help) filterCommands(
	s *state.State, ctx *plugin.Context, max HiddenLevel, cmds ...plugin.ResolvedCommand,
) []plugin.ResolvedCommand {
	filtered := make([]plugin.ResolvedCommand, 0, len(cmds))

	for _, cmd := range cmds {
		if max.Allows(h.calcCommandHiddenLevel(s, ctx, cmd)) {
			filtered = append(filtered, cmd)
		}
	}

	return filtered
}
