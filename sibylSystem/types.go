// sibylSystemGo library Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystemGo

import (
	"context"
	"net/http"
	"time"
)

type UserPermission int

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

type SibylClient interface {
	ChangeToken(token string) error
	ChangeUrl(hostUrl string) error
	ChangeToDefaultUrl()
	GetHostUrl() string
	Ban(userId int64, reason, message, srcUrl string,
		isBot bool) (*BanResult, error)
	BanUser(userId int64, reason, message, srcUrl string) (*BanResult, error)
	BanBot(userId int64, reason, message, srcUrl string) (*BanResult, error)
	RemoveBan(userId int64) (string, error)
	GetInfo(userId int64) (*GetInfoResult, error)
	GetGeneralInfo(userId int64) (*GeneralInfoResult, error)
	GetGetAllBannedUsers() (*GetBansResult, error)
	GetStats() (*GetStatsResult, error)
	CheckToken() (bool, error)
	Report(userId int64, reason, message, srcUrl string, isBot bool) (string, error)
	ReportUser(userId int64, reason, message, srcUrl string) (string, error)
	ReportBot(userId int64, reason, message, srcUrl string) (string, error)
	CreateToken(userId int64) (*TokenInfo, error)
	ChangePermission(userId int64, perm UserPermission) (*ChangePermResult, error)
	RevokeToken(userId int64) (*TokenInfo, error)
	GetToken(userId int64) (*TokenInfo, error)
	GetAllRegisteredUsers() (*GetRegisteredResult, error)
	String() string
	Println()
	Print()
}

type SibylError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Origin  string `json:"origin"`
	Date    string `json:"date"`
}

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
	UserId           int64    `json:"user_id"`
	Banned           bool     `json:"banned"`
	Reason           string   `json:"reason"`
	Message          string   `json:"message"`
	BanSourceUrl     string   `json:"ban_source_url"`
	BannedBy         int64    `json:"banned_by"`
	CrimeCoefficient int64    `json:"crime_coefficient"`
	Date             string   `json:"date"`
	BanFlags         []string `json:"ban_flags"`
	IsBot            bool     `json:"is_bot"`
}

// Remove ban types:

type RemoveBanResponse struct {
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
	UserId           int64    `json:"user_id"`
	Banned           bool     `json:"banned"`
	Reason           string   `json:"reason"`
	Message          string   `json:"message"`
	BanSourceUrl     string   `json:"ban_source_url"`
	BannedBy         int64    `json:"banned_by"`
	CrimeCoefficient int64    `json:"crime_coefficient"`
	Date             string   `json:"date"`
	BanFlags         []string `json:"ban_flags"`
	IsBot            bool     `json:"is_bot"`
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
