package capbuilder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCappedBuilderWriteRune(t *testing.T) {
	t.Parallel()

	t.Run("chunk limit", func(t *testing.T) {
		t.Parallel()

		totalCap := 20
		chunkCap := 5

		b := New(totalCap, chunkCap)

		for i := 0; i < totalCap; i++ {
			b.WriteRune('a')
		}

		expect := strings.Repeat("a", chunkCap)

		assert.Equal(t, b.used, chunkCap)
		assert.Equal(t, expect, b.String())
	})

	t.Run("global limit", func(t *testing.T) {
		t.Parallel()

		totalCap := 10
		chunkCap := 6

		b := New(totalCap, chunkCap)

		for i := 0; i < chunkCap; i++ {
			b.WriteRune('a')
		}

		b.Reset(chunkCap)

		for i := 0; i < chunkCap; i++ {
			b.WriteRune('a')
		}

		expect := strings.Repeat("a", totalCap-chunkCap)

		assert.Equal(t, b.used, b.cap)
		assert.Equal(t, expect, b.String())
	})
}

func TestCappedBuilder_WriteString(t *testing.T) {
	t.Parallel()

	t.Run("chunk limit", func(t *testing.T) {
		t.Parallel()

		t.Run("full Write", func(t *testing.T) {
			t.Parallel()

			totalCap := 20
			chunkCap := 6

			b := New(totalCap, chunkCap)

			for i := 0; i < totalCap; i++ {
				b.WriteString("ab")
			}

			expect := strings.Repeat("ab", chunkCap/2)

			assert.Equal(t, b.used, chunkCap)
			assert.Equal(t, expect, b.String())
		})

		t.Run("partial Write", func(t *testing.T) {
			t.Parallel()

			totalCap := 20
			chunkCap := 5

			b := New(totalCap, chunkCap)

			for i := 0; i < totalCap; i++ {
				b.WriteString("ab")
			}

			expect := "ababa"

			assert.Equal(t, b.used, chunkCap)
			assert.Equal(t, expect, b.String())
		})
	})

	t.Run("total limit", func(t *testing.T) {
		t.Parallel()

		t.Run("full write", func(t *testing.T) {
			t.Parallel()

			totalCap := 10
			chunkCap := 6

			b := New(totalCap, chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.WriteString("ab")
			}

			b.Reset(chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.WriteString("ab")
			}

			expect := strings.Repeat("ab", (totalCap-chunkCap)/2)

			assert.Equal(t, b.used, b.cap)
			assert.Equal(t, expect, b.String())
		})

		t.Run("partial write", func(t *testing.T) {
			t.Parallel()

			totalCap := 9
			chunkCap := 6

			b := New(totalCap, chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.WriteString("ab")
			}

			b.Reset(chunkCap)

			for i := 0; i < chunkCap/2; i++ {
				b.WriteString("ab")
			}

			expect := "aba"

			assert.Equal(t, b.used, b.cap)
			assert.Equal(t, expect, b.String())
		})
	})
}

func TestCappedBuilder_Reset(t *testing.T) {
	t.Parallel()

	chunkCap := 3

	b := New(10, chunkCap)
	assert.Equal(t, chunkCap, b.b.Cap())
	assert.Equal(t, 0, b.b.Len())

	b.WriteRune('a')
	assert.Equal(t, 1, b.b.Len())

	b.Reset(chunkCap + 1)
	assert.Equal(t, chunkCap+1, b.b.Cap())
	assert.Equal(t, 0, b.b.Len())

	b.WriteRune('a')
	assert.Equal(t, 1, b.b.Len())
}

func TestCappedBuilderRem(t *testing.T) {
	t.Parallel()

	t.Run("chunk", func(t *testing.T) {
		t.Parallel()

		chunkCap := 5

		b := New(10, chunkCap)
		b.WriteRune('a')

		assert.Equal(t, chunkCap-1, b.Rem())
	})

	t.Run("total", func(t *testing.T) {
		t.Parallel()

		totalCap := 10
		chunkCap := 7

		b := New(totalCap, chunkCap)
		b.WriteString(strings.Repeat("a", chunkCap))
		b.Reset(chunkCap)

		assert.Equal(t, totalCap-chunkCap, b.Rem())
	})
}
