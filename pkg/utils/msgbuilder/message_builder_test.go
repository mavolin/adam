package msgbuilder

import (
	"context"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestBuilder_Await(t *testing.T) {
	t.Parallel()

	t.Run("response", func(t *testing.T) {
		t.Parallel()
		t.Run("initial timeout", func(t *testing.T) {
			t.Parallel()

			_, s := state.NewMocker(t)

			s.ApplyGateways(func(g *gateway.Gateway) {
				g.AddIntents(gateway.IntentGuildMessageTyping)
			})

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

			expect := &TimeoutError{UserID: ctx.Author.ID, Cause: context.DeadlineExceeded}

			actual := New(s, ctx).
				WithAwaitedResponse(new(discord.Message), 1, 1).
				Await(10, false)
			assert.Equal(t, expect, actual)
		})
	})
}

func TestReplyWaiter_handleMessages(t *testing.T) {
	t.Parallel()

	testError := errors.New("abc")

	testCases := []struct {
		name      string
		waiter    *Builder
		e         *event.MessageCreate
		expect    *discord.Message
		expectErr error
	}{
		{
			name: "middleware block",
			waiter: New(nil, &plugin.Context{Base: event.NewBase()}).
				WithResponseMiddlewares(func(*state.State, *event.MessageCreate) error {
					return testError
				}),
			e: &event.MessageCreate{
				Base:               event.NewBase(),
				MessageCreateEvent: new(gateway.MessageCreateEvent),
			},
			expectErr: testError,
		},
		{
			name: "channel not matching",
			waiter: New(nil, &plugin.Context{
				Base:    event.NewBase(),
				Message: discord.Message{ChannelID: 123},
			}),
			e: &event.MessageCreate{
				Base: event.NewBase(),
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{ChannelID: 321},
				},
			},
			expectErr: &TimeoutError{Cause: context.DeadlineExceeded},
		},
		{
			name: "author not matching",
			waiter: New(nil, &plugin.Context{
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
			expectErr: &TimeoutError{Cause: context.DeadlineExceeded},
		},
		{
			name:   "success",
			waiter: New(nil, &plugin.Context{Base: event.NewBase()}),
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
			c.waiter.pluginCtx = &plugin.Context{
				Localizer: mock.NewLocalizer(t).Build(),
			}

			var doneChan chan error
			// cause a nil pointer dereference, if something gets sent anyway
			// although c.expect == nil
			if c.expect != nil {
				doneChan = make(chan error)
			}

			var actual discord.Message

			c.waiter.WithAwaitedResponse(&actual, 50*time.Millisecond, 0)

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			rm := c.waiter.handleMessages(ctx, doneChan)

			s.Call(c.e)

			select {
			case <-ctx.Done():
			case actualErr := <-doneChan:
				assert.Equal(t, c.expectErr, actualErr)
			}

			if c.expect != nil {
				assert.Equal(t, *c.expect, actual)
			}

			rm()
		})
	}
}
