package duration

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type ParseErrorCode string

const (
	// ErrSyntax indicates a syntax error.
	ErrSyntax ParseErrorCode = "invalid duration"
	// ErrSize indicates that the parsed duration exceeds the maximum size of
	// a duration.
	ErrSize ParseErrorCode = "the duration is too large"
	// ErrMissing units indicates that a number is missing its unit.
	ErrMissingUnit ParseErrorCode = "missing unit in duration"
	// ErrInvalidUnit indicates an unknown unit.
	// Val will contain the invalid unit.
	ErrInvalidUnit ParseErrorCode = "invalid unit '%s' in duration"
)

// ParseError is the error returned by Parse, if the passed duration is
// faulty.
type ParseError struct {
	// Code is the error code describing further what error occurred.
	Code ParseErrorCode
	// Val is the faulty part of the raw duration.
	// This will be empty, unless the doc of the error code specifically
	// states it,
	Val string
	// RawDuration is the raw duration, as attempted to parse
	RawDuration string
}

func (p ParseError) Error() string {
	if p.Code == ErrInvalidUnit {
		return "duration: " + fmt.Sprintf(string(ErrInvalidUnit), p.Val) + ": " + p.RawDuration
	}

	return "duration: " + string(p.Code) + ": " + p.RawDuration
}

// Parse parses the passed duration.
// A duration string is a sequence of decimal numbers and units.
// Both numbers and units may be followed by spaces or tabs.
//
// Valid units are "ms", "s", "min", "h", "d", "w", "m", "y" or "yr".
//
// If the passed string does not represent a valid duration, Parse will return
// a *ParseError.
func Parse(s string) (time.Duration, error) {
	return newDurationParser(s).parse()
}

type parser struct {
	raw []rune
	pos int
}

func newDurationParser(s string) *parser {
	s = strings.ToLower(s)

	return &parser{
		raw: []rune(s),
		pos: 0,
	}
}

func (p *parser) parse() (time.Duration, error) {
	p.skipWhitespace()

	var d int64

	for p.has(1) {
		val, err := p.nextDecimal()
		if err != nil {
			return 0, err
		}

		frac, scale, err := p.nextFraction()
		if err != nil {
			return 0, err
		}

		p.skipWhitespace()

		unit, err := p.nextUnit()
		if err != nil {
			return 0, err
		}

		unitMul, ok := units[unit]
		if !ok {
			return 0, &ParseError{
				Code:        ErrInvalidUnit,
				Val:         unit,
				RawDuration: string(p.raw),
			}
		}

		val *= unitMul
		if val < 0 { // overflow
			return 0, &ParseError{
				Code:        ErrSize,
				RawDuration: string(p.raw),
			}
		}

		if frac > 0 {
			val += int64(float64(frac) * (float64(unitMul) / float64(scale)))
			if val < 0 { // overflow
				return 0, &ParseError{
					Code:        ErrSize,
					RawDuration: string(p.raw),
				}
			}
		}

		d += val

		p.skipWhitespace()
	}

	return time.Duration(d), nil
}

// ================================ Helpers ================================

// has checks if there are at least min runes remaining.
func (p *parser) has(min int) bool {
	return p.pos <= len(p.raw)-min
}

func (p *parser) next() rune {
	if !p.has(1) {
		return 0
	}

	p.pos++

	return p.raw[p.pos-1]
}

// backup goes one character back.
func (p *parser) backup() {
	p.pos--
}

// peek peeks numAhead characters ahead, without incrementing the position.
func (p *parser) peek(numAhead int) rune {
	if !p.has(numAhead) {
		return 0
	}

	return p.raw[p.pos+numAhead-1]
}

// skip skips the next num characters.
func (p *parser) skip() {
	if p.has(1) {
		p.pos++
	}
}

// ================================ Parsing ================================

func (p *parser) skipWhitespace() {
	for p.has(1) {
		next := p.next()

		if next != ' ' && next != '\t' {
			p.backup()
			return
		}
	}
}

// nextDecimal consumes the part left of the decimal sign.
func (p *parser) nextDecimal() (val int64, err error) {
	start := p.pos

	for p.has(1) {
		next := p.next()

		if next < '0' || next > '9' {
			p.backup()
			break
		}

		val = val*10 + int64(next-'0')

		if val < 0 { // overflow
			return 0, &ParseError{
				Code:        ErrSize,
				RawDuration: string(p.raw),
			}
		}
	}

	if p.pos == start { // we didn't consume any digits
		return 0, &ParseError{
			Code:        ErrSyntax,
			RawDuration: string(p.raw),
		}
	}

	return
}

// nextFraction consumes the possible decimal point and fraction.
// If there is no fraction, frac will be 0.
func (p *parser) nextFraction() (frac int64, scale int, err error) {
	if p.peek(1) == '.' {
		p.skip()

		scale = 1

		start := p.pos

		var overflow bool

		for p.has(1) {
			next := p.next()

			if next < '0' || next > '9' {
				p.backup()
				break
			}
			if overflow { // ignore remaining digits
				continue
			}

			if frac > int64(math.MaxInt64)/10 {
				overflow = true
				continue
			}

			x := frac*10 + int64(next-'0')

			if frac < 0 { // overflow
				overflow = true
				continue
			}

			frac = x
			scale *= 10
		}

		if p.pos == start { // we didn't consume any digits
			return 0, 0, &ParseError{
				Code:        ErrSyntax,
				RawDuration: string(p.raw),
			}
		}
	}

	return
}

func (p *parser) nextUnit() (unit string, err error) {
	for p.has(1) {
		next := p.next()

		if next < 'a' || next > 'z' {
			p.backup()
			break
		}

		unit += string(next)
	}

	if len(unit) == 0 {
		return "", &ParseError{
			Code:        ErrMissingUnit,
			RawDuration: string(p.raw),
		}
	}

	return
}
