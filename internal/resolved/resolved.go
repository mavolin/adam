// Package resolved provides implementations for plugin.ResolvedCommand and
// plugin.ResolvedModule as well a plugin.Provider implementation that
// generates the plugins.
package resolved

import "github.com/mavolin/adam/pkg/plugin"

func replaceModuleProvider(
	parent plugin.ResolvedModule, mods []plugin.ResolvedModule, p *PluginProvider,
) []plugin.ResolvedModule {
	if len(mods) == 0 {
		return nil
	}

	cp := make([]plugin.ResolvedModule, len(mods))

	for i, mod := range mods {
		typedMod := mod.(*Module)

		modCp := *typedMod
		modCp.parent = parent
		modCp.commands = replaceCommandProvider(&modCp, modCp.commands, p)
		modCp.modules = replaceModuleProvider(&modCp, typedMod.modules, p)

		cp[i] = &modCp
	}

	return cp
}

func replaceCommandProvider(
	parent plugin.ResolvedModule, cmds []plugin.ResolvedCommand, p *PluginProvider,
) []plugin.ResolvedCommand {
	if len(cmds) == 0 {
		return nil
	}

	cp := make([]plugin.ResolvedCommand, len(cmds))

	for i, cmd := range cmds {
		typedCmd := cmd.(*Command)

		cmdCp := *typedCmd
		cmdCp.provider = p
		cmdCp.parent = parent

		cp[i] = &cmdCp
	}

	return cp
}
