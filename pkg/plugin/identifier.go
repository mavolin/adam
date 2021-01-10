package plugin

import (
	"strings"
)

// Identifier is the unique identifier of a plugin.
// The root/base is '.'.
// All plugins are dot-separated, e.g. '.mod.ban'.
type Identifier string

const whitespace = " \t\n"

// NewIdentifierFromInvoke creates a new Identifier from the passed invoke.
func NewIdentifierFromInvoke(invoke string) Identifier {
	invoke = strings.Trim(invoke, whitespace)

	var b strings.Builder
	b.Grow(len(invoke) + 1)
	b.WriteRune('.')

	var prevWhitespace bool

	for _, r := range invoke {
		if strings.ContainsRune(whitespace, r) {
			if !prevWhitespace {
				b.WriteRune('.')
				prevWhitespace = true
			}
		} else {
			prevWhitespace = false
			b.WriteRune(r)
		}
	}

	return Identifier(b.String())
}

// Parent returns the parent module of the plugin, or '.' if this Identifier
// already represents root.
//
// If the Identifier is invalid, Parent returns an empty string.
func (id Identifier) Parent() Identifier {
	if id == "." {
		return id
	}

	i := strings.LastIndex(string(id), ".")
	if i == -1 {
		return ""
	} else if i == 0 { // parent is root
		i = 1
	}

	return id[:i]
}

// All returns a slice of all parents including root and the identifier itself
// starting with root.
//
// If the Identifier is invalid, All returns nil.
func (id Identifier) All() []Identifier {
	if id.IsRoot() {
		return []Identifier{"."}
	}

	pluginCount := strings.Count(string(id), ".")
	if pluginCount == 0 {
		return nil
	}

	parents := make([]Identifier, pluginCount+1)

	parent := id

	for i := len(parents) - 1; i >= 0; i-- {
		parents[i] = parent

		parent = parent.Parent()
	}

	return parents
}

// IsRoot checks if the identifier is the root identifier.
func (id Identifier) IsRoot() bool {
	return id == "."
}

// NumParents returns the number of parents the plugin has.
//
// Returns a negative number, if the Identifier is invalid.
func (id Identifier) NumParents() int {
	return strings.Count(string(id), ".") - 1
}

// IsParent checks if the passed Identifier is a parent of this identifier.
func (id Identifier) IsParent(target Identifier) bool {
	return len(id) > len(target) && strings.HasPrefix(string(id), string(target))
}

// IsChild checks if the passed Identifier is a child of this identifier.
func (id Identifier) IsChild(target Identifier) bool {
	return len(id) < len(target) && strings.HasPrefix(string(target), string(id))
}

// AsInvoke returns the identifier as a prefixless command invoke.
//
// Returns "" if the Identifier is root or invalid.
//
// Example:
// 	.mod.ban -> mod ban
func (id Identifier) AsInvoke() string {
	if len(id) == 0 {
		return ""
	}

	return strings.ReplaceAll(string(id[1:]), ".", " ")
}

// Name returns the name of the plugin or "" if the Identifier is root or
// invalid.
func (id Identifier) Name() string {
	if len(id) <= 1 {
		return ""
	}

	i := strings.LastIndex(string(id), ".")
	return string(id[i+1:])
}
