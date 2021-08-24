package arg

import (
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// Command
// =====================================================================================

// Command is the type used for commands.
//
// Go type: plugin.ResolvedCommand
var Command plugin.ArgType = new(commandType)

type commandType struct{}

func (c commandType) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(commandName) // we have a fallback
	return name
}

func (c commandType) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(commandDescription) // we have a fallback
	return desc
}

func (c commandType) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	cmd := ctx.FindCommand(ctx.Raw)
	if cmd != nil {
		return cmd, nil
	}

	if len(ctx.UnavailablePluginSources()) > 0 {
		return nil, newArgumentError(commandNotFoundErrorProvidersUnavailable, ctx, nil)
	}

	return nil, newArgumentError(commandNotFoundError, ctx, nil)
}

func (c commandType) GetDefault() interface{} {
	return (plugin.ResolvedCommand)(nil)
}

// =============================================================================
// Module
// =====================================================================================

// Module is the type used for modules.
//
// Go type: plugin.ResolvedModule
var Module plugin.ArgType = new(moduleType)

type moduleType struct{}

func (m moduleType) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(moduleName) // we have a fallback
	return name
}

func (m moduleType) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(moduleDescription) // we have a fallback
	return desc
}

func (m moduleType) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	mod := ctx.FindModule(ctx.Raw)
	if mod != nil {
		return mod, nil
	}

	if len(ctx.UnavailablePluginSources()) > 0 {
		return nil, newArgumentError(moduleNotFoundErrorProvidersUnavailable, ctx, nil)
	}

	return nil, newArgumentError(moduleNotFoundError, ctx, nil)
}

func (m moduleType) GetDefault() interface{} {
	return (plugin.ResolvedModule)(nil)
}

// =============================================================================
// Plugin
// =====================================================================================

// Plugin is the type used for plugins, i.e. commands and Modules.
// The generated data is guaranteed to be of one of the two go types, unless
// falling back to default, i.e. if using Plugin as type for an optional arg or
// a flag without a custom default value.
//
// Go types: plugin.ResolvedCommand or plugin.ResolvedModule
var Plugin plugin.ArgType = new(pluginType)

type pluginType struct{}

func (p pluginType) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(pluginName) // we have a fallback
	return name
}

func (p pluginType) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(pluginDescription) // we have a fallback
	return desc
}

func (p pluginType) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	if cmd := ctx.FindCommand(ctx.Raw); cmd != nil {
		return cmd, nil
	} else if mod := ctx.FindModule(ctx.Raw); mod != nil {
		return mod, nil
	}

	if len(ctx.UnavailablePluginSources()) > 0 {
		return nil, newArgumentError(pluginNotFoundErrorProvidersUnavailable, ctx, nil)
	}

	return nil, newArgumentError(pluginNotFoundError, ctx, nil)
}

func (p pluginType) GetDefault() interface{} {
	// return interface{} nil, as described in type's doc
	return nil
}
