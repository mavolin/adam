package resolved

import (
	"sort"
	"strings"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/plugin"
)

func searchCommand(rcmds []plugin.ResolvedCommand, name string) int {
	return sort.Search(len(rcmds), func(i int) bool {
		return rcmds[i].Name() >= name
	})
}

func findCommand(rcmds []plugin.ResolvedCommand, name string, searchAliases bool) plugin.ResolvedCommand {
	i := searchCommand(rcmds, name)

	if i < len(rcmds) && rcmds[i].Name() == name {
		return rcmds[i]
	}

	if !searchAliases {
		return nil
	}

	for _, cmd := range rcmds {
		for _, alias := range cmd.Aliases() {
			if alias == name {
				return cmd
			}
		}
	}

	return nil
}

func searchModule(rmods []plugin.ResolvedModule, name string) int {
	return sort.Search(len(rmods), func(i int) bool {
		return rmods[i].Name() >= name
	})
}

func findModule(rmods []plugin.ResolvedModule, name string) plugin.ResolvedModule {
	i := searchModule(rmods, name)

	if i < len(rmods) && rmods[i].Name() == name {
		return rmods[i]
	}

	return nil
}

func insertCommand(rcmds []plugin.ResolvedCommand, rcmd plugin.ResolvedCommand, i int) []plugin.ResolvedCommand {
	if i < 0 {
		i = searchCommand(rcmds, rcmd.Name())
	}

	if i >= len(rcmds) {
		return append(rcmds, rcmd)
	}

	rcmds = append(rcmds, rcmd)

	copy(rcmds[i+1:], rcmds[i:])
	rcmds[i] = rcmd

	return rcmds
}

func insertModule(rmods []plugin.ResolvedModule, rmod plugin.ResolvedModule, i int) []plugin.ResolvedModule {
	if i < 0 {
		i = searchModule(rmods, rmod.Name())
	}

	if i >= len(rmods) {
		return append(rmods, rmod)
	}

	if rmods[i].Name() == rmod.Name() {
		rmods[i] = rmod
		return rmods
	}

	rmods = append(rmods, rmod)

	copy(rmods[i+1:], rmods[i:])
	rmods[i] = rmod

	return rmods
}

func firstWord(s string) (string, string) {
	for i, r := range s {
		if strings.ContainsRune(shared.Whitespace, r) {
			return s[:i], strings.TrimLeft(s[i+1:], shared.Whitespace)
		}
	}

	return s, ""
}
