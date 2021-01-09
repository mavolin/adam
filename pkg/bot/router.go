package bot

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
)

var ErrUnknownCommand = errors.NewUserErrorl(unknownCommandErrorDescription)

// Route attempts to route the passed message.
// It aborts, if the message is not a valid invoke.
func (b *Bot) Route(base *state.Base, msg *discord.Message, member *discord.Member) {
	// Only accept regular text messages.
	// Also check if a bot wrote the message, if !b.AllowBot.
	if msg.Type != discord.DefaultMessage || (!b.AllowBot && msg.Author.Bot) {
		return
	}

	prefixes, lang := b.SettingsProvider(base, msg)

	ok, invoke := b.hasPrefix(msg.Content, prefixes)
	if !ok { // not an invoke
		return
	}

	var args string

	ctx := &plugin.Context{
		Message:          *msg,
		Member:           member,
		Base:             base,
		Localizer:        b.LocalizationManager.Localizer(lang),
		Prefixes:         prefixes,
		BotOwnerIDs:      b.Owners,
		ReplyMiddlewares: b.ReplyMiddlewares,
		Replier:          replier.WrapState(b.State),
		DiscordDataProvider: &discordDataProvider{
			s:         b.State,
			guildID:   msg.GuildID,
			channelID: msg.ChannelID,
			selfID:    b.selfID,
		},
	}
	ctx.ErrorHandler = newCtxErrorHandler(b.State, ctx, b.ErrorHandler)

	if b.AsyncPluginProviders {
		ctx.InvokedCommand, ctx.Provider, args = b.routeCommandAsync(invoke, base, msg)
		if ctx.InvokedCommand == nil {
			b.ErrorHandler(ErrUnknownCommand, b.State, ctx)
		}
	} else {
		ctx.InvokedCommand, ctx.Provider, args = b.routeCommand(invoke, base, msg)
		if ctx.InvokedCommand == nil {
			b.ErrorHandler(ErrUnknownCommand, b.State, ctx)
		}
	}

	err := ctx.InvokedCommand.IsRestricted(b.State, ctx)
	if err != nil {
		b.ErrorHandler(err, b.State, ctx)
		return
	}

	if ctx.InvokedCommand.Args != nil {
		ctx.Args, ctx.Flags, err = ctx.InvokedCommand.Args.Parse(args, b.State, ctx)
		if err != nil {
			b.ErrorHandler(err, b.State, ctx)
			return
		}
	}

	_, err = ctx.InvokedCommand.Invoke(b.State, ctx)
	if err != nil {
		b.ErrorHandler(err, b.State, ctx)
		return
	}
}

// hasPrefix checks if the passed invoke starts with one of the passed
// prefixes or a mention of the bot.
// If so it returns true and the invoke stripped of the prefix.
func (b *Bot) hasPrefix(invoke string, prefixes []string) (bool, string) {
	indexes := b.selfMentionRegexp.FindStringIndex(invoke)
	if indexes != nil {
		return true, strings.TrimLeft(invoke[indexes[1]:], whitespace)
	}

	for _, p := range prefixes {
		if strings.HasPrefix(invoke, p) {
			return true, strings.TrimLeft(invoke[len(p):], whitespace)
		}
	}

	return false, ""
}

func (b *Bot) routeCommand(
	invoke string, base *state.Base, msg *discord.Message,
) (*plugin.RegisteredCommand, plugin.Provider, string) {
	ctxprovider := &ctxPluginProvider{
		base: base,
		msg:  msg,
		repos: []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands:     b.commands,
				Modules:      b.modules,
				Defaults:     b.PluginDefaults,
			},
		},
	}

	cmd, args := findCommand(invoke, ctxprovider, ctxprovider.repos[0])
	if cmd != nil {
		ctxprovider.remProviders = b.pluginProviders
		return cmd, ctxprovider, args
	}

	for i, p := range b.pluginProviders {
		cmds, mods, err := p.provider(base, msg)
		if err != nil {
			ctxprovider.unavailableProviders = append(ctxprovider.unavailableProviders, plugin.UnavailablePluginProvider{
				Name:  p.name,
				Error: err,
			})
		} else {
			repo := plugin.Repository{
				ProviderName: p.name,
				Modules:      mods,
				Commands:     cmds,
				Defaults:     p.defaults,
			}

			ctxprovider.repos = append(ctxprovider.repos, repo)

			cmd, args := findCommand(invoke, ctxprovider, repo)
			if cmd != nil {
				ctxprovider.remProviders = b.pluginProviders[i+1:]
				return cmd, ctxprovider, args
			}
		}
	}

	return nil, nil, ""
}

func (b *Bot) routeCommandAsync(
	invoke string, base *state.Base, msg *discord.Message,
) (*plugin.RegisteredCommand, plugin.Provider, string) {
	ctxprovider := &ctxPluginProvider{
		base: base,
		msg:  msg,
		repos: []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands:     b.commands,
				Modules:      b.modules,
				Defaults:     b.PluginDefaults,
			},
		},
		async: true,
	}

	cmd, args := findCommand(invoke, ctxprovider, ctxprovider.repos[0])
	if cmd != nil {
		ctxprovider.remProviders = b.pluginProviders
		return cmd, ctxprovider, args
	}

	repos, up := pluginProvidersAsync(base, msg, b.pluginProviders)
	ctxprovider.repos = append(ctxprovider.repos, repos...)
	ctxprovider.unavailableProviders = append(ctxprovider.unavailableProviders, up...)

	for _, r := range repos {
		cmd, args := findCommand(invoke, ctxprovider, r)
		if cmd != nil {
			return cmd, ctxprovider, args
		}
	}

	return nil, nil, ""
}
