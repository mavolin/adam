package arg

import (
	"reflect"
	"strings"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// ChoiceCaseSensitive is a global flag that defines whether choices should be
// case sensitive.
// Defaults to false.
var ChoiceCaseSensitive = false

type (
	// Choice is an unlocalized enum type.
	Choice []ChoiceElement

	// ChoiceElement is an element in a choice and represents a single value.
	ChoiceElement struct {
		// Name is the name of the element.
		Name string
		// Aliases are optional aliases for the element.
		Aliases []string
		// Value is the value the element. corresponds to.
		// If this is nil, the name of the choice will be used.
		Value interface{}
	}
)

var _ Type = Choice{}

func (c Choice) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(choiceName) // we have a fallback
	return name
}

func (c Choice) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(choiceDescription)
	return desc
}

func (c Choice) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	for _, e := range c {
		if (ChoiceCaseSensitive && e.Name == ctx.Raw) || strings.EqualFold(e.Name, ctx.Raw) {
			if e.Value == nil {
				return e.Name, nil
			}

			return e.Value, nil
		}

		for _, alias := range e.Aliases {
			if (ChoiceCaseSensitive && alias == ctx.Raw) || strings.EqualFold(alias, ctx.Raw) {
				if e.Value == nil {
					return e.Name, nil
				}

				return e.Value, nil
			}
		}
	}

	return nil, newArgParsingErr(choiceInvalidError, ctx, nil)
}

// Default tries to derive the default type from the value of the first choice.
// If the choice is empty, Default returns nil.
func (c Choice) Default() interface{} {
	if len(c) > 0 {
		if c[0].Value == nil {
			return "" // fallback to Name's value, which is of type string
		}

		t := reflect.TypeOf(c[0].Value)
		return reflect.Zero(t).Interface()
	}

	return nil
}

type (
	// LocalizedChoice is an localized enum type.
	LocalizedChoice []LocalizedChoiceElement

	// LocalizedChoiceElement is an element in a localized choice and
	// represents a single value.
	LocalizedChoiceElement struct {
		// Names are the names used for the element.
		Names []*i18n.Config
		// Value is the value the element corresponds to.
		Value interface{}
	}
)

func (c LocalizedChoice) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(choiceName) // we have a fallback
	return name
}

func (c LocalizedChoice) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(choiceDescription)
	return desc
}

func (c LocalizedChoice) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	for _, e := range c {
		for _, nameConfig := range e.Names {
			name, err := ctx.Localizer.Localize(nameConfig)
			if err != nil {
				return nil, err
			}

			if (ChoiceCaseSensitive && name == ctx.Raw) || strings.EqualFold(name, ctx.Raw) {
				return e.Value, nil
			}
		}
	}

	return nil, newArgParsingErr(choiceInvalidError, ctx, nil)
}

// Default tries to derive the default type from the value of the first choice.
// If the choice is empty, Default returns nil.
func (c LocalizedChoice) Default() interface{} {
	if len(c) > 0 && c[0].Value != nil {
		t := reflect.TypeOf(c[0].Value)
		return reflect.Zero(t).Interface()
	}

	return nil
}
