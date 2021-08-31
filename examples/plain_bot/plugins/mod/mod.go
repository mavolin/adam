// Package mod provides the moderation module.
package mod

import (
	"github.com/mavolin/adam/examples/plain_bot/plugins/mod/ban"
	"github.com/mavolin/adam/examples/plain_bot/plugins/mod/kick"
	"github.com/mavolin/adam/pkg/impl/module"
	"github.com/mavolin/adam/pkg/plugin"
)

// New creates a new moderation module.
func New() plugin.Module {
	m := module.New(module.Meta{
		Name:             "mod",
		ShortDescription: "Moderate a server.",
		LongDescription:  "The moderation module provides utilities for moderating your server easily.",
	})

	m.AddCommand(ban.New())
	m.AddCommand(kick.New())

	return m
}
