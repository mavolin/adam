package resolved

import (
	"errors"
	"fmt"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
)

func TestPluginProvider_PluginSources(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		r := NewPluginResolver(nil)
		r.AddSource("another",
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return []plugin.Command{mockCommand{name: "def"}}, []plugin.Module{mockModule{name: "ghi"}}, nil
			})

		p := &PluginProvider{
			commands: make([]plugin.ResolvedCommand, 0),
			resolver: r,
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockCommand{name: "abc"},
					},
				},
			},
			usedNames: make(map[string]struct{}),
		}

		expect := []plugin.Source{
			{
				Name: plugin.BuiltInSource,
				Commands: []plugin.Command{
					mockCommand{name: "abc"},
				},
			},
			{
				Name: "another",
				Commands: []plugin.Command{
					mockCommand{name: "def"},
				},
				Modules: []plugin.Module{
					mockModule{name: "ghi"},
				},
			},
		}

		actual := p.PluginSources()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := &PluginProvider{
			resolver: NewPluginResolver(nil),
			commands: make([]plugin.ResolvedCommand, 0),
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockCommand{name: "abc"},
					},
				},
				{
					Name: "another",
					Commands: []plugin.Command{
						mockCommand{name: "def"},
					},
					Modules: []plugin.Module{
						mockModule{name: "ghi"},
					},
				},
			},
		}

		expect := p.sources

		actual := p.PluginSources()
		assert.Equal(t, expect, actual)
	})
}

