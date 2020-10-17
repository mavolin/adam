package reply

import (
	"context"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestWaiter_Await(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						ChannelID: 123,
						GuildID:   456,
						Author: discord.User{
							ID: 789,
						},
					},
				},
			},
			Localizer: mock.NoOpLocalizer,
			DiscordDataProvider: mock.DiscordDataProvider{
				ChannelReturn: &discord.Channel{},
				ChannelError:  nil,
				GuildReturn: &discord.Guild{
					Roles: []discord.Role{
						{
							ID:          012,
							Permissions: discord.PermissionSendMessages,
						},
					},
				},
				GuildError: nil,
				SelfReturn: &discord.Member{
					RoleIDs: []discord.RoleID{012},
				},
			},
			Replier: replier.WrapState(s, 0, 123),
		}

		expect := errors.NewUserInfol(timeoutInfo.
			WithPlaceholders(&timeoutInfoPlaceholders{
				ResponseUserMention: ctx.Author.Mention(),
			}))

		msg, actual := NewWaiter(s, ctx).
			Await(1)
		assert.Nil(t, msg)
		assert.Equal(t, expect, actual)

		m.Eval()
	})
}

func TestWaiter_handleMessages(t *testing.T) {
	testError := errors.New("abc")

	testCases := []struct {
		name   string
		waiter *Waiter
		e      *state.MessageCreateEvent
		expect interface{}
	}{
		{
			name: "middleware block",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: new(gateway.MessageCreateEvent),
				},
			}).
				WithMiddlewares(func(*state.State, *state.MessageCreateEvent) error {
					return testError
				}),
			e: &state.MessageCreateEvent{
				Base:               state.NewBase(),
				MessageCreateEvent: new(gateway.MessageCreateEvent),
			},
			expect: nil,
		},
		{
			name: "channel not matching",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 123,
						},
					},
				},
			}),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						ChannelID: 321,
					},
				},
			},
			expect: nil,
		},
		{
			name: "author not matching",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							Author: discord.User{
								ID: 123,
							},
						},
					},
				},
			}),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Author: discord.User{
							ID: 321,
						},
					},
				},
			},
			expect: nil,
		},
		{
			name: "canceled - case sensitive",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: new(gateway.MessageCreateEvent),
				},
			}).
				WithCancelKeyword("aBc").
				CaseSensitive(),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Content: "aBc",
					},
				},
			},
			expect: Canceled,
		},
		{
			name: "not canceled - case sensitive",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: new(gateway.MessageCreateEvent),
				},
			}).
				WithCancelKeyword("aBc").
				CaseSensitive(),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Content: "AbC",
					},
				},
			},
			expect: &discord.Message{
				Content: "AbC",
			},
		},
		{
			name: "canceled - case insensitive",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: new(gateway.MessageCreateEvent),
				},
			}).
				WithCancelKeyword("aBc"),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Content: "AbC",
					},
				},
			},
			expect: Canceled,
		},
		{
			name: "success",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: new(gateway.MessageCreateEvent),
				},
			}),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Content: "abc",
					},
				},
			},
			expect: &discord.Message{
				Content: "abc",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := state.NewMocker(t)
			c.waiter.state = s

			var result chan interface{}
			// cause a nil pointer dereference, if something gets sent anyway
			// although c.expect == nil
			if c.expect != nil {
				result = make(chan interface{})
			}

			rm, err := c.waiter.handleMessages(context.TODO(), result)
			assert.NoError(t, err)

			s.Call(c.e)

			if c.expect != nil {
				var actual interface{}

				select {
				case actual = <-result:
				case <-time.After(2 * time.Second):
					require.Fail(t, "Function timed out")
				}

				assert.Equal(t, c.expect, actual)
			}

			m.Eval()

			rm()
		})
	}
}

func TestWaiter_handleCancelReactions(t *testing.T) {
	testCases := []struct {
		name   string
		waiter *Waiter
		e      *state.MessageReactionAddEvent
		expect interface{}
	}{
		{
			name: "message id not matching",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 456,
						},
					},
				},
			}),
			e: &state.MessageReactionAddEvent{
				Base: state.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					MessageID: 321,
				},
			},
			expect: nil,
		},
		{
			name: "emoji not matching",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 456,
						},
					},
				},
			}),
			e: &state.MessageReactionAddEvent{
				Base: state.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					MessageID: 123,
					Emoji: discord.Emoji{
						Name: "🍑",
					},
				},
			},
			expect: nil,
		},
		{
			name: "user id not matching",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 456,
							Author: discord.User{
								ID: 789,
							},
						},
					},
				},
			}),
			e: &state.MessageReactionAddEvent{
				Base: state.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					UserID: 987,
				},
			},
			expect: nil,
		},
		{
			name: "success",
			waiter: NewWaiter(nil, &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							ChannelID: 456,
							Author: discord.User{
								ID: 789,
							},
						},
					},
				},
			}),
			e: &state.MessageReactionAddEvent{
				Base: state.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					MessageID: 123,
					UserID:    789,
					Emoji: discord.Emoji{
						Name: "🍆",
					},
				},
			},
			expect: Canceled,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := state.NewMocker(t)
			c.waiter.state = s

			var reactionMessageID discord.MessageID = 123
			reaction := "🍆"

			c.waiter.
				WithCancelReaction(reactionMessageID, reaction)

			m.React(c.waiter.ctx.ChannelID, reactionMessageID, reaction)

			var result chan interface{}
			// cause a nil pointer dereference, if something gets sent anyway
			// although c.expect == nil
			if c.expect != nil {
				result = make(chan interface{})
			}

			_, err := c.waiter.handleCancelReactions(context.TODO(), result)
			assert.NoError(t, err)

			s.Call(c.e)

			if c.expect != nil {
				var actual interface{}

				select {
				case actual = <-result:
				case <-time.After(2 * time.Second):
					require.Fail(t, "Function timed out")
				}

				assert.Equal(t, c.expect, actual)
			}

			m.Eval()
		})
	}
}
