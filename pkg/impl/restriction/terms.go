package restriction

import "github.com/mavolin/adam/pkg/localization"

var (
	anyMessage = localization.NewFallbackConfig("restriction.compare.any",
		"You need to fulfill at least one of these requirements:")
	allMessage = localization.NewFallbackConfig("restriction.compare.all",
		"You need to fulfill all of these requirements:")
)
