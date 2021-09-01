package shared

import "github.com/mavolin/adam/pkg/i18n"

var (
	errorTitle = i18n.NewFallbackConfig("error.title", "Error")
	infoTitle  = i18n.NewFallbackConfig("info.title", "Info")
)
