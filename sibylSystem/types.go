// sibylSystemGo library Project
// Copyright (C) 2021-2022 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystem

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AnimeKaizoku/ssg/ssg"
)

type UserPermission int
type EntityType int
type PollingUniqueId uint64
type BanFlag string
type SibylUpdateType string

type sibylCore struct {
	Token      string
	HostUrl    string
	Context    context.Context
	HttpClient *http.Client
}

type SibylConfig struct {
	HostUrl    string
	HttpClient *http.Client
	Context    context.Context
}

type SibylDispatcher struct {
	PollingId          *PollingIdentifier
	TimeoutSeconds     int
	MaxConnectionTries int
	isStopped          bool
	sibylClient        SibylClient
	totalTries         int

	onStartFailed     func(error)
	onGetUpdateFailed func(error)
	onHandlerError    func(error)

	handlers *ssg.SafeMap[SibylUpdateType, []ServerUpdateHandler]
}

type SibylClient interface {
	// ChangeToken changes token of the current SibylClient.
	// returns error if any.
	ChangeToken(token string) error

	// ChangeUrl changes host url of the current SibylClient.
	ChangeUrl(hostUrl string) error

	// ChangeToDefaultUrl changes host url of the current SibylClient to default url.
	ChangeToDefaultUrl()

	// GetHostUrl returns host url of the current SibylClient.
	GetHostUrl() string

	// Ban bans user with given id, reason and BanConfig.
	Ban(userId int64, reason string, config *BanConfig) (*BanResult, error)

	// BanUser bans a "user" with given id, reason and BanConfig.
	// entityType will be set to "user".
	BanUser(userId int64, reason string, config *BanConfig) (*BanResult, error)

	// BanBot bans a "bot" with given id, reason and BanConfig.
	// entityType will be set to "bot".
	BanBot(userId int64, reason string, config *BanConfig) (*BanResult, error)

	// RemoveBan removes ban from user with given id.
	RemoveBan(userId int64, reason string, config *RevertConfig) (string, error)

	// RevertBan reverts the ban from user with given id.
	RevertBan(userId int64, reason string, config *RevertConfig) (string, error)

	// FullRevert will fully revert the target user, they won't get `Restored` status,
	// all of their bans history will be deleted.
	// This method requires high token permission.
	FullRevert(userId int64, config *FullRevertConfig) (string, error)

	// GetInfo returns information about the user with given id.
	GetInfo(userId int64) (*GetInfoResult, error)

	// GetGeneralInfo returns information about the user with given id.
	// if the user is not a registered user at PSB, server will return error.
	GetGeneralInfo(userId int64) (*GeneralInfoResult, error)

	// GetGetAllBannedUsers returns information about all banned users.
	GetGetAllBannedUsers() (*GetBansResult, error)

	// GetStats returns current server stats.
	GetStats() (*GetStatsResult, error)

	// CheckToken checks if the token is valid.
	CheckToken() (bool, error)

	// Report reports a user with given id, reason and ReportConfig.
	Report(userId int64, reason string, config *ReportConfig) (string, error)

	// Scan scans a user with given id, reason and ReportConfig.
	Scan(userId int64, reason string, config *ReportConfig) (string, error)

	// ReportUser reports a "user" with given id, reason and ReportConfig.
	// IsBot parameter will be set to false.
	ReportUser(userId int64, reason string, config *ReportConfig) (string, error)

	// ScanUser scans a "user" with given id, reason and ReportConfig.
	// IsBot parameter will be set to false.
	ScanUser(userId int64, reason string, config *ReportConfig) (string, error)

	// ReportBot reports a "bot" with given id, reason and ReportConfig.
	// IsBot parameter will be set to true.
	ReportBot(userId int64, reason string, config *ReportConfig) (string, error)

	// CreateToken creates a new token in the server-side.
	CreateToken(userId int64) (*TokenInfo, error)

	// ChangePermission changes permission of the user with given id.
	ChangePermission(userId int64, perm UserPermission) (*ChangePermResult, error)

	// RevokeToken revokes the token of the user with given id.
	// It needs owner permission if the user-id doesn't belong to yourself.
	RevokeToken(userId int64) (*TokenInfo, error)

	// GetToken returns the token of the user with given id.
	// it needs owner permission if the user-id doesn't belong to yourself.
	GetToken(userId int64) (*TokenInfo, error)

	// GetAllRegisteredUsers returns information about all registered users.
	GetAllRegisteredUsers() (*GetRegisteredResult, error)

	// StartPolling method will sends a new StartPolling request to the server.
	// as of now, this method can only be used by users with permission more than
	// inspector. this method will return the unique id of the polling process.
	// later on, for getting updates from server, you should pass this unique-id.
	StartPolling() (*PollingIdentifier, error)

	// GetUpdates will send a GetUpdates request to the sibyl's servers, the response
	// might be (nil, nil), which means getting data got timed out. normally, you have
	// to call this method consequently if you want to remain up-to-date with server's
	// events. preferably, pass the unique-id you have got from StartPolling method as
	// second arg (second arg is not mandatory, and can be set to 0).
	GetUpdates(timeout int, uniqueId *PollingIdentifier) (*ServerUpdateContainer, error)

	// String returns string representation of the current SibylClient.
	String() string

	// Println prints string representation of the current SibylClient using fmt.Println.
	Println()

	// Print prints string representation of the current SibylClient using fmt.Print.
	Print()
}

