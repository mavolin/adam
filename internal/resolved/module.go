package resolved

import (
	"strings"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type Module struct {
	parent  plugin.ResolvedModule
	sources []plugin.SourceModule

	id     plugin.ID
	hidden bool

	commands []plugin.ResolvedCommand
	modules  []plugin.ResolvedModule
}

var _ plugin.ResolvedModule = new(Module)

func newModule(parent *Module, provider *PluginProvider, sourceName string, smod plugin.Module) *Module {
	parentInvoke := ""
	if parent != nil {
		parentInvoke = parent.id.AsInvoke() + " "
	}

	if _, ok := provider.usedNames[parentInvoke+smod.GetName()]; ok {
		return nil
	}

	rmod := &Module{
		sources:  []plugin.SourceModule{{SourceName: sourceName}},
		id:       plugin.ID("." + smod.GetName()),
		hidden:   true,
		commands: nil,
	}

	if len(smod.GetCommands()) > 0 {
		rmod.commands = make([]plugin.ResolvedCommand, len(smod.GetCommands()))
	}

	if len(smod.GetModules()) > 0 {
		rmod.modules = make([]plugin.ResolvedModule, len(smod.GetModules()))
	}

	if parent != nil {
		rmod.parent = parent

		for _, parentSource := range parent.Sources() {
			if parentSource.SourceName == sourceName {
				// append is safe since the underlying slice will never change
				rmod.sources[0].Modules = append(parentSource.Modules, smod) //nolint:gocritic
				break
			}
		}

		rmod.id = parent.id + rmod.id
	} else {
		rmod.sources[0].Modules = []plugin.Module{smod}
	}

	parentInvoke += rmod.Name() + " "

	for i, subScmd := range smod.GetCommands() {
		provider.usedNames[parentInvoke+subScmd.GetName()] = struct{}{}

		aliases := subScmd.GetAliases()
		for _, alias := range aliases {
			provider.usedNames[parentInvoke+alias] = struct{}{}
		}

		rmod.commands[i] = &Command{
			parent:        rmod,
			provider:      provider,
			sourceName:    sourceName,
			source:        subScmd,
			sourceParents: rmod.sources[0].Modules,
			id:            rmod.id + plugin.ID("."+subScmd.GetName()),
			aliases:       aliases,
		}

		if !subScmd.IsHidden() {
			rmod.hidden = false
		}
	}

	for i, subSmod := range smod.GetModules() {
		rmod.modules[i] = newModule(rmod, provider, sourceName, subSmod)
		if !rmod.modules[i].IsHidden() {
			rmod.hidden = false
		}
	}

	return rmod
}

func (mod *Module) update(provider *PluginProvider, sourceName string, smod plugin.Module) {
	if mod.Parent() != nil {
		for _, parentSource := range mod.Parent().Sources() {
			if parentSource.SourceName == sourceName {
				// append is safe since the underlying slice will never change
				mod.sources = append(mod.sources, plugin.SourceModule{
					SourceName: sourceName,
					Modules:    append(parentSource.Modules, smod),
				})
				break
			}
		}
	} else {
		mod.sources = append(mod.sources, plugin.SourceModule{
			SourceName: sourceName,
			Modules:    []plugin.Module{smod},
		})
	}

	parentInvoke := mod.id.AsInvoke() + " "

	for i, subScmd := range smod.GetCommands() {
		if _, ok := provider.usedNames[parentInvoke+subScmd.GetName()]; ok {
			continue
		}

		provider.usedNames[parentInvoke+subScmd.GetName()] = struct{}{}

		var aliases []string
		if len(subScmd.GetAliases()) > 0 {
			aliases = make([]string, len(subScmd.GetAliases()))
			copy(aliases, subScmd.GetAliases())
			for _, alias := range aliases {
				if _, ok := provider.usedNames[parentInvoke+alias]; ok {
					copy(aliases[i:], aliases[i+1:])
					aliases = aliases[:len(aliases)-1]
				}

				provider.usedNames[parentInvoke+alias] = struct{}{}
			}
		}

		subRcmd := &Command{
			parent:        mod,
			provider:      provider,
			sourceName:    sourceName,
			source:        subScmd,
			sourceParents: mod.sources[len(mod.sources)-1].Modules,
			id:            mod.id + plugin.ID("."+subScmd.GetName()),
			aliases:       aliases,
		}

		mod.commands = insertCommand(mod.commands, subRcmd, -1)
	}

	for _, subSmod := range smod.GetModules() {
		i := searchModule(mod.modules, smod.GetName())
		if i < len(mod.modules) && mod.modules[i].Name() == smod.GetName() {
			mod.modules[i].(*Module).update(provider, sourceName, subSmod)
			if !mod.modules[i].IsHidden() {
				mod.hidden = false
			}
		} else {
			subRmod := newModule(mod, provider, sourceName, subSmod)
			if subRmod != nil {
				mod.modules = insertModule(mod.modules, subRmod, i)
				if !subRmod.IsHidden() {
					mod.hidden = false
				}
			}
		}
	}
}

func (mod *Module) Parent() plugin.ResolvedModule  { return mod.parent }
func (mod *Module) Sources() []plugin.SourceModule { return mod.sources }
func (mod *Module) ID() plugin.ID                  { return mod.id }

func (mod *Module) Name() string {
	return mod.sources[0].Modules[len(mod.sources[0].Modules)-1].GetName()
}

func (mod *Module) ShortDescription(l *i18n.Localizer) string {
	for _, source := range mod.sources {
		parent := source.Modules[len(source.Modules)-1]

		if desc := parent.GetShortDescription(l); len(desc) > 0 {
			return desc
		}
	}

	return ""
}

func (mod *Module) LongDescription(l *i18n.Localizer) string {
	for _, source := range mod.sources {
		parent := source.Modules[len(source.Modules)-1]

		if desc := parent.GetLongDescription(l); len(desc) > 0 {
			return desc
		}
	}

	return mod.ShortDescription(l)
}

func (mod *Module) IsHidden() bool                     { return mod.hidden }
func (mod *Module) Commands() []plugin.ResolvedCommand { return mod.commands }
func (mod *Module) Modules() []plugin.ResolvedModule   { return mod.modules }

func (mod *Module) FindCommand(name string) plugin.ResolvedCommand {
	return findCommand(mod.commands, strings.Trim(name, shared.Whitespace), true)
}

func (mod *Module) FindModule(name string) plugin.ResolvedModule {
	return findModule(mod.modules, strings.Trim(name, shared.Whitespace))
}
