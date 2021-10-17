package arg //nolint:dupl

import (
	"strings"

	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// DelimiterParser is a plugin.ArgParser, that uses a custom delimiter to
// separate flags and arguments.
// Literal delimiters can be escaped by using the delimiter twice in a row.
//
// Flags may be placed in front of the arguments.
//
// If the first argument starts with a minus, the minus must be escaped
// through a double minus to avoid confusion with a flag.
//
// Examples
//
// The below examples use ',' as the delimiter.
//
// 	cmd -flag1 abc, -flag2, first arg, second arg,, with a comma in it
//	cmd --first arg using a minus escape
type DelimiterParser struct {
	// Delimiter is the delimiter used by the parser.
	// It can be escaped by using the delimiter twice in a row, e.g. ",," if
	// using ',' as a delimiter.
	//
	// Minus ('-'), Space (' ') and New-Lines ('\n') are not permitted and may
	// lead to unexpected behavior.
	//
	// Default: ','
	Delimiter rune
}

var _ plugin.ArgParser = new(DelimiterParser)

func (p *DelimiterParser) Parse(args string, argConfig plugin.ArgConfig, s *state.State, ctx *plugin.Context) error {
	if p.Delimiter == 0 {
		p.Delimiter = ','
	}

	return newDelimiterParser(args, argConfig, p.Delimiter, s, ctx).parse()
}

// FormatArgs formats the arguments.
// It uses the parser's Delimiter followed by a space to separate arguments and
// flags.
func (p *DelimiterParser) FormatArgs(_ plugin.ArgConfig, args []string, flags map[string]string) string {
	// add double-minus escape to the first arg, if it starts with a minus
	if len(args) > 0 && strings.HasPrefix(args[0], "-") {
		args[0] = "-" + args[0]
	}

	var n int

	for i, arg := range args {
		arg = strings.ReplaceAll(arg, string(p.Delimiter), string(p.Delimiter)+string(p.Delimiter))
		// n += len(arg+p.Delimiter+" ")
		n += len(arg) + 1 + len(" ")

		args[i] = arg
	}

	for name, val := range flags {
		val = strings.ReplaceAll(val, string(p.Delimiter), string(p.Delimiter)+string(p.Delimiter))

		// n += len("-"+name+" "+val+p.Delimiter+" ")
		n += len("-") + len(name) + len(" ") + len(val) + 1 + len(" ")

		flags[name] = val
	}

	// remove the trailing delimiter
	n -= 1 + len(" ")
	if n <= 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(n)

	for name, val := range flags {
		if b.Len() > 0 {
			b.WriteRune(p.Delimiter)
			b.WriteRune(' ')
		}

		b.WriteRune('-')
		b.WriteString(name)
		b.WriteRune(' ')
		b.WriteString(val)
	}

	for _, arg := range args {
		if b.Len() > 0 {
			b.WriteRune(p.Delimiter)
			b.WriteRune(' ')
		}

		b.WriteString(arg)
	}

	return b.String()
}

func (p *DelimiterParser) FormatUsage(_ plugin.ArgConfig, args []string) string {
	if len(args) == 0 {
		return ""
	}

	// we need to use the separator (p.Delimiter+" ") (len(args)-1) times
	n := (len(args) - 1) * 2

	for _, arg := range args {
		n += len(arg)
	}

	var b strings.Builder
	b.Grow(n)

	for _, arg := range args {
		if b.Len() > 0 {
			b.WriteRune(p.Delimiter)
			b.WriteRune(' ')
		}

		b.WriteString(arg)
	}

	return b.String()
}

func (p *DelimiterParser) FormatFlag(name string) string {
	return "-" + name
}

// =============================================================================
// Parsing Logic
// =====================================================================================

type delimiterParser struct {
	helper    *parseHelper
	lexer     *delimiterLexer
	delimiter rune
}

func newDelimiterParser(
	args string, argConfig plugin.ArgConfig, delim rune, s *state.State, ctx *plugin.Context,
) *delimiterParser {
	return &delimiterParser{
		helper: newParseHelper(argConfig.GetRequiredArgs(), argConfig.GetOptionalArgs(),
			argConfig.GetFlags(), argConfig.IsVariadic(), s, ctx),
		lexer:     newCommaLexer(args, delim),
		delimiter: delim,
	}
}

func (p *delimiterParser) parse() error {
	if err := p.startParse(); err != nil {
		return err
	}

	return p.helper.store()
}

func (p *delimiterParser) startParse() error {
	item, err := p.lexer.nextItem()
	if err != nil {
		return err
	}

	if len(p.helper.rargData)+len(p.helper.oargData)+len(p.helper.flagData) == 0 && item.typ != itemEOF {
		return plugin.NewArgumentErrorl(noArgsError)
	}

	for ; err == nil && item.typ != itemEOF; item, err = p.lexer.nextItem() {
		// the lexer keeps track of the correct ordering, so we don't need to
		// worry about that
		switch item.typ { //nolint:exhaustive
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
// The passed delimiterItem must have the delimiterItemType itemFlagName.
func (p *delimiterParser) parseFlag(flagName delimiterItem) (err error) {
	f := p.helper.flag(flagName.val)
	if f == nil {
		return plugin.NewArgumentErrorl(unknownFlagError.
			WithPlaceholders(unknownFlagErrorPlaceholders{
				Name: flagName.val,
			}))
	}

	if f.GetType() == Switch {
		if err = p.helper.addFlag(f, "", ""); err != nil {
			return err
		}
	} else {
		content, err := p.lexer.nextItem()
		if err != nil {
			return err
		} else if content.typ != itemFlagContent {
			return plugin.NewArgumentErrorl(emptyFlagError.
				WithPlaceholders(emptyFlagErrorPlaceholders{
					Name: flagName.val,
				}))
		}

		contentString := strings.ReplaceAll(content.val, string(p.delimiter)+string(p.delimiter), string(p.delimiter))

		if err = p.helper.addFlag(f, flagName.val, contentString); err != nil {
			return err
		}
	}

	finalizer, err := p.lexer.nextItem()
	switch {
	case err != nil:
		return err
	case finalizer.typ == itemFlagContent && f.GetType() == Switch:
		return plugin.NewArgumentErrorl(switchWithContentError.
			WithPlaceholders(&switchWithContentErrorPlaceholders{
				Name: flagName.val,
			}))
	case finalizer.typ != itemDelimiter && finalizer.typ != itemEOF:
		return errors.NewWithStackf("arg: unexpected item during parsing: %s", finalizer.typ)
	default:
		return nil
	}
}

func (p *delimiterParser) parseArg(content delimiterItem) error {
	if strings.HasPrefix(content.val, "--") {
		content.val = content.val[1:]
	}

	content.val = strings.ReplaceAll(content.val, string(p.delimiter)+string(p.delimiter), string(p.delimiter))

	err := p.helper.addArg(content.val)
	if err != nil {
		return err
	}

	finalizer, err := p.lexer.nextItem()
	if err != nil {
		return err
	} else if finalizer.typ != itemDelimiter && finalizer.typ != itemEOF {
		return errors.NewWithStackf("arg: unexpected item during parsing: %s", finalizer.typ)
	}

	return nil
}
