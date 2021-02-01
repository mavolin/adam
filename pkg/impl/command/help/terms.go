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
// General
// =====================================================================================

// ================================ Common ================================

var commandsFieldName = i18n.NewFallbackConfig("plugin.help.common.commands", "Commands")

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

	allModuleFieldName = i18n.NewFallbackConfig("plugin.help.all.embed.field.module.name", "`{{.module}}`-Module")
)

type allModuleFieldNamePlaceholders struct {
	Module string
}
