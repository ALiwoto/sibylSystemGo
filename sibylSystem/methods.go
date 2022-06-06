// sibylSystemGo library Project
// Copyright (C) 2021-2022 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystem

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	urlLib "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	ws "github.com/AnimeKaizoku/ssg/ssg"
)

// error methods:

func (e *SibylError) Error() string {
	return http.StatusText(e.Code) + " [" + strconv.Itoa(e.Code) + "]: " + e.Message
}

//---------------------------------------------------------

// general and private methods:

func (s *sibylCore) revokeRequest(req *http.Request, result interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return s.readResp(resp, result)
}

func (s *sibylCore) readResp(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		return err
	}

	return nil
}

func (s *sibylCore) getRequest(url string, params urlLib.Values, result interface{}) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = params.Encode()

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return err
	}

	return s.readResp(resp, result)
}

func (s *sibylCore) String() string {
	return "SibylClient (as sibylCore): " + s.HostUrl
}

func (s *sibylCore) Stringln() string {
	return "SibylClient (as sibylCore): " + s.HostUrl + "\n"
}

func (s *sibylCore) Println() {
	fmt.Println(s.String())
}

func (s *sibylCore) Print() {
	fmt.Print(s.String())
}

// general and public methods:

func (s *sibylCore) ChangeToken(token string) error {
	if len(token) < 20 {
		return ErrInvalidToken
	}
	s.Token = token
	return nil
}

func (s *sibylCore) ChangeUrl(hostUrl string) error {
	if len(hostUrl) == 0 {
		s.ChangeToDefaultUrl()
		return nil
	}
	if len(hostUrl) < 4 {
		return ErrInvalidHostUrl
	}
	s.HostUrl = validateHostUrl(hostUrl)
	return nil
}

func (s *sibylCore) ChangeToDefaultUrl() {
	s.HostUrl = DefaultUrl
}

func (s *sibylCore) GetHostUrl() string {
	return s.HostUrl
}

// ban-related methods:

func (s *sibylCore) Ban(userId int64, reason string, config *BanConfig) (*BanResult, error) {
	if len(reason) < 1 {
		return nil, ErrNoReason
	}

	if config == nil {
		config = &BanConfig{}
	}

	var myToken string
	if config.TheToken != "" {
		myToken = config.TheToken
	} else {
		myToken = s.Token
	}

	v := urlLib.Values{}

	v.Add("token", myToken)
	v.Add("user-id", strconv.FormatInt(userId, 10))
	v.Add("reason", reason)
	v.Add("message", config.Message)
	v.Add("srcUrl", config.SrcUrl)
	v.Add("entity-type", config.TargetType.ToString())

	resp := new(AddBanResponse)

	err := s.getRequest(s.HostUrl+"addBan", v, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Result, nil
}

func (s *sibylCore) BanUser(userId int64, reason string, config *BanConfig) (*BanResult, error) {
	if config == nil {
		config = &BanConfig{}
	}

	config.TargetType = EntityTypeUser
	return s.Ban(userId, reason, config)
}

func (s *sibylCore) BanBot(userId int64, reason string, config *BanConfig) (*BanResult, error) {
	if config == nil {
		config = &BanConfig{}
	}

	config.TargetType = EntityTypeBot
	return s.Ban(userId, reason, config)
}

func (s *sibylCore) RemoveBan(userId int64) (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"remBan", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))

	resp := new(RemoveBanResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return "", err
	}

	if !resp.Success && resp.Error != nil {
		return "", resp.Error
	}

	return resp.Result, nil
}

// info methods:

