package permutil

import (
	"math/bits"
	"sort"

	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/pkg/i18n"
)

// Names returns the sorted and localized names of the passed
// discord.Permissions, as found in the client.
func Names(l *i18n.Localizer, perms discord.Permissions) []string {
	names := make([]string, 0, bits.OnesCount64(uint64(perms)))

	for perm, c := range permissionConfigs {
		if perms.Has(perm) {
			permString, _ := l.Localize(c)
			names = append(names, permString)
		}
	}

	sort.Strings(names)

	return names
}
