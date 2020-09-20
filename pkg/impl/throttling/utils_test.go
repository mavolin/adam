package throttling

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/localization"
)

func Test_genError(t *testing.T) {
	testCases := []struct {
		name          string
		duration      time.Duration
		expectSeconds int
		expectMinutes int
	}{
		{
			name:          "round to 0",
			duration:      499 * time.Millisecond,
			expectSeconds: 0,
			expectMinutes: 0,
		},
		{
			name:          "second",
			duration:      90 * time.Second,
			expectSeconds: 90,
			expectMinutes: 0,
		},
		{
			name:          "minute - round up",
			duration:      91 * time.Second,
			expectSeconds: 0,
			expectMinutes: 2,
		},
		{
			name:          "minute - round down",
			duration:      120 * time.Second,
			expectSeconds: 0,
			expectMinutes: 2,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			secondConfig := localization.NewTermConfig("second")
			minuteConfig := localization.NewTermConfig("minute")

			acutal := genError(c.duration, secondConfig, minuteConfig)

			switch {
			case c.expectSeconds > 0:
				assert.Equal(t, errors.NewThrottlingErrorl(secondConfig.
					WithPlaceholders(secondPlaceholders{
						Seconds: c.expectSeconds,
					})), acutal)
			case c.expectMinutes > 0:
				assert.Equal(t, errors.NewThrottlingErrorl(minuteConfig.
					WithPlaceholders(minutePlaceholders{
						Minutes: c.expectMinutes,
					})), acutal)
			default:
				assert.Nil(t, acutal)
			}
		})
	}
}
