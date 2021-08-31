package ban

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/command"
)

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig("plugin.mod.kick.short_description", "Bans someone.")

	examples = command.LocalizedExampleArgs{
		{
			Args: []*i18n.Config{
				i18n.NewFallbackConfig("plugin.mod.kick.example.plain.arg.0", "@Wumpus"),
			},
		},
		{
			Args: []*i18n.Config{
				i18n.NewFallbackConfig("plugin.mod.kick.example.reason.arg.0", "@Wumpus"),
				i18n.NewFallbackConfig("plugin.mod.kick.example.reason.arg.1", "using offensive language"),
			},
		},
	}
)

// =============================================================================
// Arguments
// =====================================================================================

var (
	argMemberName        = i18n.NewFallbackConfig("plugin.mod.kick.arg.member.name", "Member")
	argMemberDescription = i18n.NewFallbackConfig(
		"plugin.mod.kick.args.member.description",
		"The member you want to ban.")

	argReasonName        = i18n.NewFallbackConfig("plugin.mod.kick.arg.reason.name", "Reason")
	argReasonDescription = i18n.NewFallbackConfig(
		"plugin.mod.kick.args.reason.description",
		"The reason for the ban.")
)

// =============================================================================
// Flags
// =====================================================================================

var flagDaysDescription = i18n.NewFallbackConfig(
	"plugin.mod.kick.flag.days.description",
	"The amount of days to delete messages for. You can delete 7 days at most.")

// =============================================================================
// Response
// =====================================================================================

var success = i18n.NewFallbackConfig(
	"plugin.mod.ban.response.success",
	"👮 The banhammer has been slayed, and {{.username}} is no more!")

type successPlaceholders struct {
	Username string
}

// =============================================================================
// Errors
// =====================================================================================

var selfBanError = i18n.NewFallbackConfig("plugin.mod.ban.error.self_ban", "Good try, but you can ban yourself.")
