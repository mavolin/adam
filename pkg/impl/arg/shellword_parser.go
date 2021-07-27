package arg //nolint:dupl

import (
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/plugin"
)

// ShellwordParser is a plugin.ArgConfig that roughly follows the parsing rules
// of the Bourne shell.
//
// Flags
//
// Flags can be placed both before and after arguments.
// For simplicity, flags always start with a single minus, double minuses are
// not permitted.
//
// Arguments
//
// Arguments are space separated.
// To use arguments with whitespace quotes, both single and double, can be
// used.
// Additionally, lines of code as well as code blocks will be parsed as a
// single argument.
//
// Escapes
//
// Escapes are only permitted if using double quotes.
// Valid escapes are '\\' and '\"', all other combinations will be parsed
// literally to make usage easier for users unaware of escapes.
var ShellwordParser plugin.ArgParser = new(shellwordParser)

type shellwordParser struct{}

func (p *shellwordParser) Parse(args string, argConfig plugin.ArgConfig, s *state.State, ctx *plugin.Context) error {
	return newShellwordParserState(args, argConfig, s, ctx).parse()
}

var shellwordEscapeReplacer = strings.NewReplacer(`"`, `\"`, `\`, `\\`)

func (p *shellwordParser) FormatArgs(_ plugin.ArgConfig, args []string, flags map[string]string) string {
	if len(args) > 0 && strings.HasPrefix(args[0], "-") {
		args[0] = "-" + args[0]
	}

	var n int

	for i, arg := range args {
		if strings.Contains(arg, " ") || strings.HasPrefix(arg, "-") {
			arg = shellwordEscapeReplacer.Replace(arg)
			arg = `"` + arg + `"`
		}

		n += len(arg) + len(" ")
		args[i] = arg
	}

	for name, val := range flags {
		if strings.Contains(val, " ") || strings.HasPrefix(val, "-") {
			val = shellwordEscapeReplacer.Replace(val)
			val = `"` + val + `"`
		}

		n += len("-") + len(name) + len(" ") + len(val) + len(" ")
		flags[name] = val
	}

	// remove the trailing delimiter
	n -= len(" ")

	var b strings.Builder
	b.Grow(n)

	for name, val := range flags {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}

		b.WriteRune('-')
		b.WriteString(name)
		b.WriteRune(' ')
		b.WriteString(val)
	}

	for _, arg := range args {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}

		b.WriteString(arg)
	}

	return b.String()
}

func (p *shellwordParser) FormatUsage(_ plugin.ArgConfig, args []string) string {
	// we need (len(args)-1) space-separators
	n := len(args) - 1

	for _, arg := range args {
		n += len(arg)
	}

	var b strings.Builder
	b.Grow(n)

	for _, arg := range args {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}

		b.WriteString(arg)
	}

	return b.String()
}

func (p *shellwordParser) FormatFlag(name string) string {
	return "-" + name
}

// =============================================================================
// Parsing Logic
// =====================================================================================

type groupingCharacter uint8

const (
	singleQuote groupingCharacter = iota + 1
	doubleQuote
	singleBacktick
	doubleBacktick
	tripleBacktick
)

func (c groupingCharacter) String() string {
	switch c {
	case singleQuote:
		return "'"
	case doubleQuote:
		return `"`
	case singleBacktick:
		return "\\`"
	case doubleBacktick:
		return "\\`\\`"
	case tripleBacktick:
		return "\\`\\`\\`"
	default:
		return ""
	}
}

type shellwordParserState struct {
	helper *parseHelper

	raw []rune
	pos int

	builder strings.Builder
}

func newShellwordParserState(
	args string, argConfig plugin.ArgConfig, s *state.State, ctx *plugin.Context,
) *shellwordParserState {
	p := &shellwordParserState{
		helper: newParseHelper(argConfig.GetRequiredArgs(), argConfig.GetOptionalArgs(), argConfig.GetFlags(),
			argConfig.IsVariadic(), s, ctx),
		raw: []rune(args),
	}

	p.builder.Grow(len(args))

	return p
}

func (p *shellwordParserState) parse() error {
	if len(p.helper.rargData)+len(p.helper.oargData)+len(p.helper.flagData) == 0 && len(p.raw) != 0 {
		return plugin.NewArgumentErrorl(noArgsError)
	}

	if err := p.parseFlags(); err != nil {
		return err
	}

	if err := p.parseArgs(); err != nil {
		return err
	}

	if err := p.parseFlags(); err != nil {
		return err
	}

	return p.helper.store()
}

