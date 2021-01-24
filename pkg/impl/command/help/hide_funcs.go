package help

import (
	"github.com/mavolin/disstate/v3/pkg/state"

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
	// i.e. neither through command lists nor if asking direct for it.
	Hide
)

type HideFunc func(*plugin.RegisteredCommand, *state.State, *plugin.Context) HiddenLevel

// CheckHidden returns a HideFunc that returns the passed HiddenLevel, if the
// checked command is marked as Hidden.
func CheckHidden(lvl HiddenLevel) HideFunc {
	return func(cmd *plugin.RegisteredCommand, s *state.State, ctx *plugin.Context) HiddenLevel {
		if cmd.Hidden {
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
	return func(cmd *plugin.RegisteredCommand, _ *state.State, ctx *plugin.Context) HiddenLevel {
		ok, err := cmd.ChannelTypes.Check(ctx)
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

// CheckRestriction returns a HideFunc that returns the passed HiddenLevel, if
// the checked command is restricted.
// If plugin.RegisteredCommand.IsRestricted returns an error, for which
// errors.As(err, **plugin.RestrictionError) fails, the error will be
// handled silently, and Show will be returned.
// If the error fulfils errors.As for that case lvl will be returned.
func CheckRestrictions(lvl HiddenLevel) HideFunc {
	return func(cmd *plugin.RegisteredCommand, s *state.State, ctx *plugin.Context) HiddenLevel {
		err := cmd.IsRestricted(s, ctx)
		if err != nil {
			var rerr *plugin.RestrictionError

			if errors.As(err, &rerr) {
				return lvl
			}

			ctx.HandleErrorSilently(err)
		}

		return Show
	}
}
