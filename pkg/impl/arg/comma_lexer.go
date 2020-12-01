package arg

import (
	"strings"

	"github.com/mavolin/adam/pkg/errors"
)

type (
	commaLexer struct {
		raw []rune

		start int
		pos   int // next char

		emitChan chan commaItem
		state    commaStateFunc

		// used to keep track of the necessities of minus escapes
		numRequiredArgs int
		nextArg         int
		hasFlags        bool
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

func newCommaLexer(args string, numRequiredArgs int, hasFlags bool) *commaLexer {
	l := &commaLexer{
		raw:             []rune(args),
		emitChan:        make(chan commaItem, 2), // pseudo ring buffer, see nextItem
		numRequiredArgs: numRequiredArgs,
		hasFlags:        hasFlags,
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
				return commaItem{
					typ: itemEOF,
				}, nil
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

// ignore ignores all flagContent up to this point.
// It starts at the upcoming character
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

// ================================ State functions ================================

func (l *commaLexer) item() (commaStateFunc, error) {
	if l.drained() {
		return nil, nil
	}

	// as long as the first required argument hasn't been collected, there are
	// no required args, or all required arg are collected, flags can be placed
	// and minus escapes aren't needed.
	if (l.nextArg == 0 || l.nextArg >= l.numRequiredArgs) && l.hasFlags {
		if l.peek(1) == '-' && l.peek(2) != '-' {
			return l.flag, nil
		}
	}

	return l.arg, nil
}

// flag parses a flag.
// The introducing minus is still present.
func (l *commaLexer) flag() (commaStateFunc, error) {
	l.skip()
	l.ignore()

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
	for l.has(1) { // skip whitespace
		if !strings.ContainsRune(whitespace, l.next()) {
			l.backup()
			break
		}
	}

	l.ignore()

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

	if l.pos > l.start { // make sure we actually collected some content
		l.emit(itemFlagContent)
	}

	return l.end, nil
}

func (l *commaLexer) arg() (commaStateFunc, error) {
	for l.has(1) { // skip whitespace
		if !strings.ContainsRune(whitespace, l.next()) {
			l.backup()
			break
		}
	}

	l.ignore()

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

	if l.pos == l.start { // make sure we actually collected some flagContent
		return nil, errors.NewArgumentParsingErrorl(emptyArgError.
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

	for l.has(1) { // skip whitespace
		if !strings.ContainsRune(whitespace, l.next()) {
			l.backup()
			l.ignore()

			return l.item, nil
		}
	}

	// allow a terminating comma for invokes with at least one arg or flag
	return nil, nil
}
