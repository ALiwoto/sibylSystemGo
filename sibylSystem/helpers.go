// sibylSystemGo library Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystemGo

import (
	"net/http"
	"strings"
)

func NewClient(token string, config *SibylConfig) SibylClient {
	if config == nil {
		config = GetDefaultConfig()
	}

	core := &sibylCore{
		Token:      token,
		HostUrl:    validateHostUrl(config.HostUrl),
		HttpClient: config.HttpClient,
	}
	return core
}

func GetDefaultConfig() *SibylConfig {
	return &SibylConfig{
		HostUrl:    DefaultUrl,
		HttpClient: &http.Client{},
	}
}

func ToSibylError(err error) *SibylError {
	if err == nil {
		return nil
	}

	if sibylError, ok := err.(*SibylError); ok {
		return sibylError
	}
	return nil
}

func validateHostUrl(value string) string {
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value
	}
	if strings.Contains(value, "animekaizoku") {
		return "https://" + value
	}
	return value
}
