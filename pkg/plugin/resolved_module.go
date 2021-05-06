package plugin

import (
	"sort"
	"strings"

	"github.com/mavolin/adam/pkg/i18n"
)

type (
	// ResolvedModule is a resolved module as returned by a Provider.
	// In contrast to the regular Module abstraction, ResolvedModule's plugins
	// reflect the plugins provided by all modules with the same ID, i.e. a
	// plugin with the same name provided through different bot.PluginProvider.
	ResolvedModule struct {
		// Parent is the parent of this module.
		// If the module is top-level Parent will be nil.
		Parent *ResolvedModule

		// Sources contains the Modules this module is based upon.
		// Sources[0] will contain the built in module.
		// If there is no built-in module, Sources[0] will be nil.
		//
		// If the module is top-level, this will be empty.
		Sources []SourceModule

		// ID is the identifier of the module.
		ID ID
		// Name is the name of the module.
		Name string

		// Hidden specifies if all Sources are hidden.
		// A source module is considered hidden if all of it's commands and
		// modules are hidden as well.
		Hidden bool

		// Commands are the subcommands of the module.
		// They are sorted in ascending order by name.
		Commands []*ResolvedCommand
		// Modules are the submodules of the module.
		// They are sorted in ascending order by name.
		Modules []*ResolvedModule
	}

	// SourceModule contains the parent Modules of a ResolvedModule.
	SourceModule struct {
		// ProviderName is the name of the plugin provider that
		// provided the module.
		ProviderName string
		// Modules contains the parents of the ResolvedModule.
		// They are sorted in ascending order from the most distant to the
		// closest parent.
		Modules []Module
	}
)

// GenerateResolvedModules generates ResolvedModules from the passed
// Repositories.
func GenerateResolvedModules(repos []Repository) []*ResolvedModule {
	if len(repos) == 0 {
		return nil
	}

	rmod := make([]*ResolvedModule, 0, len(repos[0].Modules))

	var mergeLen int

	for _, repo := range repos {
		mergeLen += len(repo.Modules)
	}

	smods := make([]SourceModule, mergeLen)

	i := 0

	for _, repo := range repos {
		for _, mod := range repo.Modules {
			smods[i] = SourceModule{
				ProviderName: repo.ProviderName,
				Modules:      []Module{mod},
			}

			i++
		}
	}

	for _, sm := range sortSourceModules(smods) {
		// create a ResolvedModule for every SourceModule we merged.
		rmod = append(rmod, generateResolvedModule(nil, sm, repos))
	}

	return rmod
}

// sortSourceModules sorts the passed SourceModules into slices of modules with
// the same name.
// To preserve the order of the providers, smods should be sorted by providers.
//
// The outer slice contains slices of modules with the same name.
func sortSourceModules(smods []SourceModule) [][]SourceModule {
	// assume a maximum length of len(smods) and use that as cap
	sorted := make([][]SourceModule, 0, len(smods))

	for _, smod := range smods {
		// the module to compare the name against
		cmp := smod.Modules[len(smod.Modules)-1]

		// search for the index where the cmp module belongs lexicographically
		i := sort.Search(len(sorted), func(i int) bool {
			return sorted[i][0].Modules[0].GetName() >= cmp.GetName()
		})

		insert := SourceModule{
			ProviderName: smod.ProviderName,
			Modules:      smod.Modules,
		}

		if i == len(sorted) { // append
			sorted = append(sorted, []SourceModule{insert})
		} else { // insert within the range of the slice
			searchModule := sorted[i]

			// check if we have a module with the same name, and if so extend the list of source
			// modules
			if searchModule[0].Modules[0].GetName() == insert.Modules[0].GetName() {
				sorted[i] = append(sorted[i], insert)
			} else { // otherwise put the module behind the module found by Search
				sorted = append(sorted, nil) // make space for another element
				copy(sorted[i+1:], sorted[i:])

				sorted[i] = []SourceModule{insert}
			}
		}
	}

	return sorted
}

