package messageutil

import (
	"context"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestWaiter_Await(t *testing.T) {
	t.Run("initial timeout", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Base: state.NewBase(),
			Message: discord.Message{
				ChannelID: 123,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			},
			Localizer: i18n.NewFallbackLocalizer(),
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
			Replier: replier.WrapState(s, false),
		}

		expect := &TimeoutError{UserID: ctx.Author.ID}

		msg, actual := NewReplyWaiter(s, ctx).
			Await(1, 1)
		assert.Nil(t, msg)
		assert.Equal(t, expect, actual)
	})
}

func TestWaiter_handleMessages(t *testing.T) {
	testError := errors.New("abc")

	testCases := []struct {
		name   string
		waiter *ReplyWaiter
		e      *state.MessageCreateEvent
		expect interface{}
	}{
		{
			name: "middleware block",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: state.NewBase()}).
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
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base:    state.NewBase(),
				Message: discord.Message{ChannelID: 123},
			}),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{ChannelID: 321},
				},
			},
			expect: nil,
		},
		{
			name: "author not matching",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base: state.NewBase(),
				Message: discord.Message{
					Author: discord.User{ID: 123},
				},
			}),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Author: discord.User{ID: 321}},
				},
			},
			expect: nil,
		},
		{
			name: "canceled - case sensitive",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: state.NewBase()}).
				WithCancelKeywords("aBc").
				CaseSensitive(),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "aBc"},
				},
			},
			expect: errors.Abort,
		},
		{
			name: "not canceled - case sensitive",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: state.NewBase()}).
				WithCancelKeywords("aBc").
				CaseSensitive(),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "AbC"},
				},
			},
			expect: &discord.Message{Content: "AbC"},
		},
		{
			name: "canceled - case insensitive",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: state.NewBase()}).
				WithCancelKeywords("aBc"),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "AbC"},
				},
			},
			expect: errors.Abort,
		},
		{
			name:   "success",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: state.NewBase()}),
			e: &state.MessageCreateEvent{
				Base: state.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "abc"},
				},
			},
			expect: &discord.Message{Content: "abc"},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := state.NewMocker(t)
			defer m.Eval()

			c.waiter.state = s

			var result chan interface{}
			// cause a nil pointer dereference, if something gets sent anyway
			// although c.expect == nil
			if c.expect != nil {
				result = make(chan interface{})
			}

			rm := c.waiter.handleMessages(context.TODO(), result)

			s.Call(c.e)

			if c.expect != nil {
				var actual interface{}

				//goland:noinspection GoNilness
				select {
				case actual = <-result:
				case <-time.After(2 * time.Second):
					require.Fail(t, "Function timed out")
				}

				assert.Equal(t, c.expect, actual)
			}

			rm()
		})
	}
}

func TestWaiter_handleCancelReactions(t *testing.T) {
	testCases := []struct {
		name   string
		waiter *ReplyWaiter
		e      *state.MessageReactionAddEvent
		expect interface{}
	}{
		{
			name: "message id not matching",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base:    state.NewBase(),
				Message: discord.Message{ChannelID: 456},
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
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base:    state.NewBase(),
				Message: discord.Message{ChannelID: 456},
			}),
			e: &state.MessageReactionAddEvent{
				Base: state.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					MessageID: 123,
					Emoji:     discord.Emoji{Name: "🍑"},
				},
			},
			expect: nil,
		},
		{
			name: "user id not matching",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base: state.NewBase(),
				Message: discord.Message{
					ChannelID: 456,
					Author:    discord.User{ID: 789},
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
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base: state.NewBase(),
				Message: discord.Message{
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				},
			}),
			e: &state.MessageReactionAddEvent{
				Base: state.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					MessageID: 123,
					UserID:    789,
					Emoji:     discord.Emoji{Name: "🍆"},
				},
			},
			expect: errors.Abort,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := state.NewMocker(t)
			defer m.Eval()

			c.waiter.state = s

			var reactionMessageID discord.MessageID = 123
			var reaction discord.APIEmoji = "🍆"

			c.waiter.
				WithCancelReactions(reactionMessageID, reaction)

			m.React(c.waiter.channelID, reactionMessageID, reaction)

			var result chan interface{}
			// cause a nil pointer dereference, if something gets sent anyway
			// although c.expect == nil
			if c.expect != nil {
				result = make(chan interface{})
			}

			c.waiter.handleCancelReactions(context.TODO(), result)
			s.Call(c.e)

			if c.expect != nil {
				var actual interface{}

				//goland:noinspection GoNilness
				select {
				case actual = <-result:
				case <-time.After(2 * time.Second):
					require.Fail(t, "Function timed out")
				}

				assert.Equal(t, c.expect, actual)
			}
		})
	}
}
