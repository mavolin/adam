// Package capbuilder provides the CappedBuilder.
package capbuilder

import "strings"

type CappedBuilder struct {
	used, cap int
	b         *strings.Builder
}

// New creates a new CappedBuilder with the passed total capacity and the
// passed initial chunk capacity.
func New(totalCap, chunkCap int) *CappedBuilder {
	if totalCap < chunkCap {
		chunkCap = totalCap
	}

	b := new(strings.Builder)
	b.Grow(chunkCap)

	return &CappedBuilder{cap: totalCap, b: b}
}

// WriteRune writes the passed rune to the CappedBuilder's current chunk, if
// there is space.
func (b *CappedBuilder) WriteRune(r rune) {
	if b.used < b.cap && b.b.Len() < b.b.Cap() {
		b.b.WriteRune(r)
		b.used++
	}
}

// WriteString writes th passed string to the CappedBuilder's current chunk.
// If there is insufficient space for the entire string, only part of it will
// be written to the chunk.
// If there is no space, the string will be discarded.
func (b *CappedBuilder) WriteString(s string) {
	if b.used+len(s) < b.cap && b.b.Len()+len(s) < b.b.Cap() {
		b.b.WriteString(s)
		b.used += len(s)
	} else if b.used <= b.cap || b.b.Len() <= b.b.Cap() {
		end := b.cap - b.used
		if b.b.Cap()-b.b.Len() < end {
			end = b.b.Cap() - b.b.Len()
		}

		b.b.WriteString(s[:end])
		b.used += end
	}
}

// Use manually uses up n characters of the total chunk, without writing
// anything to the current chunk, or draining it.
func (b *CappedBuilder) Use(n int) {
	b.used += n
}

// String returns the string value of the current chunk.
func (b *CappedBuilder) String() string {
	return b.b.String()
}

// Reset resets the current chunk and creates a new one with the passed
// capacity.
func (b *CappedBuilder) Reset(chunkCap int) {
	b.b.Reset()
	b.b.Grow(chunkCap)
}

// ChunkLen returns the length of the current chunk.
func (b *CappedBuilder) ChunkLen() int {
	return b.b.Len()
}

// Rem returns the total amount of remaining characters.
func (b *CappedBuilder) Rem() int {
	rem := b.cap - b.used
	if chunkRem := b.b.Cap() - b.b.Len(); chunkRem < rem {
		rem = chunkRem
	}

	return rem
}
