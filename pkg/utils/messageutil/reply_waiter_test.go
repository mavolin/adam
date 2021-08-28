package messageutil

import (
	"context"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestWaiter_Await(t *testing.T) {
	t.Parallel()

	t.Run("initial timeout", func(t *testing.T) {
		t.Parallel()

		_, s := state.NewMocker(t)

		ctx := &plugin.Context{
			Base: event.NewBase(),
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
						{ID: 12, Permissions: discord.PermissionSendMessages},
					},
				},
				GuildError: nil,
				SelfReturn: &discord.Member{RoleIDs: []discord.RoleID{12}},
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
	t.Parallel()

	testError := errors.New("abc")

	testCases := []struct {
		name   string
		waiter *ReplyWaiter
		e      *event.MessageCreate
		expect interface{}
	}{
		{
			name: "middleware block",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: event.NewBase()}).
				WithMiddlewares(func(*state.State, *event.MessageCreate) error {
					return testError
				}),
			e: &event.MessageCreate{
				Base:               event.NewBase(),
				MessageCreateEvent: new(gateway.MessageCreateEvent),
			},
			expect: nil,
		},
		{
			name: "channel not matching",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base:    event.NewBase(),
				Message: discord.Message{ChannelID: 123},
			}),
			e: &event.MessageCreate{
				Base: event.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{ChannelID: 321},
				},
			},
			expect: nil,
		},
		{
			name: "author not matching",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base: event.NewBase(),
				Message: discord.Message{
					Author: discord.User{ID: 123},
				},
			}),
			e: &event.MessageCreate{
				Base: event.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Author: discord.User{ID: 321}},
				},
			},
			expect: nil,
		},
		{
			name: "canceled - case sensitive",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: event.NewBase()}).
				WithCancelKeywords("aBc").
				CaseSensitive(),
			e: &event.MessageCreate{
				Base: event.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "aBc"},
				},
			},
			expect: errors.Abort,
		},
		{
			name: "not canceled - case sensitive",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: event.NewBase()}).
				WithCancelKeywords("aBc").
				CaseSensitive(),
			e: &event.MessageCreate{
				Base: event.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "AbC"},
				},
			},
			expect: &discord.Message{Content: "AbC"},
		},
		{
			name: "canceled - case insensitive",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: event.NewBase()}).
				WithCancelKeywords("aBc"),
			e: &event.MessageCreate{
				Base: event.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "AbC"},
				},
			},
			expect: errors.Abort,
		},
		{
			name:   "success",
			waiter: NewReplyWaiter(nil, &plugin.Context{Base: event.NewBase()}),
			e: &event.MessageCreate{
				Base: event.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{Content: "abc"},
				},
			},
			expect: &discord.Message{Content: "abc"},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			_, s := state.NewMocker(t)

			c.waiter.state = s
			c.waiter.ctx = &plugin.Context{
				Localizer: mock.NewLocalizer(t).Build(),
			}

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
	t.Parallel()

	testCases := []struct {
		name   string
		waiter *ReplyWaiter
		e      *event.MessageReactionAdd
		expect interface{}
	}{
		{
			name: "message id not matching",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base:    event.NewBase(),
				Message: discord.Message{ChannelID: 456},
			}),
			e: &event.MessageReactionAdd{
				Base: event.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					MessageID: 321,
				},
			},
			expect: nil,
		},
		{
			name: "emoji not matching",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base:    event.NewBase(),
				Message: discord.Message{ChannelID: 456},
			}),
			e: &event.MessageReactionAdd{
				Base: event.NewBase(),
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
				Base: event.NewBase(),
				Message: discord.Message{
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				},
			}),
			e: &event.MessageReactionAdd{
				Base: event.NewBase(),
				MessageReactionAddEvent: &gateway.MessageReactionAddEvent{
					UserID: 987,
				},
			},
			expect: nil,
		},
		{
			name: "success",
			waiter: NewReplyWaiter(nil, &plugin.Context{
				Base: event.NewBase(),
				Message: discord.Message{
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				},
			}),
			e: &event.MessageReactionAdd{
				Base: event.NewBase(),
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
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			m, s := state.NewMocker(t)

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
