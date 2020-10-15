package arg

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
)

// Switch is the type used for bool flags.
// If the flag gets set, Switch returns true.
// It cannot be used as an argument type.
//
// Switch flags cannot be used as multi flags.
var Switch Type = new(typeSwitch)

type typeSwitch struct{}

func (s typeSwitch) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(switchName) // we have a fallback
	return name
}

func (s typeSwitch) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(switchDescription) // we have a fallback
	return desc
}

func (s typeSwitch) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	return nil, errors.NewWithStack("arg: called Switch.Parse")
}

func (s typeSwitch) Default() interface{} {
	return false
}