func TestPluginProvider_Commands(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		r := NewPluginResolver(nil)
		r.AddSource("another",
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return []plugin.Command{
					mockCommand{
						name:         "def",
						channelTypes: plugin.GuildChannels,
					},
				}, nil, nil
			})

		p := &PluginProvider{
			resolver: r,
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockCommand{
							name:         "abc",
							channelTypes: plugin.GuildChannels,
						},
					},
				},
			},
			usedNames: make(map[string]struct{}),
		}

		p.commands = append(p.commands, &Command{
			provider:   p,
			sourceName: plugin.BuiltInSource,
			source: mockCommand{
				name:         "abc",
				channelTypes: plugin.GuildChannels,
			},
			id: ".abc",
		})

		expect := []plugin.ResolvedCommand{
			&Command{
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source: mockCommand{
					name:         "abc",
					channelTypes: plugin.GuildChannels,
				},
				id: ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source: mockCommand{
					name:         "def",
					channelTypes: plugin.GuildChannels,
				},
				id: ".def",
			},
		}

		actual := p.Commands()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := new(PluginProvider)
		p.resolver = NewPluginResolver(nil)
		p.sources = []plugin.Source{
			{
				Name:     plugin.BuiltInSource,
				Commands: []plugin.Command{mockCommand{name: "abc"}},
			},
			{
				Name:     "another",
				Commands: []plugin.Command{mockCommand{name: "def"}},
			},
		}
		p.commands = []plugin.ResolvedCommand{
			&Command{
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source:     mockCommand{name: "abc"},
				id:         ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "def"},
				id:         ".def",
			},
		}

		expect := p.commands

		actual := p.Commands()
		assert.Equal(t, expect, actual)
	})

	t.Run("resolver", func(t *testing.T) {
		t.Run("single", func(t *testing.T) {
			sources := []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockCommand{
							name:         "abc",
							hidden:       true,
							throttler:    mockThrottler{cmp: "bcd"},
							channelTypes: plugin.AllChannels,
						},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: plugin.BuiltInSource,
					source: mockCommand{
						name:         "abc",
						hidden:       true,
						throttler:    mockThrottler{cmp: "bcd"},
						channelTypes: plugin.AllChannels,
					},
					sourceParents: nil,
					id:            ".abc",
				},
			}

			assert.Equal(t, expect, p.Commands())
		})

		t.Run("merge", func(t *testing.T) {
			sources := []plugin.Source{
				{
					Name: "abc",
					Commands: []plugin.Command{
						mockCommand{name: "def", channelTypes: plugin.AllChannels},
					},
				},
				{
					Name: "ghi",
					Commands: []plugin.Command{
						mockCommand{name: "def", channelTypes: plugin.AllChannels},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source: mockCommand{
						name:         "def",
						channelTypes: plugin.AllChannels,
					},
					id: ".def",
				},
			}

			assert.Equal(t, expect, p.Commands())
		})

		t.Run("merge", func(t *testing.T) {
			sources := []plugin.Source{
				{
					Name: "ghi",
					Commands: []plugin.Command{
						mockCommand{name: "jkl", channelTypes: plugin.AllChannels},
					},
				},
				{
					Name: "abc",
					Commands: []plugin.Command{
						mockCommand{name: "def", channelTypes: plugin.AllChannels},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source:     mockCommand{name: "def", channelTypes: plugin.AllChannels},
					id:         ".def",
				},
				&Command{
					provider:   p,
					sourceName: "ghi",
					source:     mockCommand{name: "jkl", channelTypes: plugin.AllChannels},
					id:         ".jkl",
				},
			}

			assert.Equal(t, expect, p.Commands())
		})

		t.Run("skip duplicates", func(t *testing.T) {
			sources := []plugin.Source{
				{
					Name: "abc",
					Commands: []plugin.Command{
						mockCommand{name: "def", channelTypes: plugin.AllChannels},
					},
				},
				{
					Name: "ghi",
					Commands: []plugin.Command{
						mockCommand{name: "def", channelTypes: plugin.AllChannels}, // duplicate
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source:     mockCommand{name: "def", channelTypes: plugin.AllChannels},
					id:         ".def",
				},
			}

			assert.Equal(t, expect, p.Commands())
		})

		t.Run("remove duplicate aliases", func(t *testing.T) {
			sources := []plugin.Source{
				{
					Name: "abc",
					Commands: []plugin.Command{
						mockCommand{
							name:         "def",
							aliases:      []string{"ghi", "jkl"},
							channelTypes: plugin.AllChannels,
						},
					},
				},
				{
					Name: "mno",
					Commands: []plugin.Command{
						mockCommand{
							name:         "pqr",
							aliases:      []string{"jkl", "stu"}, // duplicate alias
							channelTypes: plugin.AllChannels,
						},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source: mockCommand{
						name:         "def",
						aliases:      []string{"ghi", "jkl"},
						channelTypes: plugin.AllChannels,
					},
					id:      ".def",
					aliases: []string{"ghi", "jkl"},
				},
				&Command{
					provider:   p,
					sourceName: "mno",
					source: mockCommand{
						name:         "pqr",
						aliases:      []string{"jkl", "stu"},
						channelTypes: plugin.AllChannels,
					},
					id:      ".pqr",
					aliases: []string{"stu"},
				},
			}

			assert.Equal(t, expect, p.Commands())
		})
	})
}

func TestPluginProvider_Modules(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		r := NewPluginResolver(nil)
		r.AddSource("another",
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return nil, []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{
								name:         "jkl",
								channelTypes: plugin.GuildChannels,
							},
						},
					},
				}, nil
			})
		p := &PluginProvider{
			resolver: r,
			commands: make([]plugin.ResolvedCommand, 0),
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Modules: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "def",
									channelTypes: plugin.GuildChannels,
								},
							},
						},
					},
				},
			},
			usedNames: make(map[string]struct{}),
		}

		p.modules = []plugin.ResolvedModule{
			&Module{
				sources: []plugin.SourceModule{
					{
						SourceName: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{
										name:         "def",
										channelTypes: plugin.GuildChannels,
									},
								},
							},
						},
					},
				},
				id:     ".abc",
				hidden: false,
			},
		}

		p.modules[0].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     p.modules[0],
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source: mockCommand{
					name:         "def",
					channelTypes: plugin.GuildChannels,
				},
				sourceParents: []plugin.Module{
					mockModule{
						name: "abc",
						commands: []plugin.Command{
							mockCommand{
								name:         "def",
								channelTypes: plugin.GuildChannels,
							},
						},
					},
				},
				id:      ".abc.def",
				aliases: nil,
			},
		}

		expect := make([]plugin.ResolvedModule, len(p.modules))
		copy(expect, p.modules)

		expect = append(expect, &Module{
			parent: nil,
			sources: []plugin.SourceModule{
				{
					SourceName: "another",
					Modules: []plugin.Module{
						mockModule{
							name: "ghi",
							commands: []plugin.Command{
								mockCommand{
									name:         "jkl",
									channelTypes: plugin.GuildChannels,
								},
							},
						},
					},
				},
			},
			id:     ".ghi",
			hidden: false,
		})

		expect[1].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     expect[1],
				provider:   p,
				sourceName: "another",
				source: mockCommand{
					name:         "jkl",
					channelTypes: plugin.GuildChannels,
				},
				sourceParents: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{
								name:         "jkl",
								channelTypes: plugin.GuildChannels,
							},
						},
					},
				},
				id: ".ghi.jkl",
			},
		}

		actual := p.Modules()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := new(PluginProvider)
		p.resolver = NewPluginResolver(nil)

		p.usedNames = make(map[string]struct{})

		p.sources = []plugin.Source{
			{
				Name: plugin.BuiltInSource,
				Modules: []plugin.Module{
					mockModule{
						name: "abc",
						commands: []plugin.Command{
							mockCommand{name: "def"},
						},
					},
				},
			},
			{
				Name: "another",
				Modules: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
						},
					},
				},
			},
		}

		p.modules = []plugin.ResolvedModule{
			&Module{
				parent: nil,
				sources: []plugin.SourceModule{
					{
						SourceName: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{name: "def"},
								},
							},
						},
					},
				},
				id:     ".abc",
				hidden: false,
			},
			&Module{
				sources: []plugin.SourceModule{
					{
						SourceName: "another",
						Modules: []plugin.Module{
							mockModule{
								name: "ghi",
								commands: []plugin.Command{
									mockCommand{name: "jkl"},
								},
							},
						},
					},
				},
				id: ".ghi",
			},
		}

		p.modules[0].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     p.modules[0],
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source:     mockCommand{name: "def"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "abc",
						commands: []plugin.Command{
							mockCommand{name: "def"},
						},
					},
				},
				id: ".abc.def",
			},
		}

		p.modules[1].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     p.modules[1],
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "jkl"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
						},
					},
				},
				id: "ghi.jkl",
			},
		}

		expect := p.modules

		actual := p.Modules()
		assert.Equal(t, expect, actual)
	})

	t.Run("resolver", func(t *testing.T) {
		t.Run("generate", func(t *testing.T) {
			t.Run("", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{name: "zyx", channelTypes: plugin.AllChannels},
								},
							},
							mockModule{
								name: "def",
								commands: []plugin.Command{
									mockCommand{name: "wvu", channelTypes: plugin.AllChannels},
								},
							},
						},
					},
					{
						Name: "custom_commands",
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{name: "zyx", channelTypes: plugin.AllChannels},
								},
							},
							mockModule{
								name: "def",
								commands: []plugin.Command{
									mockCommand{name: "tsr", channelTypes: plugin.AllChannels},
								},
							},
							mockModule{
								name: "ghi",
								commands: []plugin.Command{
									mockCommand{name: "qpo", channelTypes: plugin.AllChannels},
								},
							},
						},
					},
				}

				p := newProviderFromSources(sources)

				abc := &Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									commands: []plugin.Command{
										mockCommand{name: "zyx", channelTypes: plugin.AllChannels},
									},
								},
							},
						},
						{
							SourceName: "custom_commands",
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									commands: []plugin.Command{
										mockCommand{name: "zyx", channelTypes: plugin.AllChannels},
									},
								},
							},
						},
					},
					id: ".abc",
				}

				abc.commands = append(abc.commands, &Command{
					provider: p,
					parent:   abc,
					source: mockCommand{
						name:         "zyx",
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "zyx",
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".abc.zyx",
				})

				def := &Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{
											name:         "wvu",
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
						{
							SourceName: "custom_commands",
							Modules: []plugin.Module{
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{
											name:         "tsr",
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
					},
					id: ".def",
				}

				def.commands = append(def.commands, &Command{
					provider: p,
					parent:   def,
					source: mockCommand{
						name:         "tsr",
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{
									name:         "tsr",
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: "custom_commands",
					id:         ".def.tsr",
				})

				def.commands = append(def.commands, &Command{
					provider: p,
					parent:   def,
					source: mockCommand{
						name:         "wvu",
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{
									name:         "wvu",
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".def.wvu",
				})

				ghi := &Module{
					sources: []plugin.SourceModule{
						{
							SourceName: "custom_commands",
							Modules: []plugin.Module{
								mockModule{
									name: "ghi",
									commands: []plugin.Command{
										mockCommand{
											name:         "qpo",
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
					},
					id: ".ghi",
				}

				ghi.commands = append(ghi.commands, &Command{
					provider: p,
					parent:   ghi,
					source: mockCommand{
						name:         "qpo",
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "ghi",
							commands: []plugin.Command{
								mockCommand{
									name:         "qpo",
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: "custom_commands",
					id:         ".ghi.qpo",
				})

				expect := []plugin.ResolvedModule{abc, def, ghi}
				assert.Equal(t, expect, p.Modules(), "all")
			})

			t.Run("no parent", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{
										name:         "def",
										channelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				}

				p := newProviderFromSources(sources)

				expect := &Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									commands: []plugin.Command{
										mockCommand{
											name:         "def",
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
					},
					id: ".abc",
				}

				expect.commands = append(expect.commands, &Command{
					provider: p,
					parent:   expect,
					source: mockCommand{
						name:         "def",
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "def",
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".abc.def",
				})

				assert.Equal(t, expect, p.Modules()[0])
			})

			t.Run("parent", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								modules: []plugin.Module{
									mockModule{
										name: "def",
										commands: []plugin.Command{
											mockCommand{
												name:         "ghi",
												channelTypes: plugin.AllChannels,
											},
										},
									},
								},
							},
						},
					},
				}

				p := newProviderFromSources(sources)

				parent := &Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									modules: []plugin.Module{
										mockModule{
											name: "def",
											commands: []plugin.Command{
												mockCommand{
													name:         "ghi",
													channelTypes: plugin.AllChannels,
												},
											},
										},
									},
								},
							},
						},
					},
					id: ".abc",
				}

				expect := &Module{
					parent: parent,
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									modules: []plugin.Module{
										mockModule{
											name: "def",
											commands: []plugin.Command{
												mockCommand{
													name:         "ghi",
													channelTypes: plugin.AllChannels,
												},
											},
										},
									},
								},
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{
											name:         "ghi",
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
					},
					id: ".abc.def",
				}

				parent.modules = append(parent.modules, expect)

				expect.commands = append(expect.commands, &Command{
					provider: p,
					parent:   expect,
					source:   mockCommand{name: "ghi", channelTypes: plugin.AllChannels},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							modules: []plugin.Module{
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{name: "ghi", channelTypes: plugin.AllChannels},
									},
								},
							},
						},
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{
									name:         "ghi",
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".abc.def.ghi",
				})

				assert.Equal(t, expect, p.Command(".abc.def.ghi").Parent())
			})

			t.Run("children hidden", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{
										name:         "def",
										hidden:       true,
										channelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
					{
						Name: "other",
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{
										name:         "ghi",
										hidden:       true,
										channelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				}

				p := newProviderFromSources(sources)

				expect := &Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									commands: []plugin.Command{
										mockCommand{
											name:         "def",
											hidden:       true,
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
						{
							SourceName: "other",
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									commands: []plugin.Command{
										mockCommand{
											name:         "ghi",
											hidden:       true,
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
					},
					id:     ".abc",
					hidden: true,
				}

				expect.commands = append(expect.commands, &Command{
					provider: p,
					parent:   expect,
					source: mockCommand{
						name:         "def",
						hidden:       true,
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "def",
									hidden:       true,
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".abc.def",
				})

				expect.commands = append(expect.commands, &Command{
					provider: p,
					parent:   expect,
					source: mockCommand{
						name:         "ghi",
						hidden:       true,
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "ghi",
									hidden:       true,
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: "other",
					id:         ".abc.ghi",
				})

				assert.Equal(t, expect, p.Modules()[0])
			})

			t.Run("not hidden", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{
										name:         "def",
										hidden:       true,
										channelTypes: plugin.AllChannels,
									},
									mockCommand{
										name:         "ghi",
										hidden:       false,
										channelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
					{
						Name: "other",
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{
										name:         "jkl",
										hidden:       false,
										channelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
				}

				p := newProviderFromSources(sources)

				expect := &Module{
					sources: []plugin.SourceModule{
						{
							SourceName: plugin.BuiltInSource,
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									commands: []plugin.Command{
										mockCommand{
											name:         "def",
											hidden:       true,
											channelTypes: plugin.AllChannels,
										},
										mockCommand{
											name:         "ghi",
											hidden:       false,
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
						{
							SourceName: "other",
							Modules: []plugin.Module{
								mockModule{
									name: "abc",
									commands: []plugin.Command{
										mockCommand{
											name:         "jkl",
											hidden:       false,
											channelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
					},
					id:     ".abc",
					hidden: false,
				}

				expect.commands = append(expect.commands, &Command{
					provider: p,
					parent:   expect,
					source: mockCommand{
						name:         "def",
						hidden:       true,
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "def",
									hidden:       true,
									channelTypes: plugin.AllChannels,
								},
								mockCommand{
									name:         "ghi",
									hidden:       false,
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".abc.def",
				})

				expect.commands = append(expect.commands, &Command{
					provider: p,
					parent:   expect,
					source: mockCommand{
						name:         "ghi",
						hidden:       false,
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "def",
									hidden:       true,
									channelTypes: plugin.AllChannels,
								},
								mockCommand{
									name:         "ghi",
									hidden:       false,
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".abc.ghi",
				})

				expect.commands = append(expect.commands, &Command{
					provider: p,
					parent:   expect,
					source: mockCommand{
						name:         "jkl",
						hidden:       false,
						channelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockModule{
							name: "abc",
							commands: []plugin.Command{
								mockCommand{
									name:         "jkl",
									hidden:       false,
									channelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: "other",
					id:         ".abc.jkl",
				})

				fmt.Printf("%#v\n\n", expect)
				fmt.Printf("%#v\n\n", p.Modules()[0])

				assert.Equal(t, expect, p.Modules()[0])
			})
		})
	})
}

func TestPluginProvider_Command(t *testing.T) {
	t.Run("top-level", func(t *testing.T) {
		p := new(PluginProvider)
		p.commands = []plugin.ResolvedCommand{
			&Command{
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source:     mockCommand{name: "abc"},
				id:         ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "def"},
				id:         ".def",
			},
		}

		expect := p.commands[1]
		actual := p.Command(".def")
		assert.Equal(t, expect, actual)
	})

	t.Run("nested", func(t *testing.T) {
		p := new(PluginProvider)
		p.resolver = NewPluginResolver(nil)

		p.sources = []plugin.Source{
			{
				Name: plugin.BuiltInSource,
				Modules: []plugin.Module{
					mockModule{
						name: "abc",
						commands: []plugin.Command{
							mockCommand{name: "def"},
						},
					},
				},
			},
			{
				Name: "another",
				Modules: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
							mockCommand{name: "mno"},
						},
					},
				},
			},
		}

		p.modules = []plugin.ResolvedModule{
			&Module{
				sources: []plugin.SourceModule{
					{
						SourceName: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{name: "def"},
								},
							},
						},
					},
				},
				id: ".abc",
			},
			&Module{
				sources: []plugin.SourceModule{
					{
						SourceName: "another",
						Modules: []plugin.Module{
							mockModule{
								name: "ghi",
								commands: []plugin.Command{
									mockCommand{name: "jkl"},
									mockCommand{name: "mno"},
								},
							},
						},
					},
				},
				id: ".ghi",
			},
		}

		p.modules[0].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     p.modules[0],
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source:     mockCommand{name: "def"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "abc",
						commands: []plugin.Command{
							mockCommand{name: "def"},
						},
					},
				},
				id: ".abc.def",
			},
		}

		p.modules[1].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     p.modules[1],
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "jkl"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
							mockCommand{name: "mno"},
						},
					},
				},
				id: ".ghi.jkl",
			},
			&Command{
				parent:     p.modules[1],
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "mno"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
							mockCommand{name: "mno"},
						},
					},
				},
				id: ".ghi.mno",
			},
		}

		expect := p.modules[1].Commands()[1]
		actual := p.Command(".ghi.mno")
		assert.Equal(t, expect, actual)
	})
}