type SibylError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Origin  string `json:"origin"`
	Date    string `json:"date"`
}

type CymaticScanConfig struct {
	Message    string
	SrcUrl     string
	TargetType EntityType
	TheToken   string
	PollingId  *PollingIdentifier
}

type FullRevertConfig struct {
	TheToken string
}

type BanConfig = CymaticScanConfig

type ReportConfig = CymaticScanConfig

type ScanConfig = ReportConfig

type RevertConfig = CymaticScanConfig

// Add ban types:

type BanResult struct {
	PreviousBan *BanInfo `json:"previous_ban"`
	CurrentBan  *BanInfo `json:"current_ban"`
}

type AddBanResponse struct {
	Success bool        `json:"success"`
	Result  *BanResult  `json:"result"`
	Error   *SibylError `json:"error"`
}

type BanInfo struct {
	UserId           int64      `json:"user_id"`
	Banned           bool       `json:"banned"`
	Reason           string     `json:"reason"`
	Message          string     `json:"message"`
	BanSourceUrl     string     `json:"ban_source_url"`
	BannedBy         int64      `json:"banned_by"`
	CrimeCoefficient int64      `json:"crime_coefficient"`
	Date             string     `json:"date"`
	BanFlags         []string   `json:"ban_flags"`
	TargetType       EntityType `json:"target_type"`
}

// Remove ban types:

type RemoveBanResponse struct {
	Success bool        `json:"success"`
	Result  string      `json:"result"`
	Error   *SibylError `json:"error"`
}

type FullRevertResponse struct {
	Success bool        `json:"success"`
	Result  string      `json:"result"`
	Error   *SibylError `json:"error"`
}

// get info types:

type GetInfoResponse struct {
	Success bool           `json:"success"`
	Result  *GetInfoResult `json:"result"`
	Error   *SibylError    `json:"error"`
}

type GetInfoResult struct {
	UserId           int64      `json:"user_id"`
	Banned           bool       `json:"banned"`
	Reason           string     `json:"reason"`
	Message          string     `json:"message"`
	BanSourceUrl     string     `json:"ban_source_url"`
	BannedBy         int64      `json:"banned_by"`
	CrimeCoefficient int        `json:"crime_coefficient"`
	Date             string     `json:"date"`
	BanFlags         []BanFlag  `json:"ban_flags"`
	TargetType       EntityType `json:"target_type"`
}

// general info types

type GeneralInfoResponse struct {
	Success bool               `json:"success"`
	Result  *GeneralInfoResult `json:"result"`
	Error   *SibylError        `json:"error"`
}

type GeneralInfoResult struct {
	UserId         int64          `json:"user_id"`
	Division       int            `json:"division"`
	AssignedBy     int64          `json:"assigned_by"`
	AssignedReason string         `json:"assigned_reason"`
	AssignedAt     string         `json:"assigned_at"`
	Permission     UserPermission `json:"permission"`
}

// get bans types:

type GetBansResponse struct {
	Success bool           `json:"success"`
	Result  *GetBansResult `json:"result"`
	Error   *SibylError    `json:"error"`
}

type GetBansResult struct {
	Users []BanInfo `json:"users"`
}

// get stats types:

type GetStatsResponse struct {
	Success bool            `json:"success"`
	Result  *GetStatsResult `json:"result"`
	Error   *SibylError     `json:"error"`
}

