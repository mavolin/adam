package i18nutil

import (
	"sort"

	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/i18n"
)

// PermissionNames returns the sorted names of the passed discord.Permissions,
// as found in the client.
func PermissionNames(perms discord.Permissions) []string {
	l := i18n.NewManager(nil).Localizer("")
	return PermissionNamesl(perms, l)
}

// PermissionNamel returns the sorted and localized names of the passed
// discord.Permissions, as found in the client.
func PermissionNamesl(perms discord.Permissions, l *i18n.Localizer) (s []string) {
	for perm, c := range permissionConfigs {
		if perms.Has(perm) {
			permString, _ := l.Localize(c)
			s = append(s, permString)
		}
	}

	sort.Strings(s)

	return s
}
