package plugin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIDFromInvoke(t *testing.T) {
	invoke := "  abc  \ndef\njkl  mno"

	var expect ID = ".abc.def.jkl.mno"

	actual := NewIDFromInvoke(invoke)
	assert.Equal(t, expect, actual)
}

func TestID_Parent(t *testing.T) {
	testCases := []struct {
		name   string
		ID     ID
		expect []ID
	}{
		{
			name:   "root",
			ID:     ".",
			expect: []ID{"."},
		},
		{
			name:   "single level",
			ID:     ".mod",
			expect: []ID{"."},
		},
		{
			name:   "multi level",
			ID:     ".mod.infr.edit",
			expect: []ID{".mod.infr", ".mod", "."},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.ID

			for lvl, expect := range c.expect {
				actual = actual.Parent()
				assert.Equal(t, expect, actual, "unexpected ID returned on level %d", lvl+1)
			}
		})
	}
}

func TestID_All(t *testing.T) {
	testCases := []struct {
		name   string
		ID     ID
		expect []ID
	}{
		{
			name:   "root",
			ID:     ".",
			expect: []ID{"."},
		},
		{
			name:   "single level",
			ID:     ".mod",
			expect: []ID{".", ".mod"},
		},
		{
			name:   "multi level",
			ID:     ".mod.infr.edit",
			expect: []ID{".", ".mod", ".mod.infr", ".mod.infr.edit"},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.ID.All()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestID_IsRoot(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		isRoot := ID(".").IsRoot()
		assert.True(t, isRoot)
	})

	t.Run("not root", func(t *testing.T) {
		isRoot := ID(".mod").IsRoot()
		assert.False(t, isRoot)
	})
}

func TestID_NumParents(t *testing.T) {
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

func TestID_IsParent(t *testing.T) {
	testCases := []struct {
		name   string
		ID     ID
		target ID
		expect bool
	}{
		{
			name:   "parent",
			ID:     ".mod.ban",
			target: ".mod",
			expect: true,
		},
		{
			name:   "equal",
			ID:     ".mod.ban",
			target: ".mod.ban",
			expect: false,
		},
		{
			name:   "child",
			ID:     ".mod",
			target: ".mod.ban",
			expect: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.ID.IsParent(c.target)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestID_IsChild(t *testing.T) {
	testCases := []struct {
		name   string
		ID     ID
		target ID
		expect bool
	}{
		{
			name:   "parent",
			ID:     ".mod.ban",
			target: ".mod",
			expect: false,
		},
		{
			name:   "equal",
			ID:     ".mod.ban",
			target: ".mod.ban",
			expect: false,
		},
		{
			name:   "child",
			ID:     ".mod",
			target: ".mod.ban",
			expect: true,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.ID.IsChild(c.target)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestID_AsCommandInvoke(t *testing.T) {
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

func TestID_Name(t *testing.T) {
	testCases := []struct {
		name   string
		ID     ID
		expect string
	}{
		{
			name:   "root",
			ID:     ".",
			expect: "",
		},
		{
			name:   "single level",
			ID:     ".mod",
			expect: "mod",
		},
		{
			name:   "multi level",
			ID:     ".mod.infr.edit",
			expect: "edit",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.ID.Name()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func ExampleID_AsInvoke() {
	var id ID = ".mod.ban"

	fmt.Println(id.AsInvoke())
	// Output: mod ban
}
