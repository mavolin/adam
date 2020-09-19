package mock

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

type Module struct {
	MetaReturn     plugin.ModuleMeta
	CommandsReturn []plugin.Command
	ModulesReturn  []plugin.Module
}

func (c Module) Meta() plugin.ModuleMeta    { return c.MetaReturn }
func (c Module) Commands() []plugin.Command { return c.CommandsReturn }
func (c Module) Modules() []plugin.Module   { return c.ModulesReturn }

type RegisteredModule struct {
	ParentReturn plugin.RegisteredModule
	ParentError  error

	IdentifierReturn        plugin.Identifier
	NameReturn              string
	ShortDescriptionReturn  string
	LongDescriptionReturn   string
	IsHiddenReturn          bool
	ThrottlingOptionsReturn plugin.ThrottlingOptions
	CommandsReturn          []plugin.RegisteredCommand
	ModulesReturn           []plugin.RegisteredModule
}

func (r RegisteredModule) Parent() (plugin.RegisteredModule, error) {
	return r.ParentReturn, r.ParentError
}

func (r RegisteredModule) Identifier() plugin.Identifier { return r.IdentifierReturn }
func (r RegisteredModule) Name() string                  { return r.NameReturn }

func (r RegisteredModule) ShortDescription(*localization.Localizer) string {
	return r.ShortDescriptionReturn
}

func (r RegisteredModule) LongDescription(*localization.Localizer) string {
	return r.LongDescriptionReturn
}

func (r RegisteredModule) IsHidden() bool { return r.IsHiddenReturn }

func (r RegisteredModule) ThrottlingOptions() plugin.ThrottlingOptions {
	return r.ThrottlingOptionsReturn
}

func (r RegisteredModule) Commands() []plugin.RegisteredCommand {
	cp := make([]plugin.RegisteredCommand, len(r.CommandsReturn))
	copy(cp, r.CommandsReturn)

	return cp
}

func (r RegisteredModule) Modules() []plugin.RegisteredModule {
	cp := make([]plugin.RegisteredModule, len(r.ModulesReturn))
	copy(cp, r.ModulesReturn)

	return cp
}

func (r RegisteredModule) FindCommand(invoke string) plugin.RegisteredCommand {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	var mod plugin.RegisteredModule = r

Identifiers:
	for i := 0; i < len(all)-1; i++ {
		name := all[i].Name()

		for _, searchMod := range mod.Modules() {
			if searchMod.Name() == name {
				mod = searchMod
				continue Identifiers
			}
		}

		return nil
	}

	name := all[len(all)-1].Name()

	for _, cmd := range r.CommandsReturn {
		if cmd.Name() == name {
			return cmd
		}

		for _, alias := range cmd.Aliases() {
			if alias == name {
				return cmd
			}
		}
	}

	return nil
}

func (r RegisteredModule) FindModule(invoke string) plugin.RegisteredModule {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	var mod plugin.RegisteredModule = r

Identifiers:
	for _, id := range all {
		name := id.Name()

		for _, searchMod := range mod.Modules() {
			if searchMod.Name() == name {
				mod = searchMod
				continue Identifiers
			}
		}

		return nil
	}

	return nil
}

type ModuleMeta struct {
	Name              string
	ShortDescription  string
	LongDescription   string
	Hidden            bool
	ChannelTypes      plugin.ChannelTypes
	BotPermissions    *discord.Permissions
	Restrictions      plugin.RestrictionFunc
	ThrottlingOptions plugin.ThrottlingOptions
}

func (c ModuleMeta) GetName() string                                    { return c.Name }
func (c ModuleMeta) GetShortDescription(*localization.Localizer) string { return c.ShortDescription }
func (c ModuleMeta) GetLongDescription(*localization.Localizer) string  { return c.LongDescription }
func (c ModuleMeta) IsHidden() bool                                     { return c.Hidden }
func (c ModuleMeta) GetDefaultChannelTypes() plugin.ChannelTypes        { return c.ChannelTypes }
func (c ModuleMeta) GetDefaultBotPermissions() *discord.Permissions     { return c.BotPermissions }
func (c ModuleMeta) GetDefaultRestrictionFunc() plugin.RestrictionFunc  { return c.Restrictions }
func (c ModuleMeta) GetThrottlingOptions() plugin.ThrottlingOptions     { return c.ThrottlingOptions }
