package kick

import "github.com/mavolin/adam/pkg/i18n"

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig("plugin.mod.kick.short_description", "Kicks a user.")

	examples = []*i18n.Config{
		i18n.NewFallbackConfig("plugin.mod.kick.example.plain", "@Clyde"),
		i18n.NewFallbackConfig("plugin.mod.kick.example.reason", "@Clyde, self-botting"),
	}
)

// =============================================================================
// Arguments
// =====================================================================================

var (
	argMemberName        = i18n.NewFallbackConfig("plugin.mod.kick.args.member.name", "Member")
	argMemberDescription = i18n.NewFallbackConfig(
		"plugin.mod.kick.args.member.description",
		"The member you want to kick.")
)

// =============================================================================
// Response
// =====================================================================================

var success = i18n.NewFallbackConfig("plugin.mod.kick.response.success", "👮 {{.username}} has been kicked!")

type successPlaceholders struct {
	Username string
}

// =============================================================================
// Errors
// =====================================================================================

var selfKickError = i18n.NewFallbackConfig("plugin.mod.kick.error.self_kick", "You can't kick yourself!")
