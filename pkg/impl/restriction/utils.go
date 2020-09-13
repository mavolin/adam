package restriction

import (
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/pluginutil"
)

// assertChannelTypes asserts that the command with the passed context
// is used in the passed channel types.
//
// assertChannelTypes will also silently report errors in some cases.
func assertChannelTypes(ctx *plugin.Context, assertTypes plugin.ChannelTypes, noRemainingError error) error {
	if assertTypes&plugin.AllChannels == plugin.AllChannels {
		return nil
	}

	if ctx.GuildID == 0 { // we are in a DM
		// we assert a DM
		if assertTypes&plugin.DirectMessages == plugin.DirectMessages {
			return nil
		}
		// no DM falls through
	} else { // we are in a guild
		// we assert all guild channels
		if assertTypes&plugin.GuildChannels == plugin.GuildChannels {
			return nil

			// we assert something other than all guild channels
		} else if !(assertTypes&plugin.GuildChannels == 0) {
			c, err := ctx.Channel()
			if err != nil {
				return err
			}

			if assertTypes.Has(c.Type) {
				return nil
			}
		}
		// not all guild types falls through
	}

	channelTypes, err := pluginutil.ChannelTypes(ctx.CommandIdentifier, ctx.Provider)
	if err != nil {
		return err
	}

	allowed := channelTypes & assertTypes
	if allowed == 0 { // no channel types remaining
		// there is no need to prevent execution, as another restriction
		// may permit it, still we should capture this
		ctx.HandleErrorSilent(noRemainingError)

		return errors.DefaultFatalRestrictionError
	}

	return newInvalidChannelTypeError(allowed, ctx.Localizer, true)
}