func (s *sibylCore) GetInfo(userId int64) (*GetInfoResult, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"getInfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))

	resp := new(GetInfoResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) GetGeneralInfo(userId int64) (*GeneralInfoResult, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"getGeneralInfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))

	resp := new(GeneralInfoResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) GetGetAllBannedUsers() (*GetBansResult, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"getBans", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)

	resp := new(GetBansResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) GetStats() (*GetStatsResult, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"getStats", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)

	resp := new(GetStatsResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) CheckToken() (bool, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"checkToken", nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("token", s.Token)

	resp := new(CheckTokenResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return false, err
	}

	if !resp.Success && resp.Error != nil {
		return false, resp.Error
	}
	return resp.Result, nil
}

// report methods:
func (s *sibylCore) Report(userId int64, reason string, config *ReportConfig) (string, error) {
	if len(reason) < 1 {
		return "", ErrNoReason
	}

	if config == nil {
		config = &ReportConfig{}
	}

	var myToken string
	if config.TheToken != "" {
		myToken = config.TheToken
	} else {
		myToken = s.Token
	}

	v := urlLib.Values{}

	v.Add("token", myToken)
	v.Add("user-id", strconv.FormatInt(userId, 10))
	v.Add("reason", reason)
	v.Add("message", config.Message)
	v.Add("src", config.SrcUrl)
	v.Add("entity-type", config.TargetType.ToString())

	resp := new(ReportResponse)

	err := s.getRequest(s.HostUrl+"reportUser", v, resp)
	if err != nil {
		return "", err
	}

	if !resp.Success && resp.Error != nil {
		return "", resp.Error
	}

	return resp.Result, nil
}

func (s *sibylCore) Scan(userId int64, reason string, config *ReportConfig) (string, error) {
	return s.Report(userId, reason, config)
}

func (s *sibylCore) ReportUser(userId int64, reason string, config *ReportConfig) (string, error) {
	if config == nil {
		config = &ReportConfig{}
	}

	config.TargetType = EntityTypeUser
	return s.Report(userId, reason, config)
}

func (s *sibylCore) ScanUser(userId int64, reason string, config *ReportConfig) (string, error) {
	return s.ReportUser(userId, reason, config)
}

func (s *sibylCore) ReportBot(userId int64, reason string, config *ReportConfig) (string, error) {
	if config == nil {
		config = &ReportConfig{}
	}

	config.TargetType = EntityTypeBot
	return s.Report(userId, reason, config)
}

// token methods:

func (s *sibylCore) CreateToken(userId int64) (*TokenInfo, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"createToken", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))

	resp := new(CreateTokenResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) ChangePermission(userId int64, perm UserPermission) (*ChangePermResult, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"changePerm", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))

	resp := new(ChangePermResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) RevokeToken(userId int64) (*TokenInfo, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"revokeToken", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))

	resp := new(RevokeTokenResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) GetToken(userId int64) (*TokenInfo, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"getToken", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))

	resp := new(GetTokenResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) StartPolling() (uint64, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"startPolling", nil)
	if err != nil {
		return 0, err
	}

	req.Header.Add("token", s.Token)

	resp := new(StartPollingResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success && resp.Error != nil {
		return 0, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) GetUpdates(timeout int, uniqueId uint64) (*ServerUpdateContainer, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"getUpdates", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("polling-timeout", ws.ToBase10(int64(timeout)))
	req.Header.Add("polling-unique-id", ws.ToBase10(int64(uniqueId)))

	resp := new(GetUpdateResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) GetAllRegisteredUsers() (*GetRegisteredResult, error) {
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"getStats", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)

	resp := new(GetRegisteredResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

//---------------------------------------------------------

func (t *TokenInfo) SetCachedTime(tCache time.Time) {
	t.cachedTime = tCache
}

func (t *TokenInfo) IsExpired(d time.Duration) bool {
	return time.Since(t.cachedTime) > d
}

func (t *TokenInfo) IsValid() bool {
	return len(t.Hash) > 20 && t.UserId != 0
}

func (t *TokenInfo) IsCitizen() bool {
	return t.Permission == NormalUser
}

func (t *TokenInfo) IsRegistered() bool {
	return t.Permission > NormalUser
}

func (t *TokenInfo) IsEnforcer() bool {
	return t.Permission == Enforcer
}

func (t *TokenInfo) IsInspector() bool {
	return t.Permission == Inspector
}

func (t *TokenInfo) IsOwner() bool {
	return t.Permission == Owner
}

//---------------------------------------------------------

// IsOwner returns true if the token's permission
// is owner.
func (p UserPermission) IsOwner() bool {
	return p == Owner
}

// IsInspector returns true if the token's permission
// is inspector.
func (p UserPermission) IsInspector() bool {
	return p == Inspector
}

// IsEnforcer returns true if the token's permission
// is enforcer.
func (p UserPermission) IsEnforcer() bool {
	return p == Enforcer
}

// IsRegistered returns true if the owner of this token is considered as
// a valid registered user in the system.
func (p UserPermission) IsRegistered() bool {
	return p > NormalUser
}

// CanReport returns true if the token with its current
// permission can report a user to sibyl system or not.
func (p UserPermission) CanReport() bool {
	return p > NormalUser
}

// CanBeReported returns true if the token with its current
// permission can be reported to sibyl system or not.
func (p UserPermission) CanBeReported() bool {
	return p == NormalUser
}

// CanBeBanned returns true if the token with its current
// permission can be banned on sibyl system or not.
func (p UserPermission) CanBeBanned() bool {
	return p == NormalUser
}

// HasRole returns true if and only if this token belongs to a
// user which has a role in the Sibyl System (is not a normal user).
func (p UserPermission) HasRole() bool {
	return p > NormalUser
}

// CanBan returns true if the token with its current
// permission can ban/unban a user from Sibyl System or not.
func (p UserPermission) CanBan() bool {
	return p > Enforcer
}

// CanCreateToken returns true if the token with its current
// permission can create tokens in Sibyl System or not.
func (p UserPermission) CanCreateToken() bool {
	return p > Inspector
}

// CanRevokeToken returns true if the token with its current
// permission can revoke tokens in Sibyl System or not.
func (p UserPermission) CanRevokeToken() bool {
	return p > Inspector
}

// CanSeeStats returns true if the token with its current
// permission can see stats of another tokens or not.
func (p UserPermission) CanSeeStats() bool {
	return p > Enforcer
}

// CanGetToken returns true if the token with its current
// permission can get the token of another user using their id
// or not.
func (p UserPermission) CanGetToken() bool {
	return p == Owner
}

// CanGetGeneralInfo returns true if the token with its current
// permission can get general info of a registered user using their id
// or not.
func (p UserPermission) CanGetGeneralInfo() bool {
	return p > NormalUser
}

// CanGetAllBans returns true if the token with its current
// permission can get all the banned users.
func (p UserPermission) CanGetAllBans() bool {
	return p > NormalUser
}

// CanGetRegisteredList returns true if the token with its current
// permission can get all the registered users.
func (p UserPermission) CanGetRegisteredList() bool {
	return p > NormalUser
}

// CanChangePermission returns true if the token with its current
// permission can change permission of another tokens or not.
func (p UserPermission) CanChangePermission(pre, target UserPermission) bool {
	return !(p < Inspector || pre >= p || target >= p)
}

// CanTryChangePermission returns true if the token with its current
// permission can try to change permission of another tokens or not.
func (p UserPermission) CanTryChangePermission(direct bool) bool {
	if direct {
		return p > Inspector
	}

	return p > Enforcer
}

// CanGetStats returns true if the token with its current
// permission can get all stats of sibyl system or not.
func (p UserPermission) CanGetStats() bool {
	return p > Enforcer
}

//---------------------------------------------------------

func (r *GetInfoResult) IsPerma() bool {
	return strings.Contains(r.Reason, "perma")
}

func (r *GetInfoResult) HasCustomFlag() bool {
	return len(r.BanFlags) != 0 && r.BanFlags[0x0] == BanFlagCustom
}

func (r *GetInfoResult) SetAsBanReason(reason string) {
	r.Reason = reason
}

func (r *GetInfoResult) GetDateAsShort() string {
	return r.Date
}

func (r *GetInfoResult) EstimateCrimeCoefficient() string {
	c := r.CrimeCoefficient
	if c > 100 {
		str := strconv.Itoa(c)
		return "over " + str[:len(str)-2] + "00"
	}
	return "under 100"
}

func (r *GetInfoResult) GetStringCrimeCoefficient() string {
	return strconv.Itoa(r.CrimeCoefficient)
}

func (r *GetInfoResult) FormatFlags() mdparser.WMarkDown {
	md := mdparser.GetEmpty()
	if len(r.BanFlags) == 0 {
		return md
	}

	for i, current := range r.BanFlags {
		if i != 0 {
			md.Normal(", ")
		}
		md.Mono(string(current))
	}

	return md
}

func (r *GetInfoResult) FormatCuteFlags() mdparser.WMarkDown {
	md := mdparser.GetEmpty()
	if len(r.BanFlags) == 0 {
		return md
	} else if len(r.BanFlags) == 1 {
		return md.Normal(strings.ToLower(string(r.BanFlags[0x0])))
	}

	for i, current := range r.BanFlags {
		if i != 0 && i != len(r.BanFlags)-1 {
			md.Normal(", ")
		} else if i == len(r.BanFlags)-1 {
			md.Normal(" and ")
		}
		md.Normal(strings.ToLower(string(current)))
	}

	return md
}

func (r *GetInfoResult) EstimateCrimeCoefficientSep() (string, string) {
	c := r.CrimeCoefficient
	if c > 100 {
		str := strconv.Itoa(c)
		return "over ", str[:len(str)-2] + "00"
	}
	return "under ", "100"
}

//---------------------------------------------------------

func (d *SibylDispatcher) Listen() {
	d.totalTries = 0
	go d.StartListening()
}

func (d *SibylDispatcher) StartListening() {
	d.totalTries++
	if d.totalTries > d.MaxConnectionTries {
		// give up
		return
	}

	var err error
	d.PollingUniqueId, err = d.sibylClient.StartPolling()
	if err != nil {
		if d.onStartFailed != nil {
			d.onStartFailed(err)
		}
		return
	}

	var container *ServerUpdateContainer

	for !d.isStopped {
		container, err = d.sibylClient.GetUpdates(d.TimeoutSeconds, d.PollingUniqueId)
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "dial tcp") && strings.Contains(errStr, "connect: connection refused") {
				time.Sleep(time.Second)
				d.StartListening()
				return
			}
			if d.onGetUpdateFailed != nil {
				d.onGetUpdateFailed(err)
			}
		}

		// no updates, our request got timed out
		if container == nil {
			continue
		}

		// parse and handle each update in its own goroutine, to get the
		// best performance.
		go d.onUpdateReceived(container)
	}
}

func (d *SibylDispatcher) SetOnStartFailed(fn func(error)) {
	d.onStartFailed = fn
}

func (d *SibylDispatcher) SetOnGetUpdateFailed(fn func(error)) {
	d.onGetUpdateFailed = fn
}

func (d *SibylDispatcher) SetOnHandlerError(fn func(error)) {
	d.onHandlerError = fn
}

func (d *SibylDispatcher) AddHandler(uType SibylUpdateType, h ServerUpdateHandler) {
	handlers := d.handlers.GetValue(uType)
	handlers = append(handlers, h)
	d.handlers.Set(uType, handlers)
}

func (d *SibylDispatcher) onUpdateReceived(container *ServerUpdateContainer) {
	var err error
	ctx := new(SibylUpdateContext)
	switch container.UpdateType {
	case UpdateTypeScanRequestApproved:
		ctx.ScanRequestApproved = new(ScanRequestApprovedUpdate)
		err = json.Unmarshal(container.UpdateData, ctx.ScanRequestApproved)
	case UpdateTypeScanRequestRejected:
		ctx.ScanRequestRejected = new(ScanRequestApprovedUpdate)
		err = json.Unmarshal(container.UpdateData, ctx.ScanRequestRejected)
	default:
		// #TODO: add something for capturing this in future, idk
		return
	}

	if err != nil {
		if d.onGetUpdateFailed != nil {
			d.onGetUpdateFailed(err)
			return
		}
	}

	handlers := d.handlers.GetValue(container.UpdateType)
	if len(handlers) == 0 {
		return
	}

	for _, current := range handlers {
		err = current(d.sibylClient, ctx)
		if err != nil && d.onHandlerError != nil {
			d.onHandlerError(err)
		}
	}
}

//---------------------------------------------------------

func (e EntityType) ToString() string {
	return ws.ToBase10(int64(e))
}

func (e EntityType) IsBot() bool {
	return e == EntityTypeBot
}

func (e EntityType) IsBotStr() string {
	return ws.YesOrNo(e == EntityTypeBot)
}

func (e EntityType) IsUser() bool {
	return e == EntityTypeUser
}

func (e EntityType) IsUserStr() string {
	return ws.YesOrNo(e == EntityTypeUser)
}

func (e EntityType) IsAdmin() bool {
	return e == EntityTypeAdmin
}

func (e EntityType) IsAdminStr() string {
	return ws.YesOrNo(e == EntityTypeAdmin)
}

func (e EntityType) IsOwner() bool {
	return e == EntityTypeOwner
}

func (e EntityType) IsOwnerStr() string {
	return ws.YesOrNo(e == EntityTypeOwner)
}

func (e EntityType) IsGroup() bool {
	return e == EntityTypeGroup
}

func (e EntityType) IsGroupStr() string {
	return ws.YesOrNo(e == EntityTypeGroup)
}

func (e EntityType) IsChannel() bool {
	return e == EntityTypeChannel
}

func (e EntityType) IsChannelStr() string {
	return ws.YesOrNo(e == EntityTypeChannel)
}

func (e EntityType) IsChat() bool {
	return e == EntityTypeChannel || e == EntityTypeGroup
}

func (e EntityType) IsChatStr() string {
	return ws.YesOrNo(e == EntityTypeChannel || e == EntityTypeGroup)
}

func (e EntityType) IsOwnerOrAdmin() bool {
	return e == EntityTypeOwner || e == EntityTypeAdmin
}

func (e EntityType) IsOwnerOrAdminStr() string {
	return ws.YesOrNo(e == EntityTypeOwner || e == EntityTypeAdmin)
}
