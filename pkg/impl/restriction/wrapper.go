package restriction

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// Wrap wraps the passed plugin.RestrictionFunc to provide proper support for
// ALL and ANY.
func Wrap(f plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		restriction := f(s, ctx)

		switch restriction := restriction.(type) {
		case *allError:
			missing, err := restriction.format(0, ctx.Localizer)
			if err != nil {
				return err
			}

			header, _ := ctx.Localize(allMessage)

			return errors.NewRestrictionError(header + "\n\n" + missing)
		case *anyError:
			missing, err := restriction.format(0, ctx.Localizer)
			if err != nil {
				return err
			}

			header, _ := ctx.Localize(anyMessage)

			return errors.NewRestrictionError(header + "\n\n" + missing)
		default:
			return restriction
		}
	}
}
