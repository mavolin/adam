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

	invoke := parentInvoke + smod.GetName()
	if _, ok := provider.usedNames[invoke]; ok {
		return nil
	}

	provider.usedNames[invoke] = struct{}{}

	id := plugin.ID("." + smod.GetName())
	if parent != nil {
		id = parent.ID() + id
	}

	rmod := &Module{id: id, hidden: true}
	if parent != nil {
		// setting rmod.parent directly wouldn't work, as the nil's types would
		// mismatch
		rmod.parent = parent
	}

	if len(smod.GetModules()) > 0 {
		rmod.modules = make([]plugin.ResolvedModule, 0, len(smod.GetModules()))
	}

	rmod.sources = []plugin.SourceModule{{SourceName: sourceName}}

	if parent != nil {
		for _, parentSource := range parent.Sources() {
			if parentSource.SourceName == sourceName {
				rmod.sources[0].Modules = make([]plugin.Module, len(parentSource.Modules)+1)
				copy(rmod.sources[0].Modules, parentSource.Modules)
				rmod.sources[0].Modules[len(rmod.sources[0].Modules)-1] = smod

				break
			}
		}
	} else {
		rmod.sources[0].Modules = []plugin.Module{smod}
	}

	rmod.addCommands(provider, rmod.sources[0], smod.GetCommands())

	for _, subSmod := range smod.GetModules() {
		subRmod := newModule(rmod, provider, sourceName, subSmod)
		if !subRmod.IsHidden() {
			rmod.hidden = false
		}

		rmod.modules = insertModule(rmod.modules, subRmod, -1)
	}

	return rmod
}

func (mod *Module) update(provider *PluginProvider, sourceName string, smod plugin.Module) {
	if mod.Parent() != nil {
		for _, parentSource := range mod.Parent().Sources() {
			if parentSource.SourceName == sourceName {
				modules := make([]plugin.Module, len(parentSource.Modules)+1)
				copy(modules, parentSource.Modules)
				modules[len(modules)-1] = smod

				mod.sources = append(mod.sources, plugin.SourceModule{
					SourceName: sourceName,
					Modules:    modules,
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

	mod.addCommands(provider, mod.sources[len(mod.sources)-1], smod.GetCommands())

	for _, subSmod := range smod.GetModules() {
		i := searchModule(mod.modules, subSmod.GetName())
		if i < len(mod.modules) && mod.modules[i].Name() == smod.GetName() {
			subRmod := mod.modules[i].(*Module)
			subRmod.update(provider, sourceName, subSmod)

			if subRmod.IsHidden() {
				mod.hidden = false
			}
		} else {
			subRmod := newModule(mod, provider, sourceName, subSmod)
			if subRmod == nil {
				continue
			}

			mod.modules = insertModule(mod.modules, subRmod, i)

			if subRmod.IsHidden() {
				mod.hidden = false
			}
		}
	}
}

func (mod *Module) addCommands(provider *PluginProvider, sourceModule plugin.SourceModule, scmds []plugin.Command) {
	if len(scmds) == 0 {
		return
	}

	cp := make([]plugin.ResolvedCommand, len(mod.commands), len(mod.commands)+len(scmds))
	copy(cp, mod.commands)
	mod.commands = cp

	for _, scmd := range scmds {
		rcmd := newCommand(mod, provider, sourceModule.SourceName, sourceModule.Modules, scmd)
		if rcmd == nil {
			continue
		}

		if !rcmd.IsHidden() {
			mod.hidden = false
		}

		mod.commands = insertCommand(mod.commands, rcmd, -1)
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
