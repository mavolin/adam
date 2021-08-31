package say

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/command"
)

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig("plugin.say.short_description", "Repeats what you say.")

	examples = command.LocalizedExampleArgs{
		{
			Args: []*i18n.Config{
				i18n.NewFallbackConfig("plugin.say.example.hello.arg.0", "Hello"),
			},
		},
	}
)

// =============================================================================
// Arguments
// =====================================================================================

var (
	argTextName        = i18n.NewFallbackConfig("plugin.say.arg.text.name", "Text")
	argTextDescription = i18n.NewFallbackConfig("plugin.say.arg.text.description", "The text you want to say.")
)
