package arg

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// RawParser is a plugin.ArgParser that parses it's arguments literally.
//
// It requires that only one RequiredArgument or OptionalArgument is defined,
// otherwise it will panic.
// To use custom type naming, RawType or LocalizedRawType can be used.
var RawParser = new(raw)

type raw struct{}

func (r raw) Parse(args string, argConfig plugin.ArgConfig, s *state.State, ctx *plugin.Context) error {
	if rargs := argConfig.GetRequiredArgs(); len(rargs) == 1 {
		arg := rargs[0]

		parsed, err := arg.GetType().Parse(s, &plugin.ParseContext{
			Context:  ctx,
			Raw:      args,
			Name:     arg.GetName(ctx.Localizer),
			UsedName: arg.GetName(ctx.Localizer),
			Index:    0,
			Kind:     plugin.KindArg,
		})
		if err == nil {
			ctx.Args = plugin.Args{parsed}
		}

		return err
	}

	if oargs := argConfig.GetOptionalArgs(); len(oargs) == 1 {
		arg := oargs[0]

		if len(args) == 0 {
			parsed := arg.GetDefault()
			if parsed == nil {
				parsed = arg.GetType().GetDefault()
			}

			ctx.Args = plugin.Args{parsed}
			return nil
		}

		parsed, err := arg.GetType().Parse(s, &plugin.ParseContext{
			Context:  ctx,
			Raw:      args,
			Name:     arg.GetName(ctx.Localizer),
			UsedName: arg.GetName(ctx.Localizer),
			Index:    0,
			Kind:     plugin.KindArg,
		})
		if err == nil {
			ctx.Args = plugin.Args{parsed}
		}

		return err
	}

	panic("arg: RawParser: ArgConfig does not contain a single RequiredArg or a single OptionalArg")
}

func (r raw) FormatArgs(_ plugin.ArgConfig, args []string, _ map[string]string) string {
	if len(args) == 0 {
		return ""
	}

	return args[0]
}

func (r raw) FormatUsage(_ plugin.ArgConfig, args []string) string {
	if len(args) == 0 {
		return ""
	}

	return args[0]
}

func (r raw) FormatFlag(string) string {
	panic("arg.RawParser should not define flags")
}

// =============================================================================
// RawType
// =====================================================================================

// RawType is a unlocalized plugin.ArgType that allows specifying a custom name
// and description.
//
// Go type: string
type RawType struct {
	// Name is the name of the type.
	Name string
	// Description is the description of the type.
	Description string
}

func (t RawType) GetName(*i18n.Localizer) string        { return t.Name }
func (t RawType) GetDescription(*i18n.Localizer) string { return t.Description }

func (t RawType) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	return ctx.Raw, nil
}

func (t RawType) GetDefault() interface{} { return "" }

// LocalizedRawType is a localized plugin.ArgType that allows specifying a
// custom name and description.
//
// Go type: string
type LocalizedRawType struct {
	// Name is the name of the type.
	Name *i18n.Config
	// Description is the description of the type.
	Description *i18n.Config
}

func (t LocalizedRawType) GetName(l *i18n.Localizer) string {
	if name, err := l.Localize(t.Name); err == nil {
		return name
	}

	return ""
}
func (t LocalizedRawType) GetDescription(l *i18n.Localizer) string {
	if desc, err := l.Localize(t.Description); err == nil {
		return desc
	}

	return ""
}

func (t LocalizedRawType) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	return ctx.Raw, nil
}

func (t LocalizedRawType) GetDefault() interface{} { return "" }
