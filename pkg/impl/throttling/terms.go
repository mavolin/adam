package throttling

import "github.com/mavolin/adam/pkg/localization"

var (
	channelErrorSecond = localization.Config{
		Term: "throttling.channel.second",
		Fallback: localization.Fallback{
			One:   "This command can be used again in this channel in one second.",
			Other: "This command can be used again in this channel in {{.seconds}} seconds.",
		},
	}
	channelErrorMinute = localization.Config{
		Term: "throttling.channel.minute",
		Fallback: localization.Fallback{
			One:   "This command can be used again in this channel in one minute.",
			Other: "This command can be used again in this channel in {{.minutes}} minutes.",
		},
	}

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
