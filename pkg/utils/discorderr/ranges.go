package discorderr

import "github.com/diamondburned/arikawa/v2/utils/httputil"

// CodeRange is a list of code ranges.
// A range consists of two elements:
//
// Element 1 contains the inclusive first code of the range
// Element 2 contains the inclusive last code of the range.
type CodeRange [][2]httputil.ErrorCode

var (
	UnknownResource = CodeRange{
		{UnknownAccount, UnknownWebhook},
		{UnknownBan, UnknownBranch},
		{UnknownRedistributable, UnknownRedistributable},
		{UnknownGuildTemplate, UnknownGuildTemplate},
	}

	RateLimit = CodeRange{
		{AnnouncementEditRateLimit, AnnouncementEditRateLimit},
		{ChannelWriteRateLimit, ChannelWriteRateLimit},
	}

	ResourceLimit = CodeRange{
		{MaxGuilds, MaxPins},
		{MaxRoles, MaxRoles},
		{MaxWebhooks, MaxWebhooks},
		{MaxReactions, MaxReactions},
		{MaxGuildChannels, MaxGuildChannels},
		{MaxAttachments, MaxInvites},
	}

	OAuthError = CodeRange{
		{OAuthNoBot, InvalidOAuthState},
		{InvalidOAuthAccessToken, InvalidOAuthAccessToken},
	}
)
