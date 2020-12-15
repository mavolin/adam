package discorderr

import "github.com/diamondburned/arikawa/v2/utils/httputil"

const (
	// GeneralError is a general error (such as a malformed request body,
	// amongst other things).
	GeneralError httputil.ErrorCode = 0

	UnknownAccount             httputil.ErrorCode = 10001
	UnknownApplication         httputil.ErrorCode = 10002
	UnknownChannel             httputil.ErrorCode = 10003
	UnknownGuild               httputil.ErrorCode = 10004
	UnknownIntegration         httputil.ErrorCode = 10005
	UnknownInvite              httputil.ErrorCode = 10006
	UnknownMember              httputil.ErrorCode = 10007
	UnknownMessage             httputil.ErrorCode = 10008
	UnknownPermissionOverwrite httputil.ErrorCode = 10009
	UnknownProvider            httputil.ErrorCode = 10010
	UnknownRole                httputil.ErrorCode = 10011
	UnknownToken               httputil.ErrorCode = 10012
	UnknownUser                httputil.ErrorCode = 10013
	UnknownEmoji               httputil.ErrorCode = 10014
	UnknownWebhook             httputil.ErrorCode = 10015
	UnknownBan                 httputil.ErrorCode = 10026
	UnknownSKU                 httputil.ErrorCode = 10027
	UnknownStoreListing        httputil.ErrorCode = 10028
	UnknownEntitlement         httputil.ErrorCode = 10029
	UnknownBuild               httputil.ErrorCode = 10030
	UnknownLobby               httputil.ErrorCode = 10031
	UnknownBranch              httputil.ErrorCode = 10032
	UnknownRedistributable     httputil.ErrorCode = 10036
	UnknownGuildTemplate       httputil.ErrorCode = 10057

	// NoBots signals that bots cannot use the endpoint.
	NoBots httputil.ErrorCode = 20001
	// OnlyBots signals that only bots can use the endpoint.
	OnlyBots httputil.ErrorCode = 20002

	// AnnouncementRatelimit signals that the message cannot be edited due to
	// announcement rate limits.
	AnnouncementEditRateLimit httputil.ErrorCode = 20022
	// ChannelWriteRateLimit signals that the channel you are writing has hit
	// the write rate limit.
	ChannelWriteRateLimit httputil.ErrorCode = 20028

	MaxGuilds        httputil.ErrorCode = 30001
	MaxFriends       httputil.ErrorCode = 30002
	MaxPins          httputil.ErrorCode = 30003
	MaxRoles         httputil.ErrorCode = 30005
	MaxWebhooks      httputil.ErrorCode = 30007
	MaxReactions     httputil.ErrorCode = 30010
	MaxGuildChannels httputil.ErrorCode = 30013
	MaxAttachments   httputil.ErrorCode = 30015
	MaxInvites       httputil.ErrorCode = 30016

	// Unauthorized signals that you need to provide a valid token.
	Unauthorized httputil.ErrorCode = 40001
	// NeedVerification signals that you need to verify your account in order
	// to perform the action.
	NeedVerification httputil.ErrorCode = 40002
	// RequestEntitySize signals that the request entity is too large.
	RequestEntitySize httputil.ErrorCode = 40005
	// TemporarilyDisabled signals that Discord has temporarily disabled the
	// feature server-side.
	TemporarilyDisabled       httputil.ErrorCode = 40006
	UserBanned                httputil.ErrorCode = 40007
	MessageAlreadyCrossposted httputil.ErrorCode = 40033

	MissingAccess      httputil.ErrorCode = 50001
	InvalidAccountType httputil.ErrorCode = 50002

	// DMActionUnavailable signals that you cannot execute the action on a dm
	// channel.
	DMActionUnavailable httputil.ErrorCode = 50003
	GuildWidgetDisabled httputil.ErrorCode = 50004
	// SystemMessageActionUnavailable signals that you cannot execute the
	// action on a system message
	SystemMessageActionUnavailable httputil.ErrorCode = 50021
	// ChannelTypeActionUnavailable signals that you cannot execute the action
	// on the channel type.
	ChannelTypeActionUnavailable httputil.ErrorCode = 50024

	// CannotEditFromOther signals that you cannot edit a message authored by
	// another user.
	CannotEditFromOther        httputil.ErrorCode = 50005
	CannotSendEmptyMessage     httputil.ErrorCode = 50006
	CannotMessageUser          httputil.ErrorCode = 50007
	CannotSendToVoiceChannel   httputil.ErrorCode = 50008
	ChannelVerificationTooHigh httputil.ErrorCode = 50009

	OAuthNoBot              httputil.ErrorCode = 50010
	OAuthLimitReached       httputil.ErrorCode = 50011
	InvalidOAuthState       httputil.ErrorCode = 50012
	InvalidOAuthAccessToken httputil.ErrorCode = 50025

	InsufficientPermissions httputil.ErrorCode = 50013
	InvalidAuthToken        httputil.ErrorCode = 50014

	NoteTooLong httputil.ErrorCode = 50015

	// MessageDeleteBounds signals that you provided too few or too many
	// messages to delete.
	// You must provide at least 2 and fewer than 100 messages to delete.
	MessageDeleteBounds httputil.ErrorCode = 50016
	PinDifferentChannel httputil.ErrorCode = 50019
	InviteCodeInvalid   httputil.ErrorCode = 50020
	InvalidRecipient    httputil.ErrorCode = 50033
	// TooOldToDelete signals that the message is too old to be bulk deleted.
	TooOldToDelete  httputil.ErrorCode = 50034
	InvalidFormBody httputil.ErrorCode = 50035

	// BotNotInJoinedGuild signals that an invite was accepted to a guild the
	// application's bot is not in.
	BotNotInJoinedGuild httputil.ErrorCode = 50036

	InvalidAPIVersion httputil.ErrorCode = 50041

	// ChannelRequiredForCommunityGuild signals that you cannot delete the
	// channel, as it is required for Community guilds.
	ChannelRequiredForCommunityGuild httputil.ErrorCode = 50074

	InvalidSticker httputil.ErrorCode = 50081

	ReactionBlocked httputil.ErrorCode = 90001

	APIResourceOverloaded httputil.ErrorCode = 130000
)
