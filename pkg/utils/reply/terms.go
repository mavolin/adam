package reply

import "github.com/mavolin/adam/pkg/i18n"

var (
	defaultCancelKeyword = i18n.NewFallbackConfig("response_waiter.cancel.default", "cancel")

	timeoutInfo = i18n.NewFallbackConfig("response_waiter.infos.timeout",
		"{{.mention}} I haven't heard back from you, so I won't wait for an answer any longer. "+
			"Try again if you ran out of time.")
)

type timeoutInfoPlaceholders struct {
	ResponseUserMention string
}
