package arg

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// Command
// =====================================================================================

// Command is the type used for commands.
//
// Go type: *plugin.ResolvedCommand
var Command Type = new(command)

type command struct{}

func (c command) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(commandName) // we have a fallback
	return name
}

func (c command) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(commandDescription) // we have a fallback
	return desc
}

func (c command) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	cmd := ctx.FindCommand(ctx.Raw)
	if cmd != nil {
		return cmd, nil
	}

	if len(ctx.UnavailablePluginProviders()) > 0 {
		return nil, newArgumentError(commandNotFoundErrorProvidersUnavailable, ctx, nil)
	}

	return nil, newArgumentError(commandNotFoundError, ctx, nil)
}

func (c command) Default() interface{} {
	return (*plugin.ResolvedCommand)(nil)
}

// =============================================================================
// Module
// =====================================================================================

// Module is the type used for modules.
//
// Go type: *plugin.ResolvedModule
var Module Type = new(module)

type module struct{}

func (m module) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(moduleName) // we have a fallback
	return name
}

func (m module) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(moduleDescription) // we have a fallback
	return desc
}

func (m module) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	mod := ctx.FindModule(ctx.Raw)
	if mod != nil {
		return mod, nil
	}

	if len(ctx.UnavailablePluginProviders()) > 0 {
		return nil, newArgumentError(moduleNotFoundErrorProvidersUnavailable, ctx, nil)
	}

	return nil, newArgumentError(moduleNotFoundError, ctx, nil)
}

func (m module) Default() interface{} {
	return (*plugin.ResolvedModule)(nil)
}

// =============================================================================
// Plugin
// =====================================================================================

// Plugin is the type used for plugins, i.e. Commands and Modules.
// The generated data is guaranteed to be of one of the two go types.
// Fallback for default will be interface{} nil.
//
// Go types: *plugin.ResolvedCommand or *plugin.ResolvedModule
var Plugin Type = new(pluginType)

type pluginType struct{}

func (p pluginType) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(pluginName) // we have a fallback
	return name
}

func (p pluginType) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(pluginDescription) // we have a fallback
	return desc
}

func (p pluginType) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	if cmd := ctx.FindCommand(ctx.Raw); cmd != nil {
		return cmd, nil
	} else if mod := ctx.FindModule(ctx.Raw); mod != nil {
		return mod, nil
	}

	if len(ctx.UnavailablePluginProviders()) > 0 {
		return nil, newArgumentError(pluginNotFoundErrorProvidersUnavailable, ctx, nil)
	}

	return nil, newArgumentError(pluginNotFoundError, ctx, nil)
}

func (p pluginType) Default() interface{} {
	// return interface{} nil, as described in type doc
	return nil
}