type GetStatsResult struct {
	BannedCount          int64 `json:"banned_count"`
	TrollingBanCount     int64 `json:"trolling_ban_count"`
	SpamBanCount         int64 `json:"spam_ban_count"`
	EvadeBanCount        int64 `json:"evade_ban_count"`
	CustomBanCount       int64 `json:"custom_ban_count"`
	PsychoHazardBanCount int64 `json:"psycho_hazard_ban_count"`
	MalImpBanCount       int64 `json:"mal_imp_ban_count"`
	NSFWBanCount         int64 `json:"nsfw_ban_count"`
	SpamBotBanCount      int64 `json:"spam_bot_ban_count"`
	RaidBanCount         int64 `json:"raid_ban_count"`
	MassAddBanCount      int64 `json:"mass_add_ban_count"`
	CloudyCount          int64 `json:"cloudy_count"`
	TokenCount           int64 `json:"token_count"`
	InspectorsCount      int64 `json:"inspectors_count"`
	EnforcesCount        int64 `json:"enforces_count"`
}

// check token types:

type CheckTokenResponse struct {
	Success bool        `json:"success"`
	Result  bool        `json:"result"`
	Error   *SibylError `json:"error"`
}

// report types:

type ReportResponse struct {
	Success bool        `json:"success"`
	Result  string      `json:"result"`
	Error   *SibylError `json:"error"`
}

// create token types:

type CreateTokenResponse struct {
	Success bool        `json:"success"`
	Result  *TokenInfo  `json:"result"`
	Error   *SibylError `json:"error"`
}

type TokenInfo struct {
	UserId          int64          `json:"user_id" gorm:"primaryKey"`
	Hash            string         `json:"hash"`
	Permission      UserPermission `json:"permission"`
	CreatedAt       string         `json:"created_at"`
	AcceptedReports int            `json:"accepted_reports"`
	DeniedReports   int            `json:"denied_reports"`
	AssignedBy      int64          `json:"assigned_by"`
	DivisionNum     int            `json:"division_num"`
	AssignedReason  string         `json:"assigned_reason"`
	cachedTime      time.Time      `json:"-"`
}

// change permission types:

type ChangePermResponse struct {
	Success bool              `json:"success"`
	Result  *ChangePermResult `json:"result"`
	Error   *SibylError       `json:"error"`
}

type ChangePermResult struct {
	PreviousPerm UserPermission `json:"previous_perm"`
	CurrentPerm  UserPermission `json:"current_perm"`
}

// revoke token types:

type RevokeTokenResponse struct {
	Success bool        `json:"success"`
	Result  *TokenInfo  `json:"result"`
	Error   *SibylError `json:"error"`
}

// get token types:

type GetTokenResponse struct {
	Success bool        `json:"success"`
	Result  *TokenInfo  `json:"result"`
	Error   *SibylError `json:"error"`
}

// get registered users types:

type GetRegisteredResponse struct {
	Success bool                 `json:"success"`
	Result  *GetRegisteredResult `json:"result"`
	Error   *SibylError          `json:"error"`
}

type GetRegisteredResult struct {
	RegisteredUsers []int64 `json:"registered_users"`
}

// polling-related structs

// PollingIdentifier represents a unique polling identifier.
type PollingIdentifier struct {
	PollingUniqueId   PollingUniqueId `json:"polling_unique_id"`
	PollingAccessHash string          `json:"polling_access_hash"`
}

type StartPollingResponse struct {
	Success bool               `json:"success"`
	Result  *PollingIdentifier `json:"result"`
	Error   *SibylError        `json:"error"`
}

type GetUpdateResponse struct {
	Success bool                   `json:"success"`
	Result  *ServerUpdateContainer `json:"result"`
	Error   *SibylError            `json:"error"`
}

type ServerUpdateContainer struct {
	UpdateType SibylUpdateType `json:"update_type"`
	UpdateData json.RawMessage `json:"update_data"`
}

type ScanRequestApprovedUpdate struct {
	UniqueId    string     `json:"unique_id"`
	TargetUser  int64      `json:"target_user"`
	TargetType  EntityType `json:"target_type"`
	AgentReason string     `json:"agent_reason"`
}

type ScanRequestRejectedUpdate struct {
	UniqueId    string     `json:"unique_id"`
	TargetUser  int64      `json:"target_user"`
	TargetType  EntityType `json:"target_type"`
	AgentReason string     `json:"agent_reason"`
}

type SibylUpdateContext struct {
	ScanRequestApproved *ScanRequestApprovedUpdate
	ScanRequestRejected *ScanRequestApprovedUpdate
}

type ServerUpdateHandler func(client SibylClient, ctx *SibylUpdateContext) error
