package arg

import (
	"reflect"
	"strings"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

var commaConfigEscaper = strings.NewReplacer(",,", ",", "--", "-")

var interfaceType = reflect.TypeOf(func(interface{}) {}).In(0)

type commaParser struct {
	helper *parseHelper
	lexer  *commaLexer
}

func newCommaParser(args string, cfg CommaConfig, s *state.State, ctx *plugin.Context) *commaParser {
	return &commaParser{
		helper: newParseHelper(cfg.RequiredArgs, cfg.OptionalArgs, cfg.Flags, cfg.Variadic, s, ctx),
		lexer:  newCommaLexer(args, len(cfg.RequiredArgs), len(cfg.Flags) > 0),
	}
}

func (p *commaParser) parse() (plugin.Args, plugin.Flags, error) {
	err := p.startParse()
	if err != nil {
		return nil, nil, err
	}

	return p.helper.get()
}

func (p *commaParser) startParse() error {
	item, err := p.lexer.nextItem()
	if err != nil {
		return err
	}

	if len(p.helper.rargData)+len(p.helper.oargData)+len(p.helper.flagData) == 0 && item.typ != itemEOF {
		return errors.NewArgumentParsingErrorl(noArgsError)
	}

	for ; err == nil && item.typ != itemEOF; item, err = p.lexer.nextItem() {
		// the lexer keeps track of the correct ordering, so we don't need to
		// worry about that
		switch item.typ {
		case itemFlagName:
			err = p.parseFlag(item)
		case itemArgContent:
			err = p.parseArg(item)
		default:
			return errors.NewWithStackf("arg: unexpected item during parsing: %s", item.typ)
		}

		if err != nil {
			return err
		}
	}

	return err
}

// parseFlag parses a flag.
// The passed commaItem must have the commaItemType itemFlagName.
func (p *commaParser) parseFlag(flagName commaItem) (err error) {
	f := p.helper.flag(flagName.val)
	if f == nil {
		return errors.NewArgumentParsingErrorl(unknownFlagError.
			WithPlaceholders(unknownFlagErrorPlaceholders{
				Name: flagName.val,
			}))
	}

	if f.Type == Switch {
		if err = p.helper.addFlag(f, "", ""); err != nil {
			return err
		}
	} else {
		content, err := p.lexer.nextItem()
		if err != nil {
			return err
		} else if content.typ != itemFlagContent {
			return errors.NewArgumentParsingErrorl(emptyFlagError.
				WithPlaceholders(emptyFlagErrorPlaceholders{
					Name: flagName.val,
				}))
		}

		contentString := commaConfigEscaper.Replace(content.val)

		if err = p.helper.addFlag(f, flagName.val, contentString); err != nil {
			return err
		}
	}

	finalizer, err := p.lexer.nextItem()
	switch {
	case err != nil:
		return err
	case finalizer.typ == itemFlagContent && f.Type == Switch:
		return errors.NewArgumentParsingErrorl(switchWithContentError.
			WithPlaceholders(&switchWithContentErrorPlaceholders{
				Name: flagName.val,
			}))
	case finalizer.typ != itemComma && finalizer.typ != itemEOF:
		return errors.NewWithStackf("arg: unexpected item during parsing: %s", finalizer.typ)
	default:
		return nil
	}
}

func (p *commaParser) parseArg(content commaItem) error {
	err := p.helper.addArg(commaConfigEscaper.Replace(content.val))
	if err != nil {
		return err
	}

	finalizer, err := p.lexer.nextItem()
	if err != nil {
		return err
	} else if finalizer.typ != itemComma && finalizer.typ != itemEOF {
		return errors.NewWithStackf("arg: unexpected item during parsing: %s", finalizer.typ)
	}

	return nil
}
