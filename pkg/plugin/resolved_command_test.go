package plugin

import (
	"errors"
	"testing"

	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GenerateResolvedCommands(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		repos := []Repository{
			{
				ProviderName: "",
				Modules:      nil,
				Commands: []Command{
					mockCommand{
						name:      "abc",
						hidden:    true,
						throttler: mockThrottler{cmp: "bcd"},
					},
				},
			},
		}

		var nilResolvedModule *ResolvedModule = nil

		expect := []*ResolvedCommand{
			{
				parent:       &nilResolvedModule,
				ProviderName: "",
				Source: mockCommand{
					name:      "abc",
					hidden:    true,
					throttler: mockThrottler{cmp: "bcd"},
				},
				SourceParents:  nil,
				ID:             ".abc",
				Name:           "abc",
				Hidden:         true,
				ChannelTypes:   AllChannels,
				BotPermissions: 0,
				Throttler:      mockThrottler{cmp: "bcd"},
			},
		}

		actual := GenerateResolvedCommands(repos)

		assert.Equal(t, expect, actual)
	})

	t.Run("merge", func(t *testing.T) {
		repos := []Repository{
			{
				ProviderName: "abc",
				Commands:     []Command{mockCommand{name: "def"}},
			},
			{
				ProviderName: "ghi",
				Commands:     []Command{mockCommand{name: "def"}},
			},
		}

		var nilResolvedModule *ResolvedModule = nil

		expect := []*ResolvedCommand{
			{
				parent:       &nilResolvedModule,
				ProviderName: "abc",
				Source:       mockCommand{name: "def"},
				ID:           ".def",
				Name:         "def",
				ChannelTypes: AllChannels,
			},
		}

		actual := GenerateResolvedCommands(repos)
		assert.Equal(t, expect, actual)
	})

	t.Run("merge", func(t *testing.T) {
		repos := []Repository{
			{
				ProviderName: "ghi",
				Commands:     []Command{mockCommand{name: "jkl"}},
			},
			{
				ProviderName: "abc",
				Commands:     []Command{mockCommand{name: "def"}},
			},
		}

		var nilResolvedModule *ResolvedModule = nil

		expect := []*ResolvedCommand{
			{
				parent:       &nilResolvedModule,
				ProviderName: "abc",
				Source:       mockCommand{name: "def"},
				ID:           ".def",
				Name:         "def",
				ChannelTypes: AllChannels,
			},
			{
				parent:       &nilResolvedModule,
				ProviderName: "ghi",
				Source:       mockCommand{name: "jkl"},
				ID:           ".jkl",
				Name:         "jkl",
				ChannelTypes: AllChannels,
			},
		}

		actual := GenerateResolvedCommands(repos)
		assert.Equal(t, expect, actual)
	})

	t.Run("skip duplicates", func(t *testing.T) {
		repos := []Repository{
			{
				ProviderName: "abc",
				Commands:     []Command{mockCommand{name: "def"}},
			},
			{
				ProviderName: "ghi",
				Commands:     []Command{mockCommand{name: "def"}}, // duplicate
			},
		}

		var nilResolvedModule *ResolvedModule = nil

		expect := []*ResolvedCommand{
			{
				parent:       &nilResolvedModule,
				ProviderName: "abc",
				Source:       mockCommand{name: "def"},
				ID:           ".def",
				Name:         "def",
				ChannelTypes: AllChannels,
			},
		}

		actual := GenerateResolvedCommands(repos)
		assert.Equal(t, expect, actual)
	})

	t.Run("remove duplicate aliases", func(t *testing.T) {
		repos := []Repository{
			{
				ProviderName: "abc",
				Commands: []Command{
					mockCommand{
						name:    "def",
						aliases: []string{"ghi", "jkl"},
					},
				},
			},
			{
				ProviderName: "mno",
				Commands: []Command{
					mockCommand{
						name:    "pqr",
						aliases: []string{"jkl", "stu"}, // duplicate alias
					},
				},
			},
		}

		var nilResolvedModule *ResolvedModule = nil

		expect := []*ResolvedCommand{
			{
				parent:       &nilResolvedModule,
				ProviderName: "abc",
				Source: mockCommand{
					name:    "def",
					aliases: []string{"ghi", "jkl"},
				},
				ID:           ".def",
				Name:         "def",
				Aliases:      []string{"ghi", "jkl"},
				ChannelTypes: AllChannels,
			},
			{
				parent:       &nilResolvedModule,
				ProviderName: "mno",
				Source: mockCommand{
					name:    "pqr",
					aliases: []string{"jkl", "stu"},
				},
				ID:           ".pqr",
				Name:         "pqr",
				Aliases:      []string{"stu"},
				ChannelTypes: AllChannels,
			},
		}

		actual := GenerateResolvedCommands(repos)
		assert.Equal(t, expect, actual)
	})
}

func TestResolvedCommand_ShortDescription(t *testing.T) {
	expect := "abc"

	rcmd := &ResolvedCommand{
		Source: mockCommand{
			shortDesc: expect,
		},
	}

	actual := rcmd.ShortDescription(nil)
	assert.Equal(t, expect, actual)
}

func TestResolvedCommand_LongDescription(t *testing.T) {
	t.Run("long description", func(t *testing.T) {
		expect := "abc"

		rcmd := &ResolvedCommand{
			Source: mockCommand{
				longDesc: expect,
			},
		}

		actual := rcmd.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("short description", func(t *testing.T) {
		expect := "abc"

		rcmd := &ResolvedCommand{
			Source: mockCommand{
				shortDesc: expect,
				// no long description defined
			},
		}

		actual := rcmd.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})
}

func TestResolvedCommand_Examples(t *testing.T) {
	expect := []string{"abc", "def"}

	rcmd := &ResolvedCommand{Source: mockCommand{exampleArgs: expect}}

	actual := rcmd.Examples(nil)
	assert.Equal(t, expect, actual)
}

func TestResolvedCommand_IsRestricted(t *testing.T) {
	expect := errors.New("abc")

	rcmd := &ResolvedCommand{
		Source: mockCommand{
			restrictionFunc: func(*state.State, *Context) error {
				return expect
			},
		},
	}

	actual := rcmd.IsRestricted(nil, nil)
	assert.Equal(t, expect, actual)
}

func TestResolvedCommand_Invoke(t *testing.T) {
	expect := "abc"

	rcmd := &ResolvedCommand{
		Source: mockCommand{
			invokeFunc: func(*state.State, *Context) (interface{}, error) {
				return expect, nil
			},
		},
	}

	actual, err := rcmd.Invoke(nil, nil)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}
