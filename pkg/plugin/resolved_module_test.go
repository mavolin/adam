package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateResolvedModules(t *testing.T) {
	repos := []Repository{
		{
			ProviderName: "built_in",
			Modules: []Module{
				mockModule{
					name:     "abc",
					commands: []Command{mockCommand{name: "zyx"}},
				},
				mockModule{
					name:     "def",
					commands: []Command{mockCommand{name: "wvu"}},
				},
			},
		},
		{
			ProviderName: "custom_commands",
			Modules: []Module{
				mockModule{
					name:     "def",
					commands: []Command{mockCommand{name: "tsr"}},
				},
				mockModule{
					name:     "ghi",
					commands: []Command{mockCommand{name: "qpo"}},
				},
			},
		},
	}

	abc := &ResolvedModule{
		Sources: []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name:     "abc",
						commands: []Command{mockCommand{name: "zyx"}},
					},
				},
			},
		},
		ID:   ".abc",
		Name: "abc",
	}

	abc.Commands = append(abc.Commands, &ResolvedCommand{
		parent: &abc,
		Source: mockCommand{name: "zyx"},
		SourceParents: []Module{
			mockModule{
				name:     "abc",
				commands: []Command{mockCommand{name: "zyx"}},
			},
		},
		ProviderName: "built_in",
		ID:           ".abc.zyx",
		Name:         "zyx",
		ChannelTypes: AllChannels,
	})

	def := &ResolvedModule{
		Name: "def",
		Sources: []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name:     "def",
						commands: []Command{mockCommand{name: "wvu"}},
					},
				},
			},
			{
				ProviderName: "custom_commands",
				Modules: []Module{
					mockModule{
						name:     "def",
						commands: []Command{mockCommand{name: "tsr"}},
					},
				},
			},
		},
		ID: ".def",
	}

	def.Commands = append(def.Commands, &ResolvedCommand{
		parent: &def,
		Source: mockCommand{name: "tsr"},
		SourceParents: []Module{
			mockModule{
				name:     "def",
				commands: []Command{mockCommand{name: "tsr"}},
			},
		},
		ProviderName: "custom_commands",
		ID:           ".def.tsr",
		Name:         "tsr",
		ChannelTypes: AllChannels,
	})

	def.Commands = append(def.Commands, &ResolvedCommand{
		parent: &def,
		Source: mockCommand{name: "wvu"},
		SourceParents: []Module{
			mockModule{
				name:     "def",
				commands: []Command{mockCommand{name: "wvu"}},
			},
		},
		ProviderName: "built_in",
		ID:           ".def.wvu",
		Name:         "wvu",
		ChannelTypes: AllChannels,
	})

	ghi := &ResolvedModule{
		Name: "ghi",
		Sources: []SourceModule{
			{
				ProviderName: "custom_commands",
				Modules: []Module{
					mockModule{
						name:     "ghi",
						commands: []Command{mockCommand{name: "qpo"}},
					},
				},
			},
		},
		ID: ".ghi",
	}

	ghi.Commands = append(ghi.Commands, &ResolvedCommand{
		parent: &ghi,
		Source: mockCommand{name: "qpo"},
		SourceParents: []Module{
			mockModule{
				name:     "ghi",
				commands: []Command{mockCommand{name: "qpo"}},
			},
		},
		ProviderName: "custom_commands",
		ID:           ".ghi.qpo",
		Name:         "qpo",
		ChannelTypes: AllChannels,
	})

	expect := []*ResolvedModule{abc, def, ghi}

	actual := GenerateResolvedModules(repos)
	assert.Equal(t, expect, actual)
}

