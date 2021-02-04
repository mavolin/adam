package say

import "github.com/mavolin/adam/pkg/i18n"

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig("plugin.say.short_description", "Repeats what you say.")

	examples = []*i18n.Config{
		i18n.NewFallbackConfig("plugin.say.example.hello", "Hello"),
	}
)

// =============================================================================
// Arguments
// =====================================================================================

var (
	argTextName        = i18n.NewFallbackConfig("plugin.say.arg.text.name", "Text")
	argTextDescription = i18n.NewFallbackConfig("plugin.say.arg.text.description", "The text you want to say.")
)
