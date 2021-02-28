package arg

import (
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

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

type shellwordParser struct {
	helper *parseHelper

	raw []rune
	pos int

	builder strings.Builder
}

func newShellwordParser(args string, cfg ShellwordConfig, s *state.State, ctx *plugin.Context) *shellwordParser {
	p := &shellwordParser{
		helper: newParseHelper(cfg.Required, cfg.Optional, cfg.Flags, cfg.Variadic, s, ctx),
		raw:    []rune(args),
	}

	p.builder.Grow(len(args))

	return p
}

func newShellwordParserl(
	args string, cfg LocalizedShellwordConfig, s *state.State, ctx *plugin.Context,
) *shellwordParser {
	p := &shellwordParser{
		helper: newParseHelperl(cfg.Required, cfg.Optional, cfg.Flags, cfg.Variadic, s, ctx),
		raw:    []rune(args),
	}

	p.builder.Grow(len(args))

	return p
}

func (p *shellwordParser) parse() error {
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

	return p.helper.putContext()
}

// has checks if there are at least min runes remaining.
func (p *shellwordParser) has(min int) bool {
	return p.pos <= len(p.raw)-min
}

func (p *shellwordParser) drained() bool {
	return !p.has(1)
}

func (p *shellwordParser) next() rune {
	if !p.has(1) {
		return 0
	}

	p.pos++

	return p.raw[p.pos-1]
}

// backup goes one character back.
func (p *shellwordParser) backup() {
	p.pos--
}

// peek peeks numAhead characters ahead, without incrementing the position.
func (p *shellwordParser) peek(numAhead int) rune {
	if !p.has(numAhead) {
		return 0
	}

	return p.raw[p.pos+numAhead-1]
}

// skip skips the next num characters.
func (p *shellwordParser) skip(num int) {
	if p.has(num) {
		p.pos += num
	}
}

func (p *shellwordParser) skipWhitespace() {
	for p.has(1) { // skip whitespace
		if !strings.ContainsRune(whitespace, p.next()) {
			p.backup()
			break
		}
	}
}

func (p *shellwordParser) nextContent() (string, error) { //nolint:gocognit
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
		case strings.ContainsRune(whitespace, char) && gc == 0:
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

func (p *shellwordParser) parseFlags() error {
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
			if strings.ContainsRune(whitespace, char) {
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

		if f.typ == Switch {
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

func (p *shellwordParser) parseArgs() error {
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
