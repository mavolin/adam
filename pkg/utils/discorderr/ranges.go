package discorderr

import "github.com/diamondburned/arikawa/v2/utils/httputil"

var (
	UnknownResource = []httputil.ErrorCode{
		UnknownAccount, UnknownApplication, UnknownChannel, UnknownGuild,
		UnknownIntegration, UnknownInvite, UnknownMember, UnknownMessage,
		UnknownPermissionOverwrite, UnknownProvider, UnknownRole, UnknownToken,
		UnknownUser, UnknownEmoji, UnknownWebhook, UnknownBan, UnknownSKU,
		UnknownStoreListing, UnknownEntitlement, UnknownBuild, UnknownLobby,
		UnknownBranch, UnknownRedistributable, UnknownGuildTemplate,
	}

	ResourceLimit = []httputil.ErrorCode{
		MaxGuilds, MaxFriends, MaxPins, MaxRoles, MaxWebhooks, MaxReactions,
		MaxGuildChannels, MaxAttachments, MaxInvites,
	}

	OAuthError = []httputil.ErrorCode{
		OAuthNoBot, OAuthLimitReached, InvalidOAuthState,
		InvalidOAuthAccessToken,
	}
)
