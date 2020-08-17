package locutil

import "github.com/mavolin/adam/pkg/localization"

var (
	// defaultSeparatorConfig is the default separator.
	defaultSeparatorConfig = localization.NewFallbackConfig("lang.lists.default_separator", ", ")
	// lastSeparatorConfig is the last separator of a list.
	lastSepartatorConfig = localization.NewFallbackConfig("lang.lists.last_separator", " and ")
)