// generateResolvedModule generates a ResolvedModule from the passed
// SourceModules with the passed parent.
// The passed SourceModules represent a group of Modules with the same name in
// the passed parent.
//
// The passed Repositories will be used to determine the Defaults of the
// subcommands.
func generateResolvedModule(parent *ResolvedModule, smods []SourceModule, repos []Repository) *ResolvedModule {
	if len(smods) == 0 {
		return nil
	}

	// used for meta info
	referenceModule := smods[0].Modules[len(smods[0].Modules)-1]

	rmod := &ResolvedModule{
		Parent:  parent,
		Sources: smods,
		Name:    referenceModule.GetName(),
	}

	if parent == nil {
		//goland:noinspection GoRedundantConversion // GoLand is buggy, we need that conversion
		rmod.ID += "." + ID(referenceModule.GetName())
	} else {
		rmod.ID = parent.ID + ID("."+referenceModule.GetName())
	}

	fillSubmodules(rmod, repos)
	fillSubcommands(rmod)

	rmod.Hidden = true

	// mark module as hidden, unless it has one visible child
	for _, s := range rmod.Sources {
		smod := s.Modules[len(s.Modules)-1]

		for _, cmd := range smod.Commands() {
			if !cmd.IsHidden() {
				rmod.Hidden = false
				break
			}
		}

		for _, m := range smod.Modules() {
			if !rmod.FindModule(m.GetName()).Hidden {
				rmod.Hidden = false
				break
			}
		}
	}

	return rmod
}

// fillSubmodules fills the Modules field of the passed parent module.
// It generates the RegisteredModules from the parents Sources.
func fillSubmodules(parent *ResolvedModule, repos []Repository) {
	var maxLen int

	for _, smod := range parent.Sources {
		// get the number of modules in every source module
		maxLen += len(smod.Modules[len(smod.Modules)-1].Modules())
	}

	if maxLen == 0 {
		parent.Modules = nil
		return
	}

	// source modules of the modules of the parent
	subSMods := make([]SourceModule, maxLen)

	i := 0

	for _, smod := range parent.Sources { // go over all source modules
		// closest parent module
		parentSource := smod.Modules[len(smod.Modules)-1]

		// and range over the modules of the closest parent module
		for _, mod := range parentSource.Modules() {
			subSMods[i] = SourceModule{
				ProviderName: smod.ProviderName,
				// append the new inner module to the original source modules
				Modules: append(smod.Modules, mod),
			}

			i++
		}
	}

	sortedSmods := sortSourceModules(subSMods)

	parent.Modules = make([]*ResolvedModule, len(sortedSmods))

	for fillLen, smod := range sortedSmods {
		rmod := generateResolvedModule(parent, smod, repos)

		i := sort.Search(fillLen, func(i int) bool {
			return parent.Modules[i].Name >= rmod.Name
		}) // find insert index

		if i == fillLen { // append
			parent.Modules[i] = rmod
		} else { // insert
			copy(parent.Modules[i+1:], parent.Modules[i:])
			parent.Modules[i] = rmod
		}
	}
}

// fillSubcommands fills the Commands field of the passed parent module with
// the commands found in the parents Sources.
func fillSubcommands(parent *ResolvedModule) {
	var maxLen int

	for _, smod := range parent.Sources {
		// get the number of commands in every source module
		maxLen += len(smod.Modules[len(smod.Modules)-1].Commands())
	}

	if maxLen == 0 {
		parent.Commands = nil
		return
	}

	// preallocate the maximum possible amount of commands
	parent.Commands = make([]*ResolvedCommand, 0, maxLen)

	// set of aliases already used
	usedAliases := make(map[string]struct{}, maxLen)

	for _, smod := range parent.Sources {
		// generate RegisteredCommands for the current provider
		insertCmds := generateResolvedCommands(parent, smod)

		for _, rcmd := range insertCmds {
			rcmd.parent = &parent

			// remove duplicate aliases
			for i := 0; i < len(rcmd.Aliases); i++ {
				alias := rcmd.Aliases[i]

				if _, ok := usedAliases[alias]; ok { // alias is already in use, remove it
					rcmd.Aliases = append(rcmd.Aliases[:i], rcmd.Aliases[i+1:]...)
					i--
				} else { // alias unused, all good
					usedAliases[alias] = struct{}{}
				}
			}

			i := sort.Search(len(parent.Commands), func(i int) bool {
				return parent.Commands[i].Name >= rcmd.Name
			}) // find the insert index

			if len(parent.Commands) == i {
				parent.Commands = append(parent.Commands, rcmd)
			} else {
				if parent.Commands[i].Name == rcmd.Name {
					continue // skip if duplicate name
				}

				// otherwise insert

				parent.Commands = append(parent.Commands, rcmd) // make space for a new element
				copy(parent.Commands[i+1:], parent.Commands[i:])

				parent.Commands[i] = rcmd
			}
		}
	}
}

