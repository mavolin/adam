package mod

import "github.com/mavolin/adam/pkg/i18n"

var (
	shortDescription = i18n.NewFallbackConfig("plugin.mod.short_description", "Moderate a server.")

	longDescription = i18n.NewFallbackConfig("plugin.mod.long_description",
		"The moderation module provides utilities for moderating your server easily.")
)
