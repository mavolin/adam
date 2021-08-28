// Package mock provides mocks for key types to aid in unit testing.
package mock

import (
	"testing"

	"github.com/mavolin/adam/internal/mock/i18n"
	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/internal/mock/resolve"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// i18n
// =====================================================================================

type Localizer = i18n.Localizer

// NewLocalizer creates a new Localizer.
// If a term is not found, Localizer will panic.
//nolint:thelper
func NewLocalizer(t *testing.T) *Localizer {
	return i18n.NewLocalizer(t)
}

// NewLocalizerWithDefault creates a new Localizer using the passed default.
// If a term is not found, Localizer will return the default value.
func NewLocalizerWithDefault(t *testing.T, def string) *Localizer { //nolint:thelper
	return i18n.NewLocalizerWithDefault(t, def)
}

// =============================================================================
// plugin
// =====================================================================================

type DiscordDataProvider = mockplugin.DiscordDataProvider

type ErrorHandler = mockplugin.ErrorHandler

//nolint:thelper
func NewErrorHandler(t *testing.T) *ErrorHandler {
	return mockplugin.NewErrorHandler(t)
}

type (
	Command = mockplugin.Command
	Module  = mockplugin.Module
)

type Throttler = mockplugin.Throttler

// NewThrottler creates a new mocked Throttler with the given return value
// for check.
func NewThrottler(checkReturn error) *Throttler {
	return mockplugin.NewThrottler(checkReturn)
}

func RestrictionFunc(ret error) plugin.RestrictionFunc {
	return mockplugin.RestrictionFunc(ret)
}

type RestrictionErrorWrapper = mockplugin.RestrictionErrorWrapper

// =============================================================================
// resolve
// =====================================================================================

type PluginProvider = resolve.Provider

// ResolveCommand creates a plugin.ResolvedCommand from the passed
// plugin.Command using the passed source name.
func ResolveCommand(sourceName string, cmd plugin.Command) plugin.ResolvedCommand {
	return resolve.Command(sourceName, cmd)
}

// ResolveModule creates a plugin.ResolvedModule from the passed
// plugin.Module using the passed source name.
func ResolveModule(sourceName string, mod plugin.Module) plugin.ResolvedModule {
	return resolve.Module(sourceName, mod)
}
