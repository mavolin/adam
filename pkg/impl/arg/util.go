package arg

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func genArgsInfo(
	l *i18n.Localizer, rargs []RequiredArg, oargs []OptionalArg, flags []Flag, variadic bool,
) (info plugin.ArgsInfo) {
	info = plugin.ArgsInfo{
		Required:      make([]plugin.ArgInfo, len(rargs)),
		Optional:      make([]plugin.ArgInfo, len(oargs)),
		Flags:         make([]plugin.FlagInfo, len(flags)),
		FlagFormatter: func(name string) string { return "-" + name },
		Variadic:      variadic,
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
) (info plugin.ArgsInfo, ok bool) {
	info = plugin.ArgsInfo{
		Required:      make([]plugin.ArgInfo, len(rargs)),
		Optional:      make([]plugin.ArgInfo, len(oargs)),
		Flags:         make([]plugin.FlagInfo, len(flags)),
		FlagFormatter: func(name string) string { return "-" + name },
		Variadic:      variadic,
	}

	for i, arg := range rargs {
		info.Required[i], ok = requiredArgInfol(l, arg)
		if !ok {
			return plugin.ArgsInfo{}, false
		}
	}

	for i, arg := range oargs {
		info.Optional[i], ok = optionalArgInfol(l, arg)
		if !ok {
			return plugin.ArgsInfo{}, false
		}
	}

	for i, flag := range flags {
		info.Flags[i] = flagInfol(l, flag)
	}

	return info, true
}

func requiredArgInfol(l *i18n.Localizer, a LocalizedRequiredArg) (info plugin.ArgInfo, ok bool) {
	var err error

	info.Name, err = l.Localize(a.Name)
	if err != nil {
		return plugin.ArgInfo{}, false
	}

	info.Type = typeInfo(l, a.Type)

	if a.Description != nil {
		info.Description, _ = l.Localize(a.Description)
	}

	return info, true
}

func optionalArgInfol(l *i18n.Localizer, a LocalizedOptionalArg) (info plugin.ArgInfo, ok bool) {
	var err error

	info.Name, err = l.Localize(a.Name)
	if err != nil {
		return plugin.ArgInfo{}, false
	}

	info.Type = typeInfo(l, a.Type)

	if a.Description != nil {
		info.Description, _ = l.Localize(a.Description)
	}

	return info, true
}

func flagInfol(l *i18n.Localizer, f LocalizedFlag) (info plugin.FlagInfo) {
	info.Name = f.Name

	if len(f.Aliases) > 0 {
		info.Aliases = make([]string, len(f.Aliases))
		copy(info.Aliases, f.Aliases)
	}

	info.Type = typeInfo(l, f.Type)
	info.Multi = f.Multi

	if f.Description != nil {
		info.Description, _ = l.Localize(f.Description)
	}

	return info
}

func typeInfo(l *i18n.Localizer, t Type) plugin.ArgType {
	return plugin.ArgType{
		Name:        t.Name(l),
		Description: t.Description(l),
	}
}

// newArgumentError2 creates a new plugin.ArgumentError using the passed
// *i18n.Config.
// It adds the following additional placeholders: name, used_name, raw and
// position.
// If raw is longer than a 100 characters, it will be shortened.
func newArgumentError(
	cfg *i18n.Config, ctx *Context, placeholders map[string]interface{},
) *plugin.ArgumentError {
	placeholders = fillPlaceholders(placeholders, ctx)
	return plugin.NewArgumentErrorl(cfg.
		WithPlaceholders(placeholders))
}

// newArgumentError2 creates a new *plugin.ArgumentError and decides based
// on the passed Context which of the two *i18n.Configs to use.
// It adds the following additional placeholders: name, used_name, raw and
// position.
// If raw is longer than a 100 characters, it will be shortened.
func newArgumentError2(
	argConfig, flagConfig *i18n.Config, ctx *Context, placeholders map[string]interface{},
) *plugin.ArgumentError {
	placeholders = fillPlaceholders(placeholders, ctx)

	if ctx.Kind == KindArg {
		return plugin.NewArgumentErrorl(argConfig.
			WithPlaceholders(placeholders))
	}

	return plugin.NewArgumentErrorl(flagConfig.
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
