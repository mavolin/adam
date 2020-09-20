package throttling

import (
	"time"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/localization"
)

// genError generates a errors.ThrottlingError using one of the two passed
// localization.Configs, based on the value of the passed time.Duration.
// Any duration less or equal to 90 seconds will displayed using seconds.
// Otherwise the minuteConfig will be used.
func genError(
	d time.Duration, secondConfig localization.Config, minuteConfig localization.Config,
) *errors.ThrottlingError {
	d = d.Round(time.Second)

	if d <= 0 {
		return nil
	} else if d <= 90*time.Second { // display up to 90 seconds in seconds
		return errors.NewThrottlingErrorl(secondConfig.
			WithPlaceholders(secondPlaceholders{
				Seconds: int(d / time.Second),
			}))
	}

	// otherwise switch to minutes

	d = d.Round(time.Minute)

	return errors.NewThrottlingErrorl(minuteConfig.
		WithPlaceholders(minutePlaceholders{
			Minutes: int(d / time.Minute),
		}))
}
