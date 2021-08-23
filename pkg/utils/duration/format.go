package duration

import (
	"fmt"
	"strings"
	"time"
)

// Format formats the passed duration rounded to milliseconds.
// Individual Values are separated.
//
// For example: '1h 10s'.
func Format(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	b := new(strings.Builder)
	b.Grow(36)

	if d < 0 {
		b.WriteRune('-')

		d++ // prevent an overflow, since maxDuration = -minDuration - 1
		d = -d
	}

	d = d.Round(Millisecond)
	if d == 0 {
		return "0s"
	}

	writeUnit(Year, "y", b, &d)
	writeUnit(Month, "m", b, &d)
	writeUnit(Week, "w", b, &d)
	writeUnit(Day, "d", b, &d)
	writeUnit(Hour, "h", b, &d)
	writeUnit(Minute, "min", b, &d)
	writeUnit(Second, "s", b, &d)
	writeUnit(Millisecond, "ms", b, &d)

	return b.String()
}

func writeUnit(mul time.Duration, unit string, b *strings.Builder, d *time.Duration) {
	if val := *d / mul; val > 0 {
		if b.Len() > 1 { // don't add a space in front of a -
			b.WriteRune(' ')
		}

		b.WriteString(fmt.Sprintf("%d%s", val, unit))

		*d %= mul
	}
}
