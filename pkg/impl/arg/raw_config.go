package arg

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Raw is the type used for unlocalized raw arguments.
type Raw struct {
	// Name is the name of the single argument.
	//
	// This field is required.
	Name string
	// TypeName is the name of the type of the argument.
	TypeName string
	// Description is the description used for the argument.
	Description string

	// Optional, if set to true, won't enforce a minimum argument length of 0.
	Optional bool
}

func (r Raw) Parse(args string, _ *state.State, _ *plugin.Context) (plugin.Args, plugin.Flags, error) {
	if !r.Optional && len(args) == 0 {
		return nil, nil, plugin.NewArgumentErrorl(notEnoughArgsError)
	}

	return []interface{}{args}, nil, nil
}

func (r Raw) Info(*i18n.Localizer) []plugin.ArgsInfo {
	ai := []plugin.ArgInfo{
		{
			Name:        r.Name,
			Type:        plugin.TypeInfo{Name: r.TypeName},
			Description: r.Description,
		},
	}

	if r.Optional {
		return []plugin.ArgsInfo{{Optional: ai}}
	}

	return []plugin.ArgsInfo{{Required: ai}}
}

// LocalizedRaw is the type used for localized raw arguments.
type LocalizedRaw struct {
	// Name is the name of the single argument.
	//
	// This field is required.
	Name *i18n.Config
	// TypeName is the name of the type of the argument.
	TypeName *i18n.Config
	// Description is the description used for the argument.
	Description *i18n.Config

	// Optional, if set to true, won't enforce a minimum argument length of 0.
	Optional bool
}

func (r LocalizedRaw) Parse(args string, _ *state.State, _ *plugin.Context) (plugin.Args, plugin.Flags, error) {
	if !r.Optional && len(args) == 0 {
		return nil, nil, plugin.NewArgumentErrorl(notEnoughArgsError)
	}

	return []interface{}{args}, nil, nil
}

func (r LocalizedRaw) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	name, err := l.Localize(r.Name)
	if err != nil {
		return nil
	}

	var typeName string
	if r.TypeName != nil {
		typeName, err = l.Localize(r.TypeName)
		if err != nil {
			return nil
		}
	}

	var desc string
	if r.Description != nil {
		desc, err = l.Localize(r.Description)
		if err != nil {
			return nil
		}
	}

	ai := []plugin.ArgInfo{{Name: name, Description: desc}}

	if len(typeName) > 0 {
		ai[0].Type = plugin.TypeInfo{Name: typeName}
	}

	if r.Optional {
		return []plugin.ArgsInfo{{Optional: ai}}
	}

	return []plugin.ArgsInfo{{Required: ai}}
}
