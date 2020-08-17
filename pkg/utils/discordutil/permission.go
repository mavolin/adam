package discordutil

import (
	"sort"

	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
)

// PermissionNames returns the sorted names of the passed discord.Permissions,
// as found in the client.
func PermissionNames(perms discord.Permissions) []string {
	l := localization.NewManager(nil).Localizer("")
	return PermissionNamesl(perms, l)
}

// PermissionNamel returns the sorted and localized names of the passed
// discord.Permissions, as found in the client.
func PermissionNamesl(perms discord.Permissions, l *localization.Localizer) (s []string) {
	for perm, c := range permissionConfigs {
		if perms.Has(perm) {
			permString, _ := l.Localize(c)
			s = append(s, permString)
		}
	}

	sort.Strings(s)

	return s
}
