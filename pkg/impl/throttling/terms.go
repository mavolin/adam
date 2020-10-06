package throttling

import "github.com/mavolin/adam/pkg/i18n"

var (
	channelErrorSecond = i18n.Config{
		Term: "throttling.channel.second",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this channel in one second.",
			Other: "This command can be used again in this channel in {{.seconds}} seconds.",
		},
	}
	channelErrorMinute = i18n.Config{
		Term: "throttling.channel.minute",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this channel in one minute.",
			Other: "This command can be used again in this channel in {{.minutes}} minutes.",
		},
	}

	guildErrorSecond = i18n.Config{
		Term: "throttling.guild.second",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this server in one second.",
			Other: "This command can be used again in this server in {{.seconds}} seconds.",
		},
	}
	guildErrorMinute = i18n.Config{
		Term: "throttling.guild.minute",
		Fallback: i18n.Fallback{
			One:   "This command can be used again in this server in one minute.",
			Other: "This command can be used again in this server in {{.minutes}} minutes.",
		},
	}

	memberErrorSecond = i18n.Config{
		Term: "throttling.member.second",
		Fallback: i18n.Fallback{
			One:   "You can use this command again in this guild in one second.",
			Other: "You can use this command again in this guild in {{.seconds}} seconds.",
		},
	}
	memberErrorMinute = i18n.Config{
		Term: "throttling.member.minute",
		Fallback: i18n.Fallback{
			One:   "You can use this command again in this guild in one minute.",
			Other: "You can use this command again in this guild in {{.minutes}} minutes.",
		},
	}

	userErrorSecond = i18n.Config{
		Term: "throttling.user.second",
		Fallback: i18n.Fallback{
			One:   "You can use this command again in one second.",
			Other: "You can use this command again in {{.seconds}} seconds.",
		},
	}
	userErrorMinute = i18n.Config{
		Term: "throttling.user.minute",
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
