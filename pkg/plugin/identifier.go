package plugin

import "strings"

// Identifier is the unique identifier of a plugin.
// The root/base is '.'.
// All plugins are dot-separated, e.g. '.mod.ban'
type Identifier string

// Parent returns the parent module of the plugin or '.' if this Identifier
// already represents root.
func (id Identifier) Parent() Identifier {
	if id == "." {
		return id
	}

	i := strings.LastIndex(string(id), ".")
	if i == 0 { // parent is root
		i = 1
	}

	return id[:i]
}

// IsRoot checks if the identifier is the root identifier.
func (id Identifier) IsRoot() bool {
	return id == "."
}