func generateResolvedCommands(parent *ResolvedModule, smod SourceModule) []*ResolvedCommand { //nolint:funlen
	var id ID

	for _, p := range smod.Modules {
		id += ID("." + p.GetName())
	}

	// get the commands of the innermost parent
	cmds := smod.Modules[len(smod.Modules)-1].Commands()
	rcmds := make([]*ResolvedCommand, len(cmds))

	for i, cmd := range cmds {
		rcmd := &ResolvedCommand{
			parent:        &parent,
			ID:            id + ID("."+cmd.GetName()),
			Source:        cmd,
			SourceParents: smod.Modules,
			ProviderName:  smod.ProviderName,
			Name:          cmd.GetName(),
			Args:          cmd.GetArgs(),
			Hidden:        cmd.IsHidden(),
			ChannelTypes:  cmd.GetChannelTypes(),
			Throttler:     cmd.GetThrottler(),
		}

		if rcmd.ChannelTypes == 0 {
			rcmd.ChannelTypes = AllChannels
		}

		if aliases := cmd.GetAliases(); aliases != nil {
			rcmd.Aliases = make([]string, len(aliases))
			copy(rcmd.Aliases, aliases)
		}

		rcmds[i] = rcmd
	}

	return rcmds
}

// ShortDescription returns an optional one-sentence description of the
// module.
func (m *ResolvedModule) ShortDescription(l *i18n.Localizer) string {
	for _, mod := range m.Sources {
		parent := mod.Modules[len(mod.Modules)-1]

		if desc := parent.GetShortDescription(l); len(desc) > 0 {
			return desc
		}
	}

	return ""
}

// LongDescription returns an option thorough description of the
// module.
//
// If the module only provides a short description, that will be
func (m *ResolvedModule) LongDescription(l *i18n.Localizer) string {
	for _, mod := range m.Sources {
		parent := mod.Modules[len(mod.Modules)-1]

		if desc := parent.GetLongDescription(l); len(desc) > 0 {
			return desc
		}
	}

	for _, mod := range m.Sources {
		parent := mod.Modules[len(mod.Modules)-1]

		if desc := parent.GetShortDescription(l); len(desc) > 0 {
			return desc
		}
	}

	return ""
}

// FindCommand finds the command with the given name inside this module.
// A name can either be the actual name of a command, or an alias.
func (m *ResolvedModule) FindCommand(name string) *ResolvedCommand {
	name = strings.TrimSpace(name)

	// fast path, if not searching for alias
	i := sort.Search(len(m.Commands), func(i int) bool {
		return m.Commands[i].Name >= name
	})

	if i < len(m.Commands) && m.Commands[i].Name == name {
		return m.Commands[i]
	}

	// no command with matching name, try matching alias
	for _, c := range m.Commands {
		for _, alias := range c.Aliases {
			if alias == name {
				return c
			}
		}
	}

	return nil
}

// FindModule finds the module with the given name inside the module.
func (m *ResolvedModule) FindModule(name string) *ResolvedModule {
	name = strings.TrimSpace(name)

	i := sort.Search(len(m.Modules), func(i int) bool {
		return m.Modules[i].Name >= name
	})

	if i < len(m.Modules) && m.Modules[i].Name == name {
		return m.Modules[i]
	}

	return nil
}
