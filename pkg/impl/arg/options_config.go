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
		Prefix string
		// Config is the underlying plugin.ArgConfig.
		Config plugin.ArgConfig
	}
)

var (
	_ plugin.ArgConfig  = Options{}
	_ plugin.ArgsInfoer = Options{}
)

func (o Options) Parse(args string, s *state.State, ctx *plugin.Context) (plugin.Args, plugin.Flags, error) {
	if len(args) == 0 {
		return nil, nil, plugin.NewArgumentErrorl(notEnoughArgsError)
	}

	prefix := firstWord(args)

	for _, o := range o {
		if o.Prefix == prefix {
			args := strings.TrimLeft(args[len(prefix):], whitespace)

			if o.Config == nil {
				if len(args) != 0 {
					return nil, nil, plugin.NewArgumentErrorl(tooManyArgsError)
				}

				return nil, nil, nil
			}

			return o.Config.Parse(args, s, ctx)
		}
	}

	return nil, nil, plugin.NewArgumentErrorl(unknownPrefixError.
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

		info := infoer.Info(l)
		if len(info) != 1 {
			return nil
		}

		info[0].Prefix = o.Prefix

		infos = append(infos, info[0])
	}

	return infos
}
