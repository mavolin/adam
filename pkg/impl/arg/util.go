package arg

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// newArgumentError2 creates a new plugin.ArgumentError using the passed
// *i18n.Config.
// It adds the following additional placeholders: name, used_name, raw and
// position.
// If raw is longer than a 100 characters, it will be shortened.
func newArgumentError(
	cfg *i18n.Config, ctx *plugin.ParseContext, placeholders map[string]interface{},
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
	argConfig, flagConfig *i18n.Config, ctx *plugin.ParseContext, placeholders map[string]interface{},
) *plugin.ArgumentError {
	placeholders = fillPlaceholders(placeholders, ctx)

	if ctx.Kind == plugin.KindArg {
		return plugin.NewArgumentErrorl(argConfig.
			WithPlaceholders(placeholders))
	}

	return plugin.NewArgumentErrorl(flagConfig.
		WithPlaceholders(placeholders))
}

func fillPlaceholders(placeholders map[string]interface{}, ctx *plugin.ParseContext) map[string]interface{} {
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
