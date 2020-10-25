// Package durationutil provides utilities for interacting with durations.
//
// In contrast to the functions provided by the time package, durationutil also
// supports the units Day, Week, Month and Year, but does not provide support
// for units smaller than a millisecond.
package durationutil

import "time"

const (
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = 24 * Hour
	Week        = 7 * Day
	Month       = 30 * Day
	Year        = 365 * Day
)

// units are units used during parsing.
var units = map[string]int64{
	"ms":  int64(Millisecond),
	"s":   int64(Second),
	"min": int64(Minute),
	"h":   int64(Hour),
	"d":   int64(Day),
	"w":   int64(Week),
	"m":   int64(Month),
	"y":   int64(Year),
	"yr":  int64(Year),
}
