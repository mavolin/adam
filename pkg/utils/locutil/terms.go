package locutil

import "github.com/mavolin/adam/pkg/localization"

var (
	// defaultSeparatorConfig is the default separator.
	defaultSeparatorConfig = localization.NewFallbackConfig("common.lists.default_separator", ", ")
	// lastSeparatorConfig is the last separator of a list.
	lastSeparatorConfig = localization.NewFallbackConfig("common.lists.last_and", " and ")
)