func TestPluginProvider_Module(t *testing.T) {
	p := new(PluginProvider)
	p.resolver = NewPluginResolver(nil)

	p.sources = []plugin.Source{
		{
			Name: plugin.BuiltInSource,
			Modules: []plugin.Module{
				mockModule{
					name: "abc",
					modules: []plugin.Module{
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{name: "ghi"},
							},
						},
					},
				},
			},
		},
		{
			Name: "another",
			Modules: []plugin.Module{
				mockModule{
					name: "jkl",
					modules: []plugin.Module{
						mockModule{
							name: "mno",
							commands: []plugin.Command{
								mockCommand{name: "pqr"},
							},
						},
					},
				},
			},
		},
	}

	p.modules = []plugin.ResolvedModule{
		&Module{
			sources: []plugin.SourceModule{
				{
					SourceName: plugin.BuiltInSource,
					Modules: []plugin.Module{
						mockModule{
							name: "abc",
							modules: []plugin.Module{
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{name: "ghi"},
									},
								},
							},
						},
					},
				},
			},
			id: ".abc",
		},
		&Module{
			sources: []plugin.SourceModule{
				{
					SourceName: "another",
					Modules: []plugin.Module{
						mockModule{
							name: "jkl",
							modules: []plugin.Module{
								mockModule{
									name: "mno",
									commands: []plugin.Command{
										mockCommand{name: "pqr"},
									},
								},
							},
						},
					},
				},
			},
			id: ".jkl",
		},
	}

	p.modules[0].(*Module).modules = []plugin.ResolvedModule{
		&Module{
			parent: p.modules[0],
			sources: []plugin.SourceModule{
				{
					SourceName: plugin.BuiltInSource,
					Modules: []plugin.Module{
						mockModule{
							name: "abc",
							modules: []plugin.Module{
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{name: "ghi"},
									},
								},
							},
						},
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{name: "ghi"},
							},
						},
					},
				},
			},
			id: ".abc.def",
		},
	}

	p.modules[0].(*Module).modules[0].(*Module).commands = []plugin.ResolvedCommand{
		&Command{
			parent:     p.modules[0].Modules()[0],
			provider:   p,
			sourceName: plugin.BuiltInSource,
			source:     mockCommand{name: "ghi"},
			sourceParents: []plugin.Module{
				mockModule{
					name: "abc",
					modules: []plugin.Module{
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{name: "ghi"},
							},
						},
					},
				},
				mockModule{
					name: "def",
					commands: []plugin.Command{
						mockCommand{name: "ghi"},
					},
				},
			},
			id: ".abc.def.ghi",
		},
	}

	p.modules[1].(*Module).modules = []plugin.ResolvedModule{
		&Module{
			parent: p.modules[1],
			sources: []plugin.SourceModule{
				{
					SourceName: "another",
					Modules: []plugin.Module{
						mockModule{
							name: "jkl",
							modules: []plugin.Module{
								mockModule{
									name: "mno",
									commands: []plugin.Command{
										mockCommand{name: "pqr"},
									},
								},
							},
						},
						mockModule{
							name: "mno",
							commands: []plugin.Command{
								mockCommand{name: "pqr"},
							},
						},
					},
				},
			},
			id: ".jkl.mno",
		},
	}

	p.modules[1].(*Module).modules[0].(*Module).commands = []plugin.ResolvedCommand{
		&Command{
			parent:     p.modules[1].Modules()[0],
			provider:   p,
			sourceName: "another",
			source:     mockCommand{name: "pqr"},
			sourceParents: []plugin.Module{
				mockModule{
					name: "jkl",
					modules: []plugin.Module{
						mockModule{
							name: "mno",
							commands: []plugin.Command{
								mockCommand{name: "pqr"},
							},
						},
					},
				},
				mockModule{
					name: "mno",
					commands: []plugin.Command{
						mockCommand{name: "pqr"},
					},
				},
			},
			id: ".jkl.mno.pqr",
		},
	}

	expect := p.modules[1].Modules()[0]
	actual := p.Module(".jkl.mno")
	assert.Equal(t, expect, actual)
}

