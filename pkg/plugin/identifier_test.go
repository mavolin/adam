package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier_Parent(t *testing.T) {
	testCases := []struct {
		name       string
		identifier Identifier
		expect     []Identifier
	}{
		{
			name:       "root",
			identifier: ".",
			expect:     []Identifier{"."},
		},
		{
			name:       "single level",
			identifier: ".mod",
			expect:     []Identifier{"."},
		},
		{
			name:       "multi level",
			identifier: ".mod.infr.edit",
			expect:     []Identifier{".mod.infr", ".mod", "."},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.identifier

			for lvl, expect := range c.expect {
				actual = actual.Parent()
				assert.Equal(t, expect, actual, "unexpected identifier returned on level %d", lvl+1)
			}
		})
	}
}

func TestIdentifier_IsRoot(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		isRoot := Identifier(".").IsRoot()
		assert.True(t, isRoot)
	})

	t.Run("not root", func(t *testing.T) {
		isRoot := Identifier(".mod").IsRoot()
		assert.False(t, isRoot)
	})
}

func TestIdentifier_IsParent(t *testing.T) {
	testCases := []struct {
		name       string
		identifier Identifier
		target     Identifier
		expect     bool
	}{
		{
			name:       "parent",
			identifier: ".mod.ban",
			target:     ".mod",
			expect:     true,
		},
		{
			name:       "equal",
			identifier: ".mod.ban",
			target:     ".mod.ban",
			expect:     false,
		},
		{
			name:       "child",
			identifier: ".mod",
			target:     ".mod.ban",
			expect:     false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.identifier.IsParent(c.target)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestIdentifier_IsChild(t *testing.T) {
	testCases := []struct {
		name       string
		identifier Identifier
		target     Identifier
		expect     bool
	}{
		{
			name:       "parent",
			identifier: ".mod.ban",
			target:     ".mod",
			expect:     false,
		},
		{
			name:       "equal",
			identifier: ".mod.ban",
			target:     ".mod.ban",
			expect:     false,
		},
		{
			name:       "child",
			identifier: ".mod",
			target:     ".mod.ban",
			expect:     true,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.identifier.IsChild(c.target)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestIdentifier_AsCommandInvoke(t *testing.T) {
	testCases := []struct {
		name   string
		id     Identifier
		expect string
	}{
		{
			name:   "root",
			id:     ".",
			expect: "",
		},
		{
			name:   "root command",
			id:     ".help",
			expect: "help",
		},
		{
			name:   "module command",
			id:     ".mod.infr.edit",
			expect: "mod infr edit",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.id.AsCommandInvoke()
			assert.Equal(t, c.expect, actual)
		})
	}
}