func Test_mergeSourceModules(t *testing.T) {
	merge := []SourceModule{
		{
			ProviderName: "built_in",
			Modules: []Module{
				mockModule{
					name:     "def",
					commands: []Command{mockCommand{name: "wvu"}},
				},
			},
		},
		{
			ProviderName: "built_in",
			Modules: []Module{
				mockModule{
					name:     "abc",
					commands: []Command{mockCommand{name: "zyx"}},
				},
			},
		},
		{
			ProviderName: "custom_commands",
			Modules: []Module{
				mockModule{
					name:     "def",
					commands: []Command{mockCommand{name: "tsr"}},
				},
			},
		},
	}

	expect := [][]SourceModule{
		{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name:     "abc",
						commands: []Command{mockCommand{name: "zyx"}},
					},
				},
			},
		},
		{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name:     "def",
						commands: []Command{mockCommand{name: "wvu"}},
					},
				},
			},
			{
				ProviderName: "custom_commands",
				Modules: []Module{
					mockModule{
						name:     "def",
						commands: []Command{mockCommand{name: "tsr"}},
					},
				},
			},
		},
	}

	actual := sortSourceModules(merge)

	assert.Equal(t, expect, actual)
}

func Test_generateResolvedModule(t *testing.T) {
	t.Run("no parent", func(t *testing.T) {
		smods := []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name:     "abc",
						commands: []Command{mockCommand{name: "def"}},
					},
				},
			},
		}

		expect := &ResolvedModule{
			Parent:  nil,
			Sources: smods,
			ID:      ".abc",
			Name:    "abc",
		}

		expect.Commands = append(expect.Commands, &ResolvedCommand{
			parent:   &expect,
			provider: nil,
			Source:   mockCommand{name: "def"},
			SourceParents: []Module{
				mockModule{
					name:     "abc",
					commands: []Command{mockCommand{name: "def"}},
				},
			},
			ProviderName: "built_in",
			ID:           ".abc.def",
			Name:         "def",
			ChannelTypes: AllChannels,
		})

		actual := generateResolvedModule(nil, smods, nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("parent", func(t *testing.T) {
		parent := &ResolvedModule{
			Parent: nil,
			Sources: []SourceModule{
				{
					ProviderName: "built_in",
					Modules: []Module{
						mockModule{
							name: "abc",
							modules: []Module{
								mockModule{
									name:     "def",
									commands: []Command{mockCommand{name: "ghi"}},
								},
							},
						},
					},
				},
			},
			ID:   ".abc",
			Name: "abc",
		}

		smods := []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name: "abc",
						modules: []Module{
							mockModule{
								name:     "def",
								commands: []Command{mockCommand{name: "ghi"}},
							},
						},
					},
					mockModule{
						name:     "def",
						commands: []Command{mockCommand{name: "ghi"}},
					},
				},
			},
		}

		expect := &ResolvedModule{
			Parent:  parent,
			Sources: smods,
			ID:      ".abc.def",
			Name:    "def",
		}

		expect.Commands = append(expect.Commands, &ResolvedCommand{
			parent:   &expect,
			provider: nil,
			Source:   mockCommand{name: "ghi"},
			SourceParents: []Module{
				mockModule{
					name: "abc",
					modules: []Module{
						mockModule{
							name:     "def",
							commands: []Command{mockCommand{name: "ghi"}},
						},
					},
				},
				mockModule{
					name:     "def",
					commands: []Command{mockCommand{name: "ghi"}},
				},
			},
			ProviderName: "built_in",
			ID:           ".abc.def.ghi",
			Name:         "ghi",
			ChannelTypes: AllChannels,
		})

		actual := generateResolvedModule(parent, smods, nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("children hidden", func(t *testing.T) {
		smods := []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name:     "abc",
						commands: []Command{mockCommand{name: "def", hidden: true}},
					},
				},
			},
			{
				ProviderName: "other",
				Modules: []Module{
					mockModule{
						name:     "abc",
						commands: []Command{mockCommand{name: "ghi", hidden: true}},
					},
				},
			},
		}

		expect := &ResolvedModule{
			Parent:  nil,
			Sources: smods,
			ID:      ".abc",
			Name:    "abc",
			Hidden:  true,
		}

		expect.Commands = append(expect.Commands, &ResolvedCommand{
			parent:   &expect,
			provider: nil,
			Source:   mockCommand{name: "def", hidden: true},
			SourceParents: []Module{
				mockModule{
					name:     "abc",
					commands: []Command{mockCommand{name: "def", hidden: true}},
				},
			},
			ProviderName: "built_in",
			ID:           ".abc.def",
			Name:         "def",
			Hidden:       true,
			ChannelTypes: AllChannels,
		})

		expect.Commands = append(expect.Commands, &ResolvedCommand{
			parent:   &expect,
			provider: nil,
			Source:   mockCommand{name: "ghi", hidden: true},
			SourceParents: []Module{
				mockModule{
					name:     "abc",
					commands: []Command{mockCommand{name: "ghi", hidden: true}},
				},
			},
			ProviderName: "other",
			ID:           ".abc.ghi",
			Name:         "ghi",
			Hidden:       true,
			ChannelTypes: AllChannels,
		})

		actual := generateResolvedModule(nil, smods, nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("not hidden", func(t *testing.T) {
		smods := []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name: "abc",
						commands: []Command{
							mockCommand{name: "def", hidden: true},
							mockCommand{name: "ghi", hidden: false},
						},
					},
				},
			},
			{
				ProviderName: "other",
				Modules: []Module{
					mockModule{
						name:     "abc",
						commands: []Command{mockCommand{name: "jkl", hidden: false}},
					},
				},
			},
		}

		expect := &ResolvedModule{
			Parent:  nil,
			Sources: smods,
			ID:      ".abc",
			Name:    "abc",
			Hidden:  false,
		}

		expect.Commands = append(expect.Commands, &ResolvedCommand{
			parent:   &expect,
			provider: nil,
			Source:   mockCommand{name: "def", hidden: true},
			SourceParents: []Module{
				mockModule{
					name: "abc",
					commands: []Command{
						mockCommand{name: "def", hidden: true},
						mockCommand{name: "ghi", hidden: false},
					},
				},
			},
			ProviderName: "built_in",
			ID:           ".abc.def",
			Name:         "def",
			Hidden:       true,
			ChannelTypes: AllChannels,
		})

		expect.Commands = append(expect.Commands, &ResolvedCommand{
			parent:   &expect,
			provider: nil,
			Source:   mockCommand{name: "ghi", hidden: false},
			SourceParents: []Module{
				mockModule{
					name: "abc",
					commands: []Command{
						mockCommand{name: "def", hidden: true},
						mockCommand{name: "ghi", hidden: false},
					},
				},
			},
			ProviderName: "built_in",
			ID:           ".abc.ghi",
			Name:         "ghi",
			Hidden:       false,
			ChannelTypes: AllChannels,
		})

		expect.Commands = append(expect.Commands, &ResolvedCommand{
			parent:   &expect,
			provider: nil,
			Source:   mockCommand{name: "jkl", hidden: false},
			SourceParents: []Module{
				mockModule{
					name:     "abc",
					commands: []Command{mockCommand{name: "jkl", hidden: false}},
				},
			},
			ProviderName: "other",
			ID:           ".abc.jkl",
			Name:         "jkl",
			ChannelTypes: AllChannels,
		})

		actual := generateResolvedModule(nil, smods, nil)
		assert.Equal(t, expect, actual)
	})
}