// has checks if there are at least min runes remaining.
func (p *shellwordParserState) has(min int) bool {
	return p.pos <= len(p.raw)-min
}

func (p *shellwordParserState) drained() bool {
	return !p.has(1)
}

func (p *shellwordParserState) next() rune {
	if !p.has(1) {
		return 0
	}

	p.pos++

	return p.raw[p.pos-1]
}

// backup goes one character back.
func (p *shellwordParserState) backup() {
	p.pos--
}

// peek peeks numAhead characters ahead, without incrementing the position.
func (p *shellwordParserState) peek(numAhead int) rune {
	if !p.has(numAhead) {
		return 0
	}

	return p.raw[p.pos+numAhead-1]
}

// skip skips the next num characters.
func (p *shellwordParserState) skip(num int) {
	if p.has(num) {
		p.pos += num
	}
}

func (p *shellwordParserState) skipWhitespace() {
	for p.has(1) { // skip whitespace
		if !strings.ContainsRune(shared.Whitespace, p.next()) {
			p.backup()
			break
		}
	}
}

//nolint:gocognit
func (p *shellwordParserState) nextContent() (string, error) {
	var (
		gc       groupingCharacter
		upEscape bool
	)

	defer p.builder.Reset()

	switch p.peek(1) {
	case '\'', '‘', '’', '‚', '‛':
		gc = singleQuote
		p.skip(1)
	case '"', '“', '”', '„', '‟':
		gc = doubleQuote
		p.skip(1)
	case '`':
		if p.peek(2) == '`' {
			if p.peek(3) == '`' {
				gc = tripleBacktick
				p.skip(3)
			} else {
				gc = doubleBacktick
				p.skip(2)
			}
		} else {
			gc = singleBacktick
			p.skip(1)
		}
	}

	for char := p.next(); char != 0; char = p.next() {
		switch {
		case upEscape:
			if char != '\\' && strings.ContainsRune(`“”„‟`, char) {
				p.builder.WriteRune('"')
			}

			p.builder.WriteRune(char)

			upEscape = false
			continue
		case char == '\\' && gc == doubleQuote:
			upEscape = true
			continue
		case strings.ContainsRune(shared.Whitespace, char) && gc == 0:
			p.backup()
			return p.builder.String(), nil
		case char == '\'' && gc == singleQuote:
			return p.builder.String(), nil
		case char == '"' && gc == doubleQuote:
			return p.builder.String(), nil
		case char == '`':
			if p.peek(1) == '`' {
				if p.peek(2) == '`' {
					if gc == tripleBacktick {
						p.skip(2)
						return p.builder.String(), nil
					}
				} else {
					if gc == doubleBacktick {
						p.skip(1)
						return p.builder.String(), nil
					}
				}
			} else if gc == singleBacktick {
				return p.builder.String(), nil
			}
		}

		p.builder.WriteRune(char)
	}

	if gc != 0 {
		return "", plugin.NewArgumentErrorl(groupNotClosedError.
			WithPlaceholders(groupNotClosedErrorPlaceholders{
				Quote: gc.String(),
			}))
	}

	return p.builder.String(), nil
}

func (p *shellwordParserState) parseFlags() error {
	for {
		p.skipWhitespace()
		if p.drained() {
			return nil
		}

		if p.peek(1) != '-' {
			return nil
		}

		p.skip(1)

		start := p.pos

		for char := p.next(); !p.drained(); char = p.next() {
			if strings.ContainsRune(shared.Whitespace, char) {
				p.backup()
				break
			}
		}

		if p.pos == start { // interpret the minus literally
			p.backup()
			return nil
		}

		name := string(p.raw[start:p.pos])

		f := p.helper.flag(name)
		if f == nil {
			return plugin.NewArgumentErrorl(unknownFlagError.
				WithPlaceholders(unknownFlagErrorPlaceholders{
					Name: name,
				}))
		}

		if f.GetType() == Switch {
			if err := p.helper.addFlag(f, "", ""); err != nil {
				return err
			}

			continue
		}

		p.skipWhitespace()
		if p.drained() {
			return nil
		}

		content, err := p.nextContent()
		if err != nil {
			return err
		}

		if err = p.helper.addFlag(f, name, content); err != nil {
			return err
		}
	}
}

func (p *shellwordParserState) parseArgs() error {
	for {
		p.skipWhitespace()
		if p.drained() {
			return nil
		}

		if p.peek(1) == '-' {
			return nil
		}

		content, err := p.nextContent()
		if err != nil {
			return err
		}

		if err = p.helper.addArg(content); err != nil {
			return err
		}
	}
}
