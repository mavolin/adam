package arg

import (
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

const whitespace = " \t\n"

func genArgsInfo(
	l *i18n.Localizer, rargs []RequiredArg, oargs []OptionalArg, flags []Flag, variadic bool,
) (info plugin.ArgsInfo) {
	info = plugin.ArgsInfo{
		Required: make([]plugin.ArgInfo, len(rargs)),
		Optional: make([]plugin.ArgInfo, len(oargs)),
		Flags:    make([]plugin.FlagInfo, len(flags)),
		Variadic: variadic,
	}

	for i, arg := range rargs {
		info.Required[i] = plugin.ArgInfo{
			Name:        arg.Name,
			Type:        typeInfo(l, arg.Type),
			Description: arg.Description,
		}
	}

	for i, arg := range oargs {
		info.Optional[i] = plugin.ArgInfo{
			Name:        arg.Name,
			Type:        typeInfo(l, arg.Type),
			Description: arg.Description,
		}
	}

	for i, flag := range flags {
		info.Flags[i] = plugin.FlagInfo{
			Name:        flag.Name,
			Type:        typeInfo(l, flag.Type),
			Description: flag.Description,
			Multi:       flag.Multi,
		}

		if len(flag.Aliases) > 0 {
			info.Flags[i].Aliases = make([]string, len(flag.Aliases))
			copy(info.Flags[i].Aliases, flag.Aliases)
		}
	}

	return info
}

func genArgsInfol(
	l *i18n.Localizer,
	rargs []LocalizedRequiredArg, oargs []LocalizedOptionalArg, flags []LocalizedFlag, variadic bool,
) (plugin.ArgsInfo, error) {
	info := plugin.ArgsInfo{
		Required: make([]plugin.ArgInfo, len(rargs)),
		Optional: make([]plugin.ArgInfo, len(oargs)),
		Flags:    make([]plugin.FlagInfo, len(flags)),
		Variadic: variadic,
	}

	var err error

	for i, arg := range rargs {
		info.Required[i], err = requiredArgInfol(l, arg)
		if err != nil {
			return plugin.ArgsInfo{}, err
		}
	}

	for i, arg := range oargs {
		info.Optional[i], err = optionalArgInfol(l, arg)
		if err != nil {
			return plugin.ArgsInfo{}, err
		}
	}

	for i, flag := range flags {
		info.Flags[i], err = flagInfol(l, flag)
		if err != nil {
			return plugin.ArgsInfo{}, err
		}
	}

	return info, nil
}

func requiredArgInfol(l *i18n.Localizer, a LocalizedRequiredArg) (info plugin.ArgInfo, err error) {
	info.Name, err = l.Localize(a.Name)
	if err != nil {
		return plugin.ArgInfo{}, err
	}

	info.Type = typeInfo(l, a.Type)

	info.Description, err = l.Localize(a.Description)
	return
}

func optionalArgInfol(l *i18n.Localizer, a LocalizedOptionalArg) (info plugin.ArgInfo, err error) {
	info.Name, err = l.Localize(a.Name)
	if err != nil {
		return
	}

	info.Type = typeInfo(l, a.Type)

	info.Description, err = l.Localize(a.Description)
	return
}

func flagInfol(l *i18n.Localizer, f LocalizedFlag) (info plugin.FlagInfo, err error) {
	info.Name = f.Name

	if len(f.Aliases) > 0 {
		info.Aliases = make([]string, len(f.Aliases))
		copy(info.Aliases, f.Aliases)
	}

	info.Type = typeInfo(l, f.Type)
	info.Multi = f.Multi

	info.Description, err = l.Localize(f.Description)
	return
}

func typeInfo(l *i18n.Localizer, t Type) plugin.TypeInfo {
	return plugin.TypeInfo{
		Name:        t.Name(l),
		Description: t.Description(l),
	}
}

// newArgParsingErr2 creates a new errors.ArgumentError using the passed
// i18n.Config.
// It adds the following additional placeholders: name, used_name, raw and
// position.
// If raw is longer than a 100 characters, it will be shortened.
func newArgParsingErr(
	cfg *i18n.Config, ctx *Context, placeholders map[string]interface{},
) *errors.ArgumentError {
	placeholders = fillPlaceholders(placeholders, ctx)
	return errors.NewArgumentErrorl(cfg.
		WithPlaceholders(placeholders))
}

// newArgParsingErr2 creates a new errors.ArgumentError and decides based
// on the passed Context which of the two i18n.Configs to use.
// It adds the following additional placeholders: name, used_name, raw and
// position.
// If raw is longer than a 100 characters, it will be shortened.
func newArgParsingErr2(
	argConfig, flagConfig *i18n.Config, ctx *Context, placeholders map[string]interface{},
) *errors.ArgumentError {
	placeholders = fillPlaceholders(placeholders, ctx)

	if ctx.Kind == KindArg {
		return errors.NewArgumentErrorl(argConfig.
			WithPlaceholders(placeholders))
	}

	return errors.NewArgumentErrorl(flagConfig.
		WithPlaceholders(placeholders))
}

func fillPlaceholders(placeholders map[string]interface{}, ctx *Context) map[string]interface{} {
	if placeholders == nil {
		placeholders = make(map[string]interface{}, 4)
	}

	placeholders["name"] = ctx.Name
	placeholders["used_name"] = ctx.UsedName
	placeholders["position"] = ctx.Index + 1

	raw := ctx.Raw
	if len(raw) > 100 {
		raw = raw[:100]
	}
	placeholders["raw"] = raw

	return placeholders
}
