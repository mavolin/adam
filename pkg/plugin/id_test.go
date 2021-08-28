package plugin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIDFromInvoke(t *testing.T) {
	t.Parallel()

	invoke := "  abc  \ndef\njkl  mno"

	var expect ID = ".abc.def.jkl.mno"

	actual := NewIDFromInvoke(invoke)
	assert.Equal(t, expect, actual)
}

func TestID_Parent(t *testing.T) {
	t.Parallel()

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
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := c.ID

			for lvl, expect := range c.expect {
				actual = actual.Parent()
				assert.Equal(t, expect, actual, "unexpected id returned on level %d", lvl+1)
			}
		})
	}
}

func TestID_All(t *testing.T) {
	t.Parallel()

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
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := c.ID.All()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestID_IsRoot(t *testing.T) {
	t.Parallel()

	t.Run("root", func(t *testing.T) {
		t.Parallel()

		isRoot := ID(".").IsRoot()
		assert.True(t, isRoot)
	})

	t.Run("not root", func(t *testing.T) {
		t.Parallel()

		isRoot := ID(".mod").IsRoot()
		assert.False(t, isRoot)
	})
}

func TestID_NumParents(t *testing.T) {
	t.Parallel()

	t.Run("root", func(t *testing.T) {
		t.Parallel()

		actual := ID(".").NumParents()
		assert.Equal(t, 0, actual)
	})

	t.Run("no parents", func(t *testing.T) {
		t.Parallel()

		actual := ID(".mod").NumParents()
		assert.Equal(t, 0, actual)
	})

	t.Run("parents", func(t *testing.T) {
		t.Parallel()

		actual := ID(".mod.ban").NumParents()
		assert.Equal(t, 1, actual)
	})
}

func TestID_IsParentOf(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		id     ID
		target ID
		expect bool
	}{
		{
			name:   "parent",
			id:     ".mod",
			target: ".mod.ban",
			expect: true,
		},
		{
			name:   "equal",
			id:     ".mod.ban",
			target: ".mod.ban",
			expect: false,
		},
		{
			name:   "child",
			id:     ".mod.ban",
			target: ".mod",
			expect: false,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := c.id.IsParentOf(c.target)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestID_IsChild(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		id     ID
		target ID
		expect bool
	}{
		{
			name:   "parent",
			id:     ".mod",
			target: ".mod.ban",
			expect: false,
		},
		{
			name:   "equal",
			id:     ".mod.ban",
			target: ".mod.ban",
			expect: false,
		},
		{
			name:   "child",
			id:     ".mod.ban",
			target: ".mod",
			expect: true,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := c.id.IsChildOf(c.target)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestID_AsCommandInvoke(t *testing.T) {
	t.Parallel()

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
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := c.id.AsInvoke()
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestID_Name(t *testing.T) {
	t.Parallel()

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
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

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
