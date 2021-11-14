// sibylSystemGo library Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystemGo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// error methods:

func (e *SibylError) Error() string {
	return http.StatusText(e.Code) + " [" + strconv.Itoa(e.Code) + "]: " + e.Message
}

// general and private methods:

func (s *sibylCore) revokeRequest(req *http.Request, result interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var b []byte

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		return err
	}
	return nil
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

func (s *sibylCore) Ban(userId int64, reason, message, srcUrl string,
	isBot bool) (*BanResult, error) {
	if len(reason) < 1 {
		return nil, ErrNoReason
	}

	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"addBan", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))
	req.Header.Add("reason", reason)
	req.Header.Add("message", message)
	req.Header.Add("srcUrl", srcUrl)
	req.Header.Add("is-bot", strconv.FormatBool(isBot))

	resp := new(AddBanResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success && resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) BanUser(userId int64, reason, message, srcUrl string) (*BanResult, error) {
	return s.Ban(userId, reason, message, srcUrl, false)
}

func (s *sibylCore) BanBot(userId int64, reason, message, srcUrl string) (*BanResult, error) {
	return s.Ban(userId, reason, message, srcUrl, true)
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
func (s *sibylCore) Report(userId int64, reason, message, srcUrl string, isBot bool) (string, error) {
	if len(reason) < 1 {
		return "", ErrNoReason
	}

	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"reportUser", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("token", s.Token)
	req.Header.Add("user-id", strconv.FormatInt(userId, 10))
	req.Header.Add("reason", reason)
	req.Header.Add("message", message)
	req.Header.Add("src", srcUrl)
	req.Header.Add("is-bot", strconv.FormatBool(isBot))

	resp := new(ReportResponse)

	err = s.revokeRequest(req, resp)
	if err != nil {
		return "", err
	}

	if !resp.Success && resp.Error != nil {
		return "", resp.Error
	}
	return resp.Result, nil
}

func (s *sibylCore) ReportUser(userId int64, reason, message, srcUrl string) (string, error) {
	return s.Report(userId, reason, message, srcUrl, false)
}

func (s *sibylCore) ReportBot(userId int64, reason, message, srcUrl string) (string, error) {
	return s.Report(userId, reason, message, srcUrl, true)
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
	req, err := http.NewRequest(http.MethodGet, s.HostUrl+"revokeToken", nil)
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

//---------------------------------------------------------