func Test_fillSubmodules(t *testing.T) {
	parent := &ResolvedModule{
		ID:   ".abc",
		Name: "abc",
		Sources: []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name: "abc",
						modules: []Module{
							mockModule{name: "ghi"}, // should get sorted
							mockModule{name: "def"},
						},
					},
				},
			},
		},
	}

	expect := &ResolvedModule{
		ID:   ".abc",
		Name: "abc",
		Sources: []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name: "abc",
						modules: []Module{
							mockModule{name: "ghi"},
							mockModule{name: "def"},
						},
					},
				},
			},
		},
		Modules: []*ResolvedModule{
			{
				Parent: parent,
				Sources: []SourceModule{
					{
						ProviderName: "built_in",
						Modules: []Module{
							mockModule{
								name: "abc",
								modules: []Module{
									mockModule{name: "ghi"},
									mockModule{name: "def"},
								},
							},
							mockModule{name: "def"},
						},
					},
				},
				ID:     ".abc.def",
				Name:   "def",
				Hidden: true,
			},
			{
				Parent: parent,
				Sources: []SourceModule{
					{
						ProviderName: "built_in",
						Modules: []Module{
							mockModule{
								name: "abc",
								modules: []Module{
									mockModule{name: "ghi"},
									mockModule{name: "def"},
								},
							},
							mockModule{name: "ghi"},
						},
					},
				},
				ID:     ".abc.ghi",
				Name:   "ghi",
				Hidden: true,
			},
		},
	}

	fillSubmodules(parent, nil)
	assert.Equal(t, expect, parent)
}

