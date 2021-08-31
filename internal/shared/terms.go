package shared

import "github.com/mavolin/adam/pkg/i18n"

var (
	ErrorTitle = i18n.NewFallbackConfig("error.title", "Error")
	InfoTitle  = i18n.NewFallbackConfig("info.title", "Info")
)
