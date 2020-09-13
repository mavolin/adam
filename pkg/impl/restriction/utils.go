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
	var has bool

	if assertTypes == plugin.GuildChannels {
		has = ctx.GuildID != 0
	} else if assertTypes == plugin.DirectMessages {
		has = ctx.GuildID == 0
	} else if assertTypes == plugin.AllChannels {
		has = true
	} else if assertTypes == plugin.GuildTextChannels || assertTypes == plugin.GuildNewsChannels {
		if ctx.GuildID == 0 {
			has = false
		} else {
			c, err := ctx.Channel()
			if err != nil {
				return err
			}

			has = assertTypes.Has(c.Type)
		}
	} else {
		c, err := ctx.Channel()
		if err != nil {
			return err
		}

		has = assertTypes.Has(c.Type)
	}

	if !has {
		channelTypes, err := pluginutil.ChannelTypes(ctx.CommandIdentifier, ctx.Provider)
		if err != nil {
			return err
		} else if channelTypes == 0 {
			return errors.DefaultFatalRestrictionError
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

	return nil
}
