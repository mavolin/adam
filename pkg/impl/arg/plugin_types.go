package arg

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// Command
// =====================================================================================

// Command is the type used for commands.
//
// Go type: *plugin.RegisteredCommand
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
		return nil, newArgParsingErr(commandNotFoundProvidersUnavailable, ctx, map[string]interface{}{
			"invoke": plugin.IdentifierFromInvoke(ctx.Raw).AsInvoke(), // remove whitespaces
		})
	}

	return nil, newArgParsingErr(commandNotFound, ctx, map[string]interface{}{
		"invoke": plugin.IdentifierFromInvoke(ctx.Raw).AsInvoke(), // remove whitespaces
	})
}

func (c command) Default() interface{} {
	return (*plugin.RegisteredCommand)(nil)
}

// =============================================================================
// Module
// =====================================================================================

// Module is the type used for modules.
//
// Go type: *plugin.RegisteredModule
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
	cmd := ctx.FindModule(ctx.Raw)
	if cmd != nil {
		return cmd, nil
	}

	if len(ctx.UnavailablePluginProviders()) > 0 {
		return nil, newArgParsingErr(moduleNotFoundProvidersUnavailable, ctx, map[string]interface{}{
			"invoke": plugin.IdentifierFromInvoke(ctx.Raw).AsInvoke(), // remove whitespaces
		})
	}

	return nil, newArgParsingErr(moduleNotFound, ctx, map[string]interface{}{
		"invoke": plugin.IdentifierFromInvoke(ctx.Raw).AsInvoke(), // remove whitespaces
	})
}

func (m module) Default() interface{} {
	return (*plugin.RegisteredModule)(nil)
}
