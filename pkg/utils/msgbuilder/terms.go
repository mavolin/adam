package msgbuilder

import "github.com/mavolin/adam/pkg/i18n"

var timeoutInfo = i18n.NewFallbackConfig("message_component.await.info.timeout",
	"{{.mention}} I haven't heard back from you, so I won't wait for an answer any longer. "+
		"Try again if you ran out of time.")

type timeoutInfoPlaceholders struct {
	Mention string
}