func TestPluginProvider_FindCommand(t *testing.T) {
	t.Run("top-level", func(t *testing.T) {
		p := new(PluginProvider)
		p.commands = []plugin.ResolvedCommand{
			&Command{
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source:     mockCommand{name: "abc"},
				id:         ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "def"},
				id:         ".def",
			},
		}

		expect := p.commands[1]
		actual := p.FindCommand(" def  \n")
		assert.Equal(t, expect, actual)
	})

	t.Run("nested", func(t *testing.T) {
		p := new(PluginProvider)
		p.resolver = NewPluginResolver(nil)

		p.sources = []plugin.Source{
			{
				Name: plugin.BuiltInSource,
				Modules: []plugin.Module{
					mockModule{
						name: "abc",
						commands: []plugin.Command{
							mockCommand{name: "def"},
						},
					},
				},
			},
			{
				Name: "another",
				Modules: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
							mockCommand{name: "mno"},
						},
					},
				},
			},
		}

		p.modules = []plugin.ResolvedModule{
			&Module{
				parent: nil,
				sources: []plugin.SourceModule{
					{
						SourceName: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockModule{
								name: "abc",
								commands: []plugin.Command{
									mockCommand{name: "def"},
								},
							},
						},
					},
				},
				id: ".abc",
			},
			&Module{
				parent: nil,
				sources: []plugin.SourceModule{
					{
						SourceName: "another",
						Modules: []plugin.Module{
							mockModule{
								name: "ghi",
								commands: []plugin.Command{
									mockCommand{name: "jkl"},
									mockCommand{name: "mno"},
								},
							},
						},
					},
				},
				id: ".ghi",
			},
		}

		p.modules[0].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     p.modules[0],
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source:     mockCommand{name: "def"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "abc",
						commands: []plugin.Command{
							mockCommand{name: "def"},
						},
					},
				},
				id: ".abc.def",
			},
		}

		p.modules[1].(*Module).commands = []plugin.ResolvedCommand{
			&Command{
				parent:     p.modules[1],
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "jkl"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
							mockCommand{name: "mno"},
						},
					},
				},
				id: ".ghi.jkl",
			},
			&Command{
				parent:     p.modules[1],
				provider:   p,
				sourceName: "another",
				source:     mockCommand{name: "mno"},
				sourceParents: []plugin.Module{
					mockModule{
						name: "ghi",
						commands: []plugin.Command{
							mockCommand{name: "jkl"},
							mockCommand{name: "mno"},
						},
					},
				},
				id: ".ghi.mno",
			},
		}

		expect := p.modules[1].Commands()[1]
		actual := p.FindCommand(" ghi  mno")
		assert.Equal(t, expect, actual)
	})
}

