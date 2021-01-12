// Package module provides implementations for the module abstractions found
// in package plugin.
package module

import (
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
)

// Module is an implementation of plugin.Module with support for middlewares.
type Module struct {
	plugin.ModuleMeta
	*bot.MiddlewareManager

	commands []plugin.Command
	modules  []plugin.Module
}

var _ plugin.Module = new(Module)

// New creates a new *Module using the passed plugin.ModuleMeta.
func New(meta plugin.ModuleMeta) *Module {
	return &Module{
		ModuleMeta:        meta,
		MiddlewareManager: new(bot.MiddlewareManager),
	}
}

// AddCommands adds the passed command to the module.
func (m *Module) AddCommand(cmd plugin.Command) {
	m.commands = append(m.commands, cmd)
}

// AddModule adds the passed module to the module.
func (m *Module) AddModule(mod plugin.Module) {
	m.modules = append(m.modules, mod)
}

func (m *Module) Commands() []plugin.Command {
	return m.commands
}

func (m *Module) Modules() []plugin.Module {
	return m.modules
}