func Test_fillSubcommands(t *testing.T) {
	parent := &ResolvedModule{
		ID:   ".abc",
		Name: "abc",
		Sources: []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name: "abc",
						commands: []Command{
							mockCommand{
								name:    "ghi", // should get sorted
								aliases: []string{"mno", "pqr"},
							},
							mockCommand{
								name:    "def",
								aliases: []string{"jkl", "mno"}, // duplicate alias
								hidden:  true,
							},
						},
					},
				},
			},
			{
				ProviderName: "custom_commands",
				Modules: []Module{
					mockModule{
						name: "abc",
						commands: []Command{
							mockCommand{name: "def"}, // duplicate name
						},
					},
				},
			},
		},
	}

	expect := &ResolvedModule{
		Sources: []SourceModule{
			{
				ProviderName: "built_in",
				Modules: []Module{
					mockModule{
						name: "abc",
						commands: []Command{
							mockCommand{
								name:    "ghi",
								aliases: []string{"mno", "pqr"},
							},
							mockCommand{
								name:    "def",
								aliases: []string{"jkl", "mno"},
								hidden:  true,
							},
						},
					},
				},
			},
			{
				ProviderName: "custom_commands",
				Modules: []Module{
					mockModule{
						name:     "abc",
						commands: []Command{mockCommand{name: "def"}},
					},
				},
			},
		},
		ID:   ".abc",
		Name: "abc",
		Commands: []*ResolvedCommand{
			{
				parent:   &parent,
				provider: nil,
				Source: mockCommand{
					name:    "def",
					aliases: []string{"jkl", "mno"},
					hidden:  true,
				},
				SourceParents: []Module{
					mockModule{
						name: "abc",
						commands: []Command{
							mockCommand{
								name:    "ghi",
								aliases: []string{"mno", "pqr"},
							},
							mockCommand{
								name:    "def",
								aliases: []string{"jkl", "mno"},
								hidden:  true,
							},
						},
					},
				},
				ProviderName: "built_in",
				ID:           ".abc.def",
				Name:         "def",
				Aliases:      []string{"jkl"},
				Hidden:       true,
				ChannelTypes: AllChannels,
			},
			{
				parent:   &parent,
				provider: nil,
				Source: mockCommand{
					name:    "ghi",
					aliases: []string{"mno", "pqr"},
				},
				SourceParents: []Module{
					mockModule{
						name: "abc",
						commands: []Command{
							mockCommand{
								name:    "ghi",
								aliases: []string{"mno", "pqr"},
							},
							mockCommand{
								name:    "def",
								aliases: []string{"jkl", "mno"},
								hidden:  true,
							},
						},
					},
				},
				ProviderName: "built_in",
				ID:           ".abc.ghi",
				Name:         "ghi",
				Aliases:      []string{"mno", "pqr"},
				ChannelTypes: AllChannels,
			},
		},
	}

	fillSubcommands(parent)
	assert.Equal(t, expect, parent)
}

