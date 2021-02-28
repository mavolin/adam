package arg

import (
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// Options is a set of plugin.ArgConfigs with prefixes distinguishing them.
	Options []Option

	// Option is a single option.
	Option struct {
		// Prefix is the prefix the arguments must start with.
		// It will also be used as plugin.Context.ArgCombinationID.
		Prefix string
		// Config is the underlying plugin.ArgConfig.
		Config plugin.ArgConfig
	}
)

var (
	_ plugin.ArgConfig  = Options{}
	_ plugin.ArgsInfoer = Options{}
)

func (o Options) Parse(args string, s *state.State, ctx *plugin.Context) error {
	if len(args) == 0 {
		return plugin.NewArgumentErrorl(notEnoughArgsError)
	}

	prefix := firstWord(args)

	for _, o := range o {
		if o.Prefix == prefix {
			if len(ctx.ArgCombinationID) > 0 {
				ctx.ArgCombinationID += "." + o.Prefix
			} else {
				ctx.ArgCombinationID = o.Prefix
			}

			args := strings.TrimLeft(args[len(prefix):], whitespace)

			if o.Config == nil {
				if len(args) != 0 {
					return plugin.NewArgumentErrorl(tooManyArgsError)
				}

				return nil
			}

			return o.Config.Parse(args, s, ctx)
		}
	}

	return plugin.NewArgumentErrorl(unknownPrefixError.
		WithPlaceholders(unknownPrefixErrorPlaceholders{
			Name: prefix,
		}))
}

// firstWord extracts the first word of the given string.
// A word ends if it is followed by a space, tab or newline.
func firstWord(s string) string {
	for i, char := range s {
		if strings.ContainsRune(whitespace, char) {
			return s[:i]
		}
	}

	return s
}

func (o Options) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	infos := make([]plugin.ArgsInfo, 0, len(o))

	for _, o := range o {
		if o.Config == nil { // special case
			infos = append(infos, plugin.ArgsInfo{Prefix: o.Prefix})
			continue
		}

		infoer, ok := o.Config.(plugin.ArgsInfoer)
		if !ok || infoer == nil {
			return nil
		}

		subInfos := infoer.Info(l)
		if len(subInfos) == 0 {
			return nil
		}

		for i, info := range subInfos {
			if len(info.Prefix) > 0 {
				subInfos[i].Prefix += " " + o.Prefix
			} else {
				subInfos[i].Prefix = o.Prefix
			}
		}

		infos = append(infos, subInfos...)
	}

	return infos
}
