package throttler

import "github.com/mavolin/adam/pkg/i18n"

// ================================ channel ================================

var (
	channelThrottledErrorSecond = &i18n.Config{
		Term: "throttler.channel.errors.throttled.second",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this channel in one second.",
			Other: "This command can be used again in this channel in {{.seconds}} seconds.",
		},
	}
	channelThrottledErrorMinute = &i18n.Config{
		Term: "throttler.channel.errors.throttled.minute",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this channel in one minute.",
			Other: "This command can be used again in this channel in {{.minutes}} minutes.",
		},
	}
)

// ================================ guild ================================

var (
	guildThrottledErrorSecond = &i18n.Config{
		Term: "throttler.guild.errors.throttled.second",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this server in one second.",
			Other: "This command can be used again in this server in {{.seconds}} seconds.",
		},
	}
	guildThrottledErrorMinute = &i18n.Config{
		Term: "throttler.guild.errors.throttled.minute",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this server in one minute.",
			Other: "This command can be used again in this server in {{.minutes}} minutes.",
		},
	}
)

// ================================ member ================================

var (
	memberThrottledErrorSecond = &i18n.Config{
		Term: "throttler.member.errors.throttled.second",
		Fallback: i18n.Fallback{
			One:   "You can use this command again in one second.",
			Other: "You can use this command again in {{.seconds}} seconds.",
		},
	}
	memberThrottledErrorMinute = &i18n.Config{
		Term: "throttler.member.errors.throttled.minute",
		Fallback: i18n.Fallback{
			One:   "You can use this command again in one minute.",
			Other: "You can use this command again in {{.minutes}} minutes.",
		},
	}
)

// ================================ user ================================

var (
	userThrottledErrorSecond = &i18n.Config{
		Term: "throttler.user.errors.throttled.second",
		Fallback: i18n.Fallback{
			One:   "You can use this command again in one second.",
			Other: "You can use this command again in {{.seconds}} seconds.",
		},
	}
	userThrottledErrorMinute = &i18n.Config{
		Term: "throttler.user.errors.throttled.minute",
		Fallback: i18n.Fallback{
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
