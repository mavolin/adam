package arg

import (
	"sort"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	emojiutil "github.com/mavolin/adam/pkg/utils/emoji"
)

const whitespace = " \t\n"

func genArgsInfo(
	l *i18n.Localizer, rargs []RequiredArg, oargs []OptionalArg, flags []Flag, variadic bool,
) (plugin.ArgsInfo, error) {
	info := plugin.ArgsInfo{
		Required: make([]plugin.ArgInfo, len(rargs)),
		Optional: make([]plugin.ArgInfo, len(oargs)),
		Flags:    make([]plugin.FlagInfo, len(flags)),
		Variadic: variadic,
	}

	var err error

	for i, arg := range rargs {
		info.Required[i], err = requiredArgInfo(arg, l)
		if err != nil {
			return plugin.ArgsInfo{}, err
		}
	}

	for i, arg := range oargs {
		info.Optional[i], err = optionalArgInfo(arg, l)
		if err != nil {
			return plugin.ArgsInfo{}, err

		}
	}

	for i, flag := range flags {
		info.Flags[i], err = flagInfo(flag, l)
		if err != nil {
			return plugin.ArgsInfo{}, err
		}
	}
	return info, nil
}

func requiredArgInfo(a RequiredArg, l *i18n.Localizer) (info plugin.ArgInfo, err error) {
	info.Name, err = a.Name.Get(l)
	if err != nil {
		return
	}

	var ok bool

	info.Type, ok = typeInfo(a.Type, l)
	if !ok {
		return
	}

	info.Description, err = a.Description.Get(l)
	return
}

func optionalArgInfo(a OptionalArg, l *i18n.Localizer) (info plugin.ArgInfo, err error) {
	info.Name, err = a.Name.Get(l)
	if err != nil {
		return
	}

	var ok bool

	info.Type, ok = typeInfo(a.Type, l)
	if !ok {
		return
	}

	info.Description, err = a.Description.Get(l)
	return
}

func flagInfo(f Flag, l *i18n.Localizer) (info plugin.FlagInfo, err error) {
	info.Name = f.Name

	if len(f.Aliases) > 0 {
		info.Aliases = make([]string, len(f.Aliases))
		copy(info.Aliases, f.Aliases)
	}

	var ok bool

	info.Type, ok = typeInfo(f.Type, l)
	if !ok {
		return
	}

	info.Multi = f.Multi

	info.Description, err = f.Description.Get(l)
	return
}

func typeInfo(t Type, l *i18n.Localizer) (info plugin.TypeInfo, ok bool) {
	info.Name = t.Name(l)
	if info.Name == "" {
		return
	}

	info.Description = t.Description(l)
	if info.Description == "" {
		return
	}

	ok = true
	return
}

// newArgParsingErr2 creates a new errors.ArgumentParsingError using the passed
// i18n.Config.
// It adds the following additional placeholders: name, used_name, raw and
// position.
// If raw is longer than a 100 characters, it will be shortened.
func newArgParsingErr(
	cfg *i18n.Config, ctx *Context, placeholders map[string]interface{},
) *errors.ArgumentParsingError {
	placeholders = fillPlaceholders(placeholders, ctx)
	return errors.NewArgumentParsingErrorl(cfg.
		WithPlaceholders(placeholders))
}

// newArgParsingErr2 creates a new errors.ArgumentParsingError and decides based
// on the passed Context which of the two i18n.Configs to use.
// It adds the following additional placeholders: name, used_name, raw and
// position.
// If raw is longer than a 100 characters, it will be shortened.
func newArgParsingErr2(
	argConfig, flagConfig *i18n.Config, ctx *Context, placeholders map[string]interface{},
) *errors.ArgumentParsingError {
	placeholders = fillPlaceholders(placeholders, ctx)

	if ctx.Kind == KindArg {
		return errors.NewArgumentParsingErrorl(argConfig.
			WithPlaceholders(placeholders))
	}

	return errors.NewArgumentParsingErrorl(flagConfig.
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

// insertChannelSorted inserts the passed channel into the passed slice channels
// so that the positions of the channels are sorted in ascending order.
func insertChannelSorted(c *discord.Channel, channels []*discord.Channel) []*discord.Channel {
	i := sort.Search(len(channels), func(i int) bool { return channels[i].Position >= c.Position })

	channels = append(channels, new(discord.Channel))
	copy(channels[i+1:], channels[i:])

	channels[i] = c

	return channels
}

// chooserEmojis are the emojis used as options in a chooser embed.
var chooserEmojis = []api.Emoji{emojiutil.Keycap1, emojiutil.Keycap2, emojiutil.Keycap3, emojiutil.Keycap4,
	emojiutil.Keycap5, emojiutil.Keycap6, emojiutil.Keycap7, emojiutil.Keycap8, emojiutil.Keycap9, emojiutil.Keycap10}
