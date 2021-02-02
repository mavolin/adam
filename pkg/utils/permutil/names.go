package permutil

import (
	"sort"

	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/i18n"
)

// Names returns the sorted names of the passed discord.Permissions, as found
// in the client.
func Names(perms discord.Permissions) []string {
	return Namesl(perms, i18n.NewFallbackLocalizer())
}

// PermissionNamel returns the sorted and localized names of the passed
// discord.Permissions, as found in the client.
func Namesl(perms discord.Permissions, l *i18n.Localizer) (s []string) {
	for perm, c := range permissionConfigs {
		if perms.Has(perm) {
			permString, _ := l.Localize(c)
			s = append(s, permString)
		}
	}

	sort.Strings(s)

	return s
}
