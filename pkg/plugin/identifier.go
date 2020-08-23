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

// IsParent checks if the passed Identifier is a parent of this identifier.
func (id Identifier) IsParent(target Identifier) bool {
	return id > target && strings.HasPrefix(string(id), string(target))
}

// IsChild checks if the passed Identifier is a child of this identifier.
func (id Identifier) IsChild(target Identifier) bool {
	return target > id && strings.HasPrefix(string(target), string(id))
}

// AsCommandInvoke returns the identifier as a prefixless command invoke.
//
// Returns "" if the Identifier is root.
//
// Example:
// 	.mod.ban -> mod ban
func (id Identifier) AsCommandInvoke() string {
	return strings.ReplaceAll(string(id)[1:], ".", " ")
}
