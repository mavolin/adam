package help

import "github.com/mavolin/adam/pkg/i18n"

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig(
		"plugin.help.short_description",
		"Lists all commands or shows you the usage of a command.")

	longDescription = i18n.NewFallbackConfig(
		"plugin.help.long_description",
		"Lists all commands or the commands in a certain module. "+
			"Furthermore, it can also show detailed information about a command.")

	examples = []*i18n.Config{
		i18n.NewFallbackConfig("plugin.help.examples.list_all", "help"),
		i18n.NewFallbackConfig("plugin.help.examples.command", "help some_command"),
		i18n.NewFallbackConfig("plugin.help.examples.module", "help some_module"),
	}
)

// =============================================================================
// Args
// =====================================================================================

var (
	argsPluginName        = i18n.NewFallbackConfig("plugin.help.args.plugin.name", "Command or Module")
	argsPluginDescription = i18n.NewFallbackConfig(
		"plugin.help.args.plugin.description",
		"The name of the command or module you need help with.")
)

// =============================================================================
// Text
// =====================================================================================

// ================================ Common ================================

var (
	commandsFieldName = i18n.NewFallbackConfig("plugin.help.common.commands", "Commands")

	moduleTitle = i18n.NewFallbackConfig("plugin.help.common.title", "`{{.module}}`-Module")

	// ===== below two are copy of variables in the arg package =====

	pluginNotFoundError = i18n.NewFallbackConfig(
		"arg.types.plugin.error.not_found",
		"I don't know any commands or modules with the name `{{.invoke}}`. Make sure you spelled it right.")

	pluginNotFoundErrorProvidersUnavailable = i18n.NewFallbackConfig(
		"arg.types.plugin.error.not_found.providers_unavailable",
		"I couldn't find any commands or modules with the name `{{.invoke}}`, "+
			"however, I'm having trouble accessing some of my commands, so this may be why. "+
			"Try again later or check your spelling.")
)

type moduleTitlePlaceholders struct {
	Module string
}

// ================================ All ================================

var (
	allTitle = i18n.NewFallbackConfig("plugin.help.all.embed.title", "Help")

	allDescriptionDM = i18n.NewFallbackConfig(
		"plugin.help.all.embed.description.dm",
		"Below is a list of all commands accessible through direct messages.")
	allDescriptionGuild = i18n.NewFallbackConfig(
		"plugin.help.all.embed.description.guild",
		"Below is a list of all commands on this server.")

	allPrefixesFieldName = i18n.NewFallbackConfig("plugin.help.all.embed.field.prefix.name", "Prefixes")
)
