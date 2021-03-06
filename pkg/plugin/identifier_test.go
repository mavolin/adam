package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIdentifierFromInvoke(t *testing.T) {
	invoke := "  abc \tdef\njkl  mno"

	var expect ID = ".abc.def.jkl.mno"

	actual := NewIDFromInvoke(invoke)
	assert.Equal(t, expect, actual)
}

func TestIdentifier_Parent(t *testing.T) {
	testCases := []struct {
		name       string
		identifier ID
		expect     []ID
	}{
		{
			name:       "root",
			identifier: ".",
			expect:     []ID{"."},
		},
		{
			name:       "single level",
			identifier: ".mod",
			expect:     []ID{"."},
		},
		{
			name:       "multi level",
			identifier: ".mod.infr.edit",
			expect:     []ID{".mod.infr", ".mod", "."},
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

func TestIdentifier_All(t *testing.T) {
	testCases := []struct {
		name       string
		identifier ID
		expect     []ID
	}{
		{
			name:       "root",
			identifier: ".",
			expect:     []ID{"."},
		},
		{
			name:       "single level",
			identifier: ".mod",
			expect:     []ID{".", ".mod"},
		},
		{
			name:       "multi level",
			identifier: ".mod.infr.edit",
			expect:     []ID{".", ".mod", ".mod.infr", ".mod.infr.edit"},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.identifier.All()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestIdentifier_IsRoot(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		isRoot := ID(".").IsRoot()
		assert.True(t, isRoot)
	})

	t.Run("not root", func(t *testing.T) {
		isRoot := ID(".mod").IsRoot()
		assert.False(t, isRoot)
	})
}

func TestIdentifier_NumParents(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		actual := ID(".").NumParents()
		assert.Equal(t, 0, actual)
	})

	t.Run("no parents", func(t *testing.T) {
		actual := ID(".mod").NumParents()
		assert.Equal(t, 0, actual)
	})

	t.Run("parents", func(t *testing.T) {
		actual := ID(".mod.ban").NumParents()
		assert.Equal(t, 1, actual)
	})
}

func TestIdentifier_IsParent(t *testing.T) {
	testCases := []struct {
		name       string
		identifier ID
		target     ID
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
		identifier ID
		target     ID
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
		id     ID
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
			actual := c.id.AsInvoke()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestIdentifier_Name(t *testing.T) {
	testCases := []struct {
		name       string
		identifier ID
		expect     string
	}{
		{
			name:       "root",
			identifier: ".",
			expect:     "",
		},
		{
			name:       "single level",
			identifier: ".mod",
			expect:     "mod",
		},
		{
			name:       "multi level",
			identifier: ".mod.infr.edit",
			expect:     "edit",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.identifier.Name()
			assert.Equal(t, c.expect, actual)
		})
	}
}