func TestPluginProvider_FindModule(t *testing.T) {
	p := new(PluginProvider)
	p.resolver = NewPluginResolver(nil)

	p.sources = []plugin.Source{
		{
			Name: plugin.BuiltInSource,
			Modules: []plugin.Module{
				mockModule{
					name: "abc",
					modules: []plugin.Module{
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{name: "ghi"},
							},
						},
					},
				},
			},
		},
		{
			Name: "another",
			Modules: []plugin.Module{
				mockModule{
					name: "jkl",
					modules: []plugin.Module{
						mockModule{
							name: "mno",
							commands: []plugin.Command{
								mockCommand{name: "pqr"},
							},
						},
					},
				},
			},
		},
	}

	p.modules = []plugin.ResolvedModule{
		&Module{
			sources: []plugin.SourceModule{
				{
					SourceName: plugin.BuiltInSource,
					Modules: []plugin.Module{
						mockModule{
							name: "abc",
							modules: []plugin.Module{
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{name: "ghi"},
									},
								},
							},
						},
					},
				},
			},
			id: ".abc",
		},
		&Module{
			sources: []plugin.SourceModule{
				{
					SourceName: "another",
					Modules: []plugin.Module{
						mockModule{
							name: "jkl",
							modules: []plugin.Module{
								mockModule{
									name: "mno",
									commands: []plugin.Command{
										mockCommand{name: "pqr"},
									},
								},
							},
						},
					},
				},
			},
			id: ".jkl",
		},
	}

	p.modules[0].(*Module).modules = []plugin.ResolvedModule{
		&Module{
			parent: p.modules[0],
			sources: []plugin.SourceModule{
				{
					SourceName: plugin.BuiltInSource,
					Modules: []plugin.Module{
						mockModule{
							name: "abc",
							modules: []plugin.Module{
								mockModule{
									name: "def",
									commands: []plugin.Command{
										mockCommand{name: "ghi"},
									},
								},
							},
						},
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{name: "ghi"},
							},
						},
					},
				},
			},
			id: ".abc.def",
		},
	}

	p.modules[0].(*Module).modules[0].(*Module).commands = []plugin.ResolvedCommand{
		&Command{
			parent:     p.modules[0].Modules()[0],
			provider:   p,
			sourceName: plugin.BuiltInSource,
			source:     mockCommand{name: "ghi"},
			sourceParents: []plugin.Module{
				mockModule{
					name: "abc",
					modules: []plugin.Module{
						mockModule{
							name: "def",
							commands: []plugin.Command{
								mockCommand{name: "ghi"},
							},
						},
					},
				},
				mockModule{
					name: "def",
					commands: []plugin.Command{
						mockCommand{name: "ghi"},
					},
				},
			},
			id: ".abc.def.ghi",
		},
	}

	p.modules[1].(*Module).modules = []plugin.ResolvedModule{
		&Module{
			parent: p.modules[1],
			sources: []plugin.SourceModule{
				{
					SourceName: "another",
					Modules: []plugin.Module{
						mockModule{
							name: "jkl",
							modules: []plugin.Module{
								mockModule{
									name: "mno",
									commands: []plugin.Command{
										mockCommand{name: "pqr"},
									},
								},
							},
						},
						mockModule{
							name: "mno",
							commands: []plugin.Command{
								mockCommand{name: "pqr"},
							},
						},
					},
				},
			},
			id: ".jkl.mno",
		},
	}

	p.modules[1].(*Module).modules[0].(*Module).commands = []plugin.ResolvedCommand{
		&Command{
			parent:     p.modules[1].Modules()[0],
			provider:   p,
			sourceName: "another",
			source:     mockCommand{name: "pqr"},
			sourceParents: []plugin.Module{
				mockModule{
					name: "jkl",
					modules: []plugin.Module{
						mockModule{
							name: "mno",
							commands: []plugin.Command{
								mockCommand{name: "pqr"},
							},
						},
					},
				},
				mockModule{
					name: "mno",
					commands: []plugin.Command{
						mockCommand{name: "pqr"},
					},
				},
			},
			id: ".jkl.mno.pqr",
		},
	}

	expect := p.modules[1].Modules()[0]
	actual := p.FindModule("jkl \nmno")
	assert.Equal(t, expect, actual)
}

func TestPluginProvider_UnavailablePluginSources(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		r := NewPluginResolver(nil)
		r.AddSource("another",
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return nil, nil, errors.New("abc")
			})

		p := &PluginProvider{
			resolver: r,
			commands: make([]plugin.ResolvedCommand, 0),
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockCommand{
							name:         "abc",
							channelTypes: plugin.GuildChannels,
						},
					},
				},
			},
		}

		expect := []plugin.UnavailablePluginSource{
			{
				Name:  "another",
				Error: errors.New("abc"),
			},
		}

		actual := p.UnavailablePluginSources()
		assert.Equal(t, expect, actual)
	})

	t.Run("loaded", func(t *testing.T) {
		p := &PluginProvider{
			resolver: NewPluginResolver(nil),
			commands: make([]plugin.ResolvedCommand, 0),
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockCommand{
							name:         "abc",
							channelTypes: plugin.GuildChannels,
						},
					},
				},
			},
			unavailableSources: []plugin.UnavailablePluginSource{
				{
					Name:  "another",
					Error: errors.New("abc"),
				},
			},
		}

		expect := p.unavailableSources

		actual := p.UnavailablePluginSources()
		assert.Equal(t, expect, actual)
	})
}
