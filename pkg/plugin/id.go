package plugin

import (
	"strings"

	"github.com/mavolin/adam/internal/shared"
)

// ID is the unique identifier of a plugin.
// The root/base is '.'.
// All plugins are dot-separated, e.g. '.mod.ban'.
type ID string

// RootID is the ID representing root, i.e. the parent of all plugins.
const RootID ID = "."

// NewIDFromInvoke creates a new ID from the passed invoke.
func NewIDFromInvoke(invoke string) ID {
	invoke = strings.Trim(invoke, shared.Whitespace)
	plugins := strings.FieldsFunc(invoke, func(r rune) bool {
		return strings.ContainsRune(shared.Whitespace, r)
	})

	return ID("." + strings.Join(plugins, "."))
}

// Parent returns the parent module of the plugin, or '.' if this ID
// already represents root.
//
// If the ID is invalid, Parent returns an empty string.
func (id ID) Parent() ID {
	if id == RootID {
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
// If the ID is invalid, All returns nil.
func (id ID) All() []ID {
	if id.IsRoot() {
		return []ID{RootID}
	}

	pluginCount := strings.Count(string(id), ".")
	if pluginCount == 0 {
		return nil
	}

	parents := make([]ID, pluginCount+1)

	parent := id

	for i := len(parents) - 1; i >= 0; i-- {
		parents[i] = parent

		parent = parent.Parent()
	}

	return parents
}

// IsRoot checks if the identifier is the root identifier.
func (id ID) IsRoot() bool {
	return id == RootID
}

// NumParents returns the number of parents the plugin has.
//
// Returns a negative number, if the ID is invalid.
func (id ID) NumParents() int {
	return strings.Count(string(id), ".") - 1
}

// IsParentOf checks if the this ID is a parent of target.
func (id ID) IsParentOf(target ID) bool {
	return len(id) < len(target) && strings.HasPrefix(string(target), string(id))
}

// IsChildOf checks if the this ID is a child of target.
func (id ID) IsChildOf(target ID) bool {
	return len(id) > len(target) && strings.HasPrefix(string(id), string(target))
}

// AsInvoke returns the identifier as a prefixless command invoke.
//
// Returns "" if the ID is root or invalid.
func (id ID) AsInvoke() string {
	if id == "" {
		return ""
	}

	return strings.ReplaceAll(string(id[1:]), ".", " ")
}

// Name returns the name of the plugin or "" if the ID is root or
// invalid.
func (id ID) Name() string {
	if len(id) <= 1 {
		return ""
	}

	i := strings.LastIndex(string(id), ".")
	return string(id[i+1:])
}
