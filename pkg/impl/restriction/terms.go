package restriction

import "github.com/mavolin/adam/pkg/localization"

var (
	anyMessageHeader = localization.NewFallbackConfig(
		"restrictions.compare.any.header",
		"You need to fulfill at least one of these requirements to execute the command:")
	anyMessageInline = localization.NewFallbackConfig(
		"restrictions.compare.any.inline",
		"You need to fulfill at least one of these requirements:")

	allMessageHeader = localization.NewFallbackConfig(
		"restrictions.compare.all.header",
		"You need to fulfill all of these requirements to execute the command:")
	allMessageInline = localization.NewFallbackConfig(
		"restrictions.compare.all.inline",
		"You need to fulfill all of these requirements:")
)
