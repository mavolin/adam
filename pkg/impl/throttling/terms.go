package throttling

import "github.com/mavolin/adam/pkg/localization"

var (
	userErrorSecond = localization.Config{
		Term: "throttling.user.second",
		Fallback: localization.Fallback{
			One:   "You can use this command again in one second.",
			Other: "You can use this command again in {{.seconds}} seconds.",
		},
	}
	userErrorMinute = localization.Config{
		Term: "throttling.user.minute",
		Fallback: localization.Fallback{
			One:   "You can use this command again in one minute.",
			Other: "You can use this command again in {{.minutes}} minutes.",
		},
	}
)

type (
	secondPlaceholders struct {
		Seconds int
	}

	minutePlaceholders struct {
		Minutes int
	}
)
