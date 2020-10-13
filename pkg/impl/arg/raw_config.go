package arg

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

type raw struct {
	desc i18nutil.Text
}

// Raw is a plugin.ArgConfig that returns the arguments as
var Raw = new(raw)

// RawWithDescription creates a argument config for raw arguments, that uses
// the passed description as argument config.
func RawWithDescription(description string) plugin.ArgConfig {
	return &raw{
		desc: i18nutil.NewText(description),
	}
}

// RawWithDescription creates a argument config for raw arguments, that uses
// the passed description as argument config.
func RawWithDescriptionl(description i18n.Config) plugin.ArgConfig {
	return &raw{
		desc: i18nutil.NewTextl(description),
	}
}

// RawWithDescriptionlt creates a argument config for raw arguments, that uses
// the passed description as argument config.
func RawWithDescriptionlt(description i18n.Term) plugin.ArgConfig {
	return &raw{
		desc: i18nutil.NewTextl(description.AsConfig()),
	}
}

func (r raw) Parse(args string, _ *state.State, _ *plugin.Context) (plugin.Args, plugin.Flags, error) {
	return []interface{}{args}, nil, nil
}

func (r raw) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	desc, err := r.desc.Get(l)
	if err != nil || desc == "" {
		return nil
	}

	return []plugin.ArgsInfo{
		{
			Required: []plugin.ArgInfo{
				{
					Name:        "",
					Type:        plugin.ArgTypeRaw,
					Description: desc,
				},
			},
		},
	}
}
