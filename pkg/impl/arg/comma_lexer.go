package arg

import (
	"strings"

	"github.com/mavolin/adam/pkg/plugin"
)

type (
	commaLexer struct {
		raw []rune

		start int
		pos   int // next char

		emitChan chan commaItem
		state    commaStateFunc

		nextArg      int
		parsingFlags bool
	}

	commaStateFunc func() (commaStateFunc, error)

	commaItem struct {
		typ commaItemType
		val string
	}
)

type commaItemType uint8

const (
	itemEOF commaItemType = iota
	itemFlagName
	itemFlagContent
	itemArgContent
	itemComma
)

func (i commaItemType) String() string {
	switch i {
	case itemEOF:
		return "EOF"
	case itemFlagName:
		return "flagName"
	case itemFlagContent:
		return "flagContent"
	case itemArgContent:
		return "argContent"
	case itemComma:
		return "comma"
	default:
		return ""
	}
}

func newCommaLexer(args string) *commaLexer {
	l := &commaLexer{
		raw:          []rune(args),
		emitChan:     make(chan commaItem, 2), // pseudo ring buffer, see nextItem
		parsingFlags: true,
	}

	l.state = l.item
	return l
}

func (l *commaLexer) nextItem() (commaItem, error) {
	for {
		select {
		case emit := <-l.emitChan:
			return emit, nil
		default:
			if l.state == nil {
				return commaItem{typ: itemEOF}, nil
			}

			var err error

			l.state, err = l.state()
			if err != nil {
				return commaItem{}, err
			}
		}
	}
}

// ================================ Helpers ================================

// has checks if there are at least min runes remaining.
func (l *commaLexer) has(min int) bool {
	return l.pos <= len(l.raw)-min
}

func (l *commaLexer) drained() bool {
	return !l.has(1)
}

func (l *commaLexer) next() rune {
	if !l.has(1) {
		return 0
	}

	l.pos++
	return l.raw[l.pos-1]
}

// backup goes one character back.
func (l *commaLexer) backup() {
	l.pos--
}

// peek peeks numAhead characters ahead, without incrementing the position.
func (l *commaLexer) peek(numAhead int) rune {
	if !l.has(numAhead) {
		return 0
	}

	return l.raw[l.pos+numAhead-1]
}

// skip skips the next num characters.
func (l *commaLexer) skip() {
	if l.has(1) {
		l.pos++
	}
}

// ignore ignores all content up to this point.
// It starts at the upcoming character.
func (l *commaLexer) ignore() {
	l.start = l.pos
}

func (l *commaLexer) emit(typ commaItemType) {
	l.emitChan <- commaItem{
		typ: typ,
		val: string(l.raw[l.start:l.pos]),
	}

	l.start = l.pos
}

func (l *commaLexer) ignoreWhitespace() {
	for l.has(1) { // skip whitespace
		if !strings.ContainsRune(whitespace, l.next()) {
			l.backup()
			break
		}
	}

	l.ignore()
}

func (l *commaLexer) consumeContent() {
	for l.has(1) {
		if l.next() == ',' {
			if l.peek(1) == ',' { // escaped comma
				l.skip()
			} else {
				l.backup()
				break
			}
		}
	}
}

// ================================ State functions ================================

func (l *commaLexer) item() (commaStateFunc, error) {
	if l.drained() {
		return nil, nil
	}

	if l.parsingFlags && l.peek(1) == '-' {
		l.skip()
		l.ignore()

		if l.peek(1) != '-' { // mind the l.skip() above
			return l.flag, nil
		}
		// else: this is the first argument, using minus escapes
	}

	l.parsingFlags = false
	return l.arg, nil
}

// flag parses a flag.
func (l *commaLexer) flag() (commaStateFunc, error) {
	for l.has(1) {
		next := l.next()

		if strings.ContainsRune(whitespace, next) {
			l.backup()
			l.emit(itemFlagName)

			return l.flagContent, nil
		}
	}

	l.emit(itemFlagName)
	return nil, nil
}

func (l *commaLexer) flagContent() (commaStateFunc, error) {
	l.ignoreWhitespace()
	l.consumeContent()

	if l.pos > l.start { // make sure we actually collected some content
		l.emit(itemFlagContent)
	}

	return l.end, nil
}

func (l *commaLexer) arg() (commaStateFunc, error) {
	l.consumeContent()

	if l.pos == l.start { // make sure we actually collected some flagContent
		return nil, plugin.NewArgumentErrorl(emptyArgError.
			WithPlaceholders(emptyArgErrorPlaceholders{
				Position: l.nextArg + 1,
			}))
	}

	l.emit(itemArgContent)
	l.nextArg++

	return l.end, nil
}

// end is called if the end of an argument or a flag is reached.
// It expects either EOF or a comma.
func (l *commaLexer) end() (commaStateFunc, error) {
	if l.drained() {
		return nil, nil
	}

	l.skip()
	l.emit(itemComma)

	l.ignoreWhitespace()
	return l.item, nil
}
