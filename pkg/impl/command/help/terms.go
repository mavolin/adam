package help

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/command"
)

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig(
		"plugin.help.short_description",
		"Lists all commands or shows you the usage of a command.")

	longDescription = i18n.NewFallbackConfig(
		"plugin.help.long_description",
		"Lists all commands and shows you how to use them.\n"+
			"If you don't provide any parameters, "+
			"the help command will show you a list of all commands available to you. "+
			"Optionally, you can use the name of a module to list all commands in that module, "+
			"or the name of a command to display detailed usage information.")

	exampleArgs = command.LocalizedExampleArgs{
		{},
		{
			Args: []*i18n.Config{
				i18n.NewFallbackConfig("plugin.help.example.command.arg.0", "some_command"),
			},
		},
		{
			Args: []*i18n.Config{
				i18n.NewFallbackConfig("plugin.help.example.module.arg.0", "some_module"),
			},
		},
	}
)

// =============================================================================
// Arguments
// =====================================================================================

var (
	argPluginName        = i18n.NewFallbackConfig("plugin.help.arg.plugin.name", "Command or Module")
	argPluginDescription = i18n.NewFallbackConfig(
		"plugin.help.arg.plugin.description",
		"The name of the command or module you need help with.")
)

// =============================================================================
// Response
// =====================================================================================

// ================================ Common ================================

var (
	commandsFieldName = i18n.NewFallbackConfig("plugin.help.common.commands", "Commands")

	moduleTitle = i18n.NewFallbackConfig("plugin.help.common.module_title", "`{{.module}}` Module")

	// pluginNotFoundError is a copy of the variable with the same name found
	// in the arg package.
	pluginNotFoundError = i18n.NewFallbackConfig(
		"arg.types.plugin.error.not_found",
		"I don't know any commands or modules with the name `{{.invoke}}`. Make sure you spelled it right.")
)

type (
	moduleTitlePlaceholders struct {
		Module string
	}

	pluginNotFoundErrorPlaceholder struct {
		Invoke string
	}
)

// ================================ All ================================

var (
	allTitle = i18n.NewFallbackConfig("plugin.help.all.embed.title", "Help")

	allDescriptionDM = i18n.NewFallbackConfig(
		"plugin.help.all.embed.description.dm",
		"Below is a list of all commands accessible through direct messages.")
	allDescriptionGuild = i18n.NewFallbackConfig(
		"plugin.help.all.embed.description.guild",
		"Below is a list of all commands on this server.")

	allPrefixesFieldName = i18n.NewFallbackConfig("plugin.help.all.embed.field.prefixes.name", "Prefixes")
)

// ================================ Command ================================

var (
	commandTitle = i18n.NewFallbackConfig("plugin.help.command.embed.title", "`{{.command}}` Command")

	aliasesFieldName = i18n.NewFallbackConfig("plugin.help.command.embed.fields.aliases.name", "Aliases")

	usageFieldNameSingle = i18n.NewFallbackConfig("plugin.help.command.embed.fields.usage.name", "Usage")

	argumentsFieldName = i18n.NewFallbackConfig("plugin.help.command.embed.fields.arguments.name", "Arguments")
	flagsFieldName     = i18n.NewFallbackConfig("plugin.help.command.embed.fields.flags.name", "Flags")
	examplesFieldName  = i18n.NewFallbackConfig("plugin.help.command.embed.fields.examples.name", "Examples")
)

type commandTitlePlaceholders struct {
	Command string
}
