package discordutil

import (
	"strings"

	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/utils/locutil"
)

// PermissionNames returns the names of the passed discord.Permissions, as
// found in the client.
func PermissionNames(perms discord.Permissions) []string {
	l := localization.NewManager(nil).Localizer("")
	return PermissionNamesl(perms, l)
}

// PermissionNamel returns the localized names of the passed
// discord.Permissions, as found in the client.
func PermissionNamesl(perms discord.Permissions, l *localization.Localizer) (s []string) {
	for perm, c := range permissionConfigs {
		if perms.Has(perm) {
			permString, _ := l.Localize(c)
			s = append(s, permString)
		}
	}

	return s
}

// PermissionList creates a sorted written list of the passed
// discord.Permissions.
func PermissionList(perms discord.Permissions) string {
	const (
		defaultSep = ", "
		lastSep    = " and "
	)

	names := PermissionNames(perms)

	var b strings.Builder

	if len(names) > 2 {
		b.Grow((len(names) - 2) * len(defaultSep))
	}

	b.Grow(len(lastSep))

	for _, s := range names {
		b.Grow(len(s))
	}

	for i, s := range names {
		b.WriteString(s)

		if i < len(names)-2 {
			b.WriteString(defaultSep)
		} else if i == len(names)-2 {
			b.WriteString(lastSep)
		}
	}

	return b.String()
}

// PermissionListl creates a written list from the passed permissions using the
// passed localization.Localizer.
func PermissionListl(perms discord.Permissions, l *localization.Localizer) string {
	var cfgs []localization.Config

	for perm, c := range permissionConfigs {
		if perms.Has(perm) {
			cfgs = append(cfgs, c)
		}
	}

	// we can ignore the error because all translations have fallbacks
	s, _ := locutil.ConfigsToList(cfgs, l)

	return s
}
