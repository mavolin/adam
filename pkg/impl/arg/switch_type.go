package arg

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Switch is the type used for bool flags.
// If the flag gets set, Switch returns true.
// It cannot be used as an argument type.
// Due to the special nature of this type, plugin.ArgParsers must handle it
// specially, i.e. expect no content for it.
//
// Switch flags cannot be used as multi flags.
var Switch plugin.ArgType = new(typeSwitch)

type typeSwitch struct{}

func (s typeSwitch) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(switchName) // we have a fallback
	return name
}

func (s typeSwitch) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(switchDescription) // we have a fallback
	return desc
}

func (s typeSwitch) Parse(*state.State, *plugin.ParseContext) (interface{}, error) {
	return nil, errors.NewWithStack("arg: called Switch.Parse")
}

func (s typeSwitch) GetDefault() interface{} {
	return false
}
