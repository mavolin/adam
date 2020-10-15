// Package module provides implementations for the module abstractions found
// in package plugin.
package module

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Module is an implementation of plugin.Module with support for middlewares.
type Module struct {
	Meta plugin.ModuleMeta

	bot.MiddlewareManager

	commands []plugin.Command
	modules  []plugin.Module
}

// AddCommands adds the passed command to the module.
func (m *Module) AddCommand(cmd plugin.Command) {
	m.commands = append(m.commands, cmd)
}

// AddModule adds the passed module to the module.
func (m *Module) AddModule(mod plugin.Module) {
	m.modules = append(m.modules, mod)
}

func (m *Module) GetName() string                              { return m.Meta.GetName() }
func (m *Module) GetShortDescription(l *i18n.Localizer) string { return m.Meta.GetShortDescription(l) }
func (m *Module) GetLongDescription(l *i18n.Localizer) string  { return m.Meta.GetLongDescription(l) }
func (m *Module) IsHidden() bool                               { return m.Meta.IsHidden() }
func (m *Module) GetDefaultChannelTypes() plugin.ChannelTypes  { return m.Meta.GetDefaultChannelTypes() }

func (m *Module) GetDefaultBotPermissions() *discord.Permissions {
	return m.Meta.GetDefaultBotPermissions()
}

func (m *Module) GetDefaultRestrictionFunc() plugin.RestrictionFunc {
	return m.Meta.GetDefaultRestrictionFunc()
}

func (m *Module) GetDefaultThrottler() plugin.Throttler { return m.Meta.GetDefaultThrottler() }
func (m *Module) Commands() []plugin.Command            { return m.commands }
func (m *Module) Modules() []plugin.Module              { return m.modules }
