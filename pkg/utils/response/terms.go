package response

import "github.com/mavolin/adam/pkg/localization"

var (
	defaultCancelKeyword = localization.NewFallbackConfig("response_waiter.cancel.default", "cancel")
	timeoutInfo          = localization.NewFallbackConfig("response_waiter.infos.timeout",
		"{{.response_user_mention}} I haven't heard back from you, so I won't wait for an answer any longer. "+
			"Try again if you ran out of time.")

	timeExtensionTitle = localization.NewFallbackConfig("response_waiter.time_extension.title",
		"Are you still there?")
	timeExtensionDescription = localization.NewFallbackConfig("response_waiter.time_extension.description",
		"{{.response_user_mention}} If you are still answering, click the {{.time_extension_reaction}} reaction on this "+
			"message.")
)

type (
	timeoutInfoPlaceholders struct {
		ResponseUserMention string
	}

	timeExtensionDescriptionPlaceholders struct {
		ResponseUserMention   string
		TimeExtensionReaction string
	}
)
