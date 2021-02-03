package command

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// LocalizedMeta is the localized implementation of the plugin.CommandMeta
// interface.
type LocalizedMeta struct {
	// Name is the name of the command.
	// It may not contain whitespace or dots.
	Name string
	// Aliases are the optional aliases of the command.
	// They may not contain whitespace or dots.
	Aliases []string
	// ShortDescription is an optional short description of the command.
	ShortDescription *i18n.Config
	// LongDescription is an optional long description of the command.
	LongDescription *i18n.Config
	// Examples contains optional example usages of the command.
	Examples []*i18n.Config
	// Args is the argument configuration of the command.
	// If this is left empty, the command won't accept any arguments.
	Args plugin.ArgConfig
	// Hidden specifies whether this command should be hidden from the help
	// message.
	Hidden bool
	// ChannelTypes are the plugin.ChannelTypes the command may be executed in.
	//
	// If this is not set, the channel types of the parent will be used.
	ChannelTypes plugin.ChannelTypes
	// BotPermissions are the permissions the bot needs to execute this
	// command.
	BotPermissions discord.Permissions
	// Restrictions contains the restrictions of the command.
	//
	// If this is nil, the restrictions of the parent will be used.
	// Use restriction.None to prevent inheritance.
	Restrictions plugin.RestrictionFunc
	// Throttler is the plugin.Throttler of the command.
	//
	// If none is set, the throttler of the parent will be used.
	// Use throttler.None to prevent inheritance.
	Throttler plugin.Throttler
}

var _ plugin.CommandMeta = LocalizedMeta{}

func (m LocalizedMeta) GetName() string      { return m.Name }
func (m LocalizedMeta) GetAliases() []string { return m.Aliases }

func (m LocalizedMeta) GetShortDescription(l *i18n.Localizer) string {
	desc, err := l.Localize(m.ShortDescription)
	if err != nil {
		return ""
	}

	return desc
}

func (m LocalizedMeta) GetLongDescription(l *i18n.Localizer) string {
	desc, err := l.Localize(m.LongDescription)
	if err != nil {
		return ""
	}

	return desc
}

func (m LocalizedMeta) GetExamples(l *i18n.Localizer) []string {
	if len(m.Examples) == 0 {
		return nil
	}

	examples := make([]string, 0, len(m.Examples))

	for _, lexample := range m.Examples {
		example, err := l.Localize(lexample)
		if err == nil {
			examples = append(examples, example)
		}
	}

	if len(examples) == 0 {
		return nil
	}

	return examples
}

func (m LocalizedMeta) GetArgs() plugin.ArgConfig                  { return m.Args }
func (m LocalizedMeta) IsHidden() bool                             { return m.Hidden }
func (m LocalizedMeta) GetChannelTypes() plugin.ChannelTypes       { return m.ChannelTypes }
func (m LocalizedMeta) GetBotPermissions() discord.Permissions     { return m.BotPermissions }
func (m LocalizedMeta) GetRestrictionFunc() plugin.RestrictionFunc { return m.Restrictions }
func (m LocalizedMeta) GetThrottler() plugin.Throttler             { return m.Throttler }
