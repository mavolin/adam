package resolved

import (
	"errors"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestPluginProvider_PluginSources(t *testing.T) {
	t.Run("not loaded", func(t *testing.T) {
		r := NewPluginResolver(nil)
		r.AddSource("another",
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return []plugin.Command{mockplugin.Command{Name: "def"}}, []plugin.Module{mockplugin.Module{Name: "ghi"}}, nil
			})

		p := &PluginProvider{
			commands: make([]plugin.ResolvedCommand, 0),
			resolver: r,
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockplugin.Command{Name: "abc"},
					},
				},
			},
			usedNames: make(map[string]struct{}),
		}

		expect := []plugin.Source{
			{
				Name: plugin.BuiltInSource,
				Commands: []plugin.Command{
					mockplugin.Command{Name: "abc"},
				},
			},
			{
				Name: "another",
				Commands: []plugin.Command{
					mockplugin.Command{Name: "def"},
				},
				Modules: []plugin.Module{
					mockplugin.Module{Name: "ghi"},
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
						mockplugin.Command{Name: "abc"},
					},
				},
				{
					Name: "another",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "def"},
					},
					Modules: []plugin.Module{
						mockplugin.Module{Name: "ghi"},
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
					mockplugin.Command{
						Name:         "def",
						ChannelTypes: plugin.GuildChannels,
					},
				}, nil, nil
			})

		p := &PluginProvider{
			resolver: r,
			sources: []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockplugin.Command{
							Name:         "abc",
							ChannelTypes: plugin.GuildChannels,
						},
					},
				},
			},
			usedNames: make(map[string]struct{}),
		}

		p.commands = append(p.commands, &Command{
			provider:   p,
			sourceName: plugin.BuiltInSource,
			source: mockplugin.Command{
				Name:         "abc",
				ChannelTypes: plugin.GuildChannels,
			},
			id: ".abc",
		})

		expect := []plugin.ResolvedCommand{
			&Command{
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source: mockplugin.Command{
					Name:         "abc",
					ChannelTypes: plugin.GuildChannels,
				},
				id: ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source: mockplugin.Command{
					Name:         "def",
					ChannelTypes: plugin.GuildChannels,
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
				Commands: []plugin.Command{mockplugin.Command{Name: "abc"}},
			},
			{
				Name:     "another",
				Commands: []plugin.Command{mockplugin.Command{Name: "def"}},
			},
		}
		p.commands = []plugin.ResolvedCommand{
			&Command{
				provider:   p,
				sourceName: plugin.BuiltInSource,
				source:     mockplugin.Command{Name: "abc"},
				id:         ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source:     mockplugin.Command{Name: "def"},
				id:         ".def",
			},
		}

		expect := p.commands

		actual := p.Commands()
		assert.Equal(t, expect, actual)
	})

	t.Run("provider", func(t *testing.T) {
		t.Run("single", func(t *testing.T) {
			throttler := mockplugin.NewThrottler(errors.New("abc"))

			sources := []plugin.Source{
				{
					Name: plugin.BuiltInSource,
					Commands: []plugin.Command{
						mockplugin.Command{
							Name:         "abc",
							Hidden:       true,
							Throttler:    throttler,
							ChannelTypes: plugin.AllChannels,
						},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: plugin.BuiltInSource,
					source: mockplugin.Command{
						Name:         "abc",
						Hidden:       true,
						Throttler:    throttler,
						ChannelTypes: plugin.AllChannels,
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
						mockplugin.Command{Name: "def", ChannelTypes: plugin.AllChannels},
					},
				},
				{
					Name: "ghi",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "def", ChannelTypes: plugin.AllChannels},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source: mockplugin.Command{
						Name:         "def",
						ChannelTypes: plugin.AllChannels,
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
						mockplugin.Command{Name: "jkl", ChannelTypes: plugin.AllChannels},
					},
				},
				{
					Name: "abc",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "def", ChannelTypes: plugin.AllChannels},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source:     mockplugin.Command{Name: "def", ChannelTypes: plugin.AllChannels},
					id:         ".def",
				},
				&Command{
					provider:   p,
					sourceName: "ghi",
					source:     mockplugin.Command{Name: "jkl", ChannelTypes: plugin.AllChannels},
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
						mockplugin.Command{Name: "def", ChannelTypes: plugin.AllChannels},
					},
				},
				{
					Name: "ghi",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "def", ChannelTypes: plugin.AllChannels}, // duplicate
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source:     mockplugin.Command{Name: "def", ChannelTypes: plugin.AllChannels},
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
						mockplugin.Command{
							Name:         "def",
							Aliases:      []string{"ghi", "jkl"},
							ChannelTypes: plugin.AllChannels,
						},
					},
				},
				{
					Name: "mno",
					Commands: []plugin.Command{
						mockplugin.Command{
							Name:         "pqr",
							Aliases:      []string{"jkl", "stu"}, // duplicate alias
							ChannelTypes: plugin.AllChannels,
						},
					},
				},
			}

			p := newProviderFromSources(sources)

			expect := []plugin.ResolvedCommand{
				&Command{
					provider:   p,
					sourceName: "abc",
					source: mockplugin.Command{
						Name:         "def",
						Aliases:      []string{"ghi", "jkl"},
						ChannelTypes: plugin.AllChannels,
					},
					id:      ".def",
					aliases: []string{"ghi", "jkl"},
				},
				&Command{
					provider:   p,
					sourceName: "mno",
					source: mockplugin.Command{
						Name:         "pqr",
						Aliases:      []string{"jkl", "stu"},
						ChannelTypes: plugin.AllChannels,
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
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{
								Name:         "jkl",
								ChannelTypes: plugin.GuildChannels,
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
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "def",
									ChannelTypes: plugin.GuildChannels,
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
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{
										Name:         "def",
										ChannelTypes: plugin.GuildChannels,
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
				source: mockplugin.Command{
					Name:         "def",
					ChannelTypes: plugin.GuildChannels,
				},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "abc",
						Commands: []plugin.Command{
							mockplugin.Command{
								Name:         "def",
								ChannelTypes: plugin.GuildChannels,
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
						mockplugin.Module{
							Name: "ghi",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "jkl",
									ChannelTypes: plugin.GuildChannels,
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
				source: mockplugin.Command{
					Name:         "jkl",
					ChannelTypes: plugin.GuildChannels,
				},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{
								Name:         "jkl",
								ChannelTypes: plugin.GuildChannels,
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
					mockplugin.Module{
						Name: "abc",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "def"},
						},
					},
				},
			},
			{
				Name: "another",
				Modules: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
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
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "def"},
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
							mockplugin.Module{
								Name: "ghi",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "jkl"},
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
				source:     mockplugin.Command{Name: "def"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "abc",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "def"},
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
				source:     mockplugin.Command{Name: "jkl"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
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

	t.Run("provider", func(t *testing.T) {
		t.Run("generate", func(t *testing.T) {
			t.Run("", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "zyx", ChannelTypes: plugin.AllChannels},
								},
							},
							mockplugin.Module{
								Name: "def",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "wvu", ChannelTypes: plugin.AllChannels},
								},
							},
						},
					},
					{
						Name: "custom_commands",
						Modules: []plugin.Module{
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "zyx", ChannelTypes: plugin.AllChannels},
								},
							},
							mockplugin.Module{
								Name: "def",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "tsr", ChannelTypes: plugin.AllChannels},
								},
							},
							mockplugin.Module{
								Name: "ghi",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "qpo", ChannelTypes: plugin.AllChannels},
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
								mockplugin.Module{
									Name: "abc",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "zyx", ChannelTypes: plugin.AllChannels},
									},
								},
							},
						},
						{
							SourceName: "custom_commands",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "abc",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "zyx", ChannelTypes: plugin.AllChannels},
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
					source: mockplugin.Command{
						Name:         "zyx",
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "zyx",
									ChannelTypes: plugin.AllChannels,
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
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "wvu",
											ChannelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
						{
							SourceName: "custom_commands",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "tsr",
											ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "tsr",
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "tsr",
									ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "wvu",
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "wvu",
									ChannelTypes: plugin.AllChannels,
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
								mockplugin.Module{
									Name: "ghi",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "qpo",
											ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "qpo",
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "ghi",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "qpo",
									ChannelTypes: plugin.AllChannels,
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
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{
										Name:         "def",
										ChannelTypes: plugin.AllChannels,
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
								mockplugin.Module{
									Name: "abc",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "def",
											ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "def",
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "def",
									ChannelTypes: plugin.AllChannels,
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
							mockplugin.Module{
								Name: "abc",
								Modules: []plugin.Module{
									mockplugin.Module{
										Name: "def",
										Commands: []plugin.Command{
											mockplugin.Command{
												Name:         "ghi",
												ChannelTypes: plugin.AllChannels,
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
								mockplugin.Module{
									Name: "abc",
									Modules: []plugin.Module{
										mockplugin.Module{
											Name: "def",
											Commands: []plugin.Command{
												mockplugin.Command{
													Name:         "ghi",
													ChannelTypes: plugin.AllChannels,
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
								mockplugin.Module{
									Name: "abc",
									Modules: []plugin.Module{
										mockplugin.Module{
											Name: "def",
											Commands: []plugin.Command{
												mockplugin.Command{
													Name:         "ghi",
													ChannelTypes: plugin.AllChannels,
												},
											},
										},
									},
								},
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "ghi",
											ChannelTypes: plugin.AllChannels,
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
					source:   mockplugin.Command{Name: "ghi", ChannelTypes: plugin.AllChannels},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "ghi", ChannelTypes: plugin.AllChannels},
									},
								},
							},
						},
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "ghi",
									ChannelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: plugin.BuiltInSource,
					id:         ".abc.def.ghi",
				})

				assert.Equal(t, expect, p.Command(".abc.def.ghi").Parent())
			})

			t.Run("children Hidden", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{
										Name:         "def",
										Hidden:       true,
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
					{
						Name: "other",
						Modules: []plugin.Module{
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{
										Name:         "ghi",
										Hidden:       true,
										ChannelTypes: plugin.AllChannels,
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
								mockplugin.Module{
									Name: "abc",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "def",
											Hidden:       true,
											ChannelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
						{
							SourceName: "other",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "abc",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "ghi",
											Hidden:       true,
											ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "def",
						Hidden:       true,
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "def",
									Hidden:       true,
									ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "ghi",
						Hidden:       true,
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "ghi",
									Hidden:       true,
									ChannelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: "other",
					id:         ".abc.ghi",
				})

				assert.Equal(t, expect, p.Modules()[0])
			})

			t.Run("not Hidden", func(t *testing.T) {
				sources := []plugin.Source{
					{
						Name: plugin.BuiltInSource,
						Modules: []plugin.Module{
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{
										Name:         "def",
										Hidden:       true,
										ChannelTypes: plugin.AllChannels,
									},
									mockplugin.Command{
										Name:         "ghi",
										Hidden:       false,
										ChannelTypes: plugin.AllChannels,
									},
								},
							},
						},
					},
					{
						Name: "other",
						Modules: []plugin.Module{
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{
										Name:         "jkl",
										Hidden:       false,
										ChannelTypes: plugin.AllChannels,
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
								mockplugin.Module{
									Name: "abc",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "def",
											Hidden:       true,
											ChannelTypes: plugin.AllChannels,
										},
										mockplugin.Command{
											Name:         "ghi",
											Hidden:       false,
											ChannelTypes: plugin.AllChannels,
										},
									},
								},
							},
						},
						{
							SourceName: "other",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "abc",
									Commands: []plugin.Command{
										mockplugin.Command{
											Name:         "jkl",
											Hidden:       false,
											ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "def",
						Hidden:       true,
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "def",
									Hidden:       true,
									ChannelTypes: plugin.AllChannels,
								},
								mockplugin.Command{
									Name:         "ghi",
									Hidden:       false,
									ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "ghi",
						Hidden:       false,
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "def",
									Hidden:       true,
									ChannelTypes: plugin.AllChannels,
								},
								mockplugin.Command{
									Name:         "ghi",
									Hidden:       false,
									ChannelTypes: plugin.AllChannels,
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
					source: mockplugin.Command{
						Name:         "jkl",
						Hidden:       false,
						ChannelTypes: plugin.AllChannels,
					},
					sourceParents: []plugin.Module{
						mockplugin.Module{
							Name: "abc",
							Commands: []plugin.Command{
								mockplugin.Command{
									Name:         "jkl",
									Hidden:       false,
									ChannelTypes: plugin.AllChannels,
								},
							},
						},
					},
					sourceName: "other",
					id:         ".abc.jkl",
				})

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
				source:     mockplugin.Command{Name: "abc"},
				id:         ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source:     mockplugin.Command{Name: "def"},
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
					mockplugin.Module{
						Name: "abc",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "def"},
						},
					},
				},
			},
			{
				Name: "another",
				Modules: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
							mockplugin.Command{Name: "mno"},
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
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "def"},
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
							mockplugin.Module{
								Name: "ghi",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "jkl"},
									mockplugin.Command{Name: "mno"},
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
				source:     mockplugin.Command{Name: "def"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "abc",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "def"},
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
				source:     mockplugin.Command{Name: "jkl"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
							mockplugin.Command{Name: "mno"},
						},
					},
				},
				id: ".ghi.jkl",
			},
			&Command{
				parent:     p.modules[1],
				provider:   p,
				sourceName: "another",
				source:     mockplugin.Command{Name: "mno"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
							mockplugin.Command{Name: "mno"},
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
				mockplugin.Module{
					Name: "abc",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "ghi"},
							},
						},
					},
				},
			},
		},
		{
			Name: "another",
			Modules: []plugin.Module{
				mockplugin.Module{
					Name: "jkl",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "mno",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "pqr"},
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
						mockplugin.Module{
							Name: "abc",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "ghi"},
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
						mockplugin.Module{
							Name: "jkl",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "mno",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "pqr"},
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
						mockplugin.Module{
							Name: "abc",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "ghi"},
									},
								},
							},
						},
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "ghi"},
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
			source:     mockplugin.Command{Name: "ghi"},
			sourceParents: []plugin.Module{
				mockplugin.Module{
					Name: "abc",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "ghi"},
							},
						},
					},
				},
				mockplugin.Module{
					Name: "def",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "ghi"},
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
						mockplugin.Module{
							Name: "jkl",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "mno",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "pqr"},
									},
								},
							},
						},
						mockplugin.Module{
							Name: "mno",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "pqr"},
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
			source:     mockplugin.Command{Name: "pqr"},
			sourceParents: []plugin.Module{
				mockplugin.Module{
					Name: "jkl",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "mno",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "pqr"},
							},
						},
					},
				},
				mockplugin.Module{
					Name: "mno",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "pqr"},
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
				source:     mockplugin.Command{Name: "abc"},
				id:         ".abc",
			},
			&Command{
				provider:   p,
				sourceName: "another",
				source:     mockplugin.Command{Name: "def"},
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
					mockplugin.Module{
						Name: "abc",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "def"},
						},
					},
				},
			},
			{
				Name: "another",
				Modules: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
							mockplugin.Command{Name: "mno"},
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
							mockplugin.Module{
								Name: "abc",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "def"},
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
							mockplugin.Module{
								Name: "ghi",
								Commands: []plugin.Command{
									mockplugin.Command{Name: "jkl"},
									mockplugin.Command{Name: "mno"},
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
				source:     mockplugin.Command{Name: "def"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "abc",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "def"},
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
				source:     mockplugin.Command{Name: "jkl"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
							mockplugin.Command{Name: "mno"},
						},
					},
				},
				id: ".ghi.jkl",
			},
			&Command{
				parent:     p.modules[1],
				provider:   p,
				sourceName: "another",
				source:     mockplugin.Command{Name: "mno"},
				sourceParents: []plugin.Module{
					mockplugin.Module{
						Name: "ghi",
						Commands: []plugin.Command{
							mockplugin.Command{Name: "jkl"},
							mockplugin.Command{Name: "mno"},
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
				mockplugin.Module{
					Name: "abc",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "ghi"},
							},
						},
					},
				},
			},
		},
		{
			Name: "another",
			Modules: []plugin.Module{
				mockplugin.Module{
					Name: "jkl",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "mno",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "pqr"},
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
						mockplugin.Module{
							Name: "abc",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "ghi"},
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
						mockplugin.Module{
							Name: "jkl",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "mno",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "pqr"},
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
						mockplugin.Module{
							Name: "abc",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "def",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "ghi"},
									},
								},
							},
						},
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "ghi"},
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
			source:     mockplugin.Command{Name: "ghi"},
			sourceParents: []plugin.Module{
				mockplugin.Module{
					Name: "abc",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "def",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "ghi"},
							},
						},
					},
				},
				mockplugin.Module{
					Name: "def",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "ghi"},
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
						mockplugin.Module{
							Name: "jkl",
							Modules: []plugin.Module{
								mockplugin.Module{
									Name: "mno",
									Commands: []plugin.Command{
										mockplugin.Command{Name: "pqr"},
									},
								},
							},
						},
						mockplugin.Module{
							Name: "mno",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "pqr"},
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
			source:     mockplugin.Command{Name: "pqr"},
			sourceParents: []plugin.Module{
				mockplugin.Module{
					Name: "jkl",
					Modules: []plugin.Module{
						mockplugin.Module{
							Name: "mno",
							Commands: []plugin.Command{
								mockplugin.Command{Name: "pqr"},
							},
						},
					},
				},
				mockplugin.Module{
					Name: "mno",
					Commands: []plugin.Command{
						mockplugin.Command{Name: "pqr"},
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
						mockplugin.Command{
							Name:         "abc",
							ChannelTypes: plugin.GuildChannels,
						},
					},
				},
			},
		}

		expect := []plugin.UnavailableSource{
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
						mockplugin.Command{
							Name:         "abc",
							ChannelTypes: plugin.GuildChannels,
						},
					},
				},
			},
			unavailableSources: []plugin.UnavailableSource{
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
