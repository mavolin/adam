package arg

import (
	"strings"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/plugin"
)

type (
	delimiterLexer struct {
		raw []rune

		start int
		pos   int // next char

		emitChan chan delimiterItem
		state    delimiterStateFunc

		nextArg      int
		parsingFlags bool

		delimiter rune
	}

	delimiterStateFunc func() (delimiterStateFunc, error)

	delimiterItem struct {
		typ delimiterItemType
		val string
	}
)

type delimiterItemType uint8

const (
	itemEOF delimiterItemType = iota
	itemFlagName
	itemFlagContent
	itemArgContent
	itemDelimiter
)

func (i delimiterItemType) String() string {
	switch i {
	case itemEOF:
		return "EOF"
	case itemFlagName:
		return "flagName"
	case itemFlagContent:
		return "flagContent"
	case itemArgContent:
		return "argContent"
	case itemDelimiter:
		return "comma"
	default:
		return ""
	}
}

func newCommaLexer(args string, delim rune) *delimiterLexer {
	l := &delimiterLexer{
		raw:          []rune(args),
		emitChan:     make(chan delimiterItem, 1), // pseudo ring buffer, see nextItem
		parsingFlags: true,
		delimiter:    delim,
	}

	l.state = l.item
	return l
}

func (l *delimiterLexer) nextItem() (delimiterItem, error) {
	for {
		select {
		case emit := <-l.emitChan:
			return emit, nil
		default:
			if l.state == nil {
				return delimiterItem{typ: itemEOF}, nil
			}

			var err error

			l.state, err = l.state()
			if err != nil {
				return delimiterItem{}, err
			}
		}
	}
}

// ================================ Helpers ================================

// has checks if there are at least min runes remaining.
func (l *delimiterLexer) has(min int) bool {
	return l.pos <= len(l.raw)-min
}

func (l *delimiterLexer) drained() bool {
	return !l.has(1)
}

func (l *delimiterLexer) next() rune {
	if !l.has(1) {
		return 0
	}

	l.pos++
	return l.raw[l.pos-1]
}

// backup goes one character back.
func (l *delimiterLexer) backup() {
	l.pos--
}

// peek peeks numAhead characters ahead, without incrementing the position.
func (l *delimiterLexer) peek(numAhead int) rune {
	if !l.has(numAhead) {
		return 0
	}

	return l.raw[l.pos+numAhead-1]
}

// skip skips the next character.
func (l *delimiterLexer) skip() {
	if l.has(1) {
		l.pos++
	}
}

// ignore ignores all content up to this point.
// It starts at the upcoming character.
func (l *delimiterLexer) ignore() {
	l.start = l.pos
}

func (l *delimiterLexer) emit(typ delimiterItemType) {
	l.emitChan <- delimiterItem{
		typ: typ,
		val: string(l.raw[l.start:l.pos]),
	}

	l.start = l.pos
}

func (l *delimiterLexer) ignoreWhitespace() {
	for l.has(1) {
		if !strings.ContainsRune(shared.Whitespace, l.next()) {
			l.backup()
			break
		}
	}

	l.ignore()
}

func (l *delimiterLexer) consumeContent() {
	for l.has(1) {
		if l.next() == l.delimiter {
			if l.peek(1) == l.delimiter { // escaped delimiter
				l.skip()
			} else {
				l.backup()
				break
			}
		}
	}
}

// ================================ State functions ================================

func (l *delimiterLexer) item() (delimiterStateFunc, error) {
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
func (l *delimiterLexer) flag() (delimiterStateFunc, error) {
	for l.has(1) {
		next := l.next()

		if strings.ContainsRune(shared.Whitespace, next) {
			l.backup()
			l.emit(itemFlagName)

			return l.flagContent, nil
		}
	}

	l.emit(itemFlagName)
	return nil, nil
}

func (l *delimiterLexer) flagContent() (delimiterStateFunc, error) {
	l.ignoreWhitespace()
	l.consumeContent()

	if l.pos > l.start { // make sure we actually collected some content
		l.emit(itemFlagContent)
	}

	return l.end, nil
}

func (l *delimiterLexer) arg() (delimiterStateFunc, error) {
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
func (l *delimiterLexer) end() (delimiterStateFunc, error) {
	if l.drained() {
		return nil, nil
	}

	l.skip()
	l.emit(itemDelimiter)

	l.ignoreWhitespace()
	return l.item, nil
}