func Test_generateResolvedCommands(t *testing.T) {
	parent := new(ResolvedModule)

	smod := SourceModule{
		ProviderName: "",
		Modules: []Module{
			mockModule{
				name: "abc",
				commands: []Command{
					mockCommand{
						name:            "def",
						hidden:          true,
						restrictionFunc: nil,
						throttler:       mockThrottler{cmp: "cde"},
					},
				},
			},
		},
	}

	expect := []*ResolvedCommand{
		{
			parent:       &parent,
			provider:     nil,
			ProviderName: "",
			Source: mockCommand{
				name:            "def",
				hidden:          true,
				restrictionFunc: nil,
				throttler:       mockThrottler{cmp: "cde"},
			},
			SourceParents: []Module{
				mockModule{
					name: "abc",
					commands: []Command{
						mockCommand{
							name:            "def",
							hidden:          true,
							restrictionFunc: nil,
							throttler:       mockThrottler{cmp: "cde"},
						},
					},
				},
			},
			ID:           ".abc.def",
			Name:         "def",
			Hidden:       true,
			ChannelTypes: AllChannels,
			Throttler:    mockThrottler{cmp: "cde"},
		},
	}

	actual := generateResolvedCommands(parent, smod)

	assert.Equal(t, expect, actual)
}

func TestResolvedModule_ShortDescription(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := "abc"

		rmod := &ResolvedModule{
			Sources: []SourceModule{
				{Modules: []Module{mockModule{shortDesc: expect}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("fallback", func(t *testing.T) {
		expect := "def"

		rmod := &ResolvedModule{
			Sources: []SourceModule{
				{Modules: []Module{mockModule{shortDesc: ""}}},
				{Modules: []Module{mockModule{shortDesc: expect}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("none", func(t *testing.T) {
		rmod := &ResolvedModule{
			Sources: []SourceModule{
				{Modules: []Module{mockModule{shortDesc: ""}}},
				{Modules: []Module{mockModule{shortDesc: ""}}},
			},
		}

		actual := rmod.ShortDescription(nil)
		assert.Empty(t, actual)
	})
}

func TestResolvedModule_LongDescription(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := "abc"

		rmod := &ResolvedModule{
			Sources: []SourceModule{
				{Modules: []Module{mockModule{longDesc: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("fallback", func(t *testing.T) {
		expect := "def"

		rmod := &ResolvedModule{
			Sources: []SourceModule{
				{Modules: []Module{mockModule{}}},
				{Modules: []Module{mockModule{longDesc: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("short description", func(t *testing.T) {
		expect := "abc"

		rmod := &ResolvedModule{
			Sources: []SourceModule{
				{Modules: []Module{mockModule{}}},
				{Modules: []Module{mockModule{shortDesc: expect}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Equal(t, expect, actual)
	})

	t.Run("none", func(t *testing.T) {
		rmod := &ResolvedModule{
			Sources: []SourceModule{
				{Modules: []Module{mockModule{}}},
				{Modules: []Module{mockModule{longDesc: ""}}},
			},
		}

		actual := rmod.LongDescription(nil)
		assert.Empty(t, actual)
	})
}

func TestResolvedModule_FindCommand(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		expect := &ResolvedCommand{
			Name: "def",
		}

		rmod := &ResolvedModule{
			Commands: []*ResolvedCommand{{Name: "abc"}, expect, {Name: "ghi"}},
		}

		actual := rmod.FindCommand(expect.Name)
		assert.Equal(t, expect, actual)
	})

	t.Run("alias", func(t *testing.T) {
		expect := &ResolvedCommand{
			Name:    "def",
			Aliases: []string{"mno"},
		}

		rmod := &ResolvedModule{
			Commands: []*ResolvedCommand{
				{Name: "abc", Aliases: []string{"jkl"}},
				expect,
				{Name: "ghi"},
			},
		}

		actual := rmod.FindCommand(expect.Aliases[0])
		assert.Equal(t, expect, actual)
	})

	t.Run("not found", func(t *testing.T) {
		actual := new(ResolvedModule).FindCommand("abc")
		assert.Nil(t, actual)
	})
}

func TestResolvedModule_FindModule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := &ResolvedModule{
			Name: "def",
		}

		rmod := &ResolvedModule{
			Modules: []*ResolvedModule{{Name: "abc"}, expect, {Name: "ghi"}},
		}

		actual := rmod.FindModule(expect.Name)
		assert.Equal(t, expect, actual)
	})

	t.Run("not found", func(t *testing.T) {
		actual := new(ResolvedModule).FindModule("abc")
		assert.Nil(t, actual)
	})
}
