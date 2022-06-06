// sibylSystemGo library Project
// Copyright (C) 2021-2022 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystem

import (
	"context"
	"net/http"
	"strings"

	"github.com/AnimeKaizoku/ssg/ssg"
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

	if core.Context == nil {
		core.Context = context.Background()
	}

	if core.HostUrl[len(core.HostUrl)-1] != '/' {
		core.HostUrl += "/"
	}

	return core
}

func GetNewDispatcher(client SibylClient) *SibylDispatcher {
	return &SibylDispatcher{
		TimeoutSeconds:     DefaultDispatcherTimeout,
		MaxConnectionTries: 50,
		sibylClient:        client,
		handlers:           ssg.NewSafeMap[SibylUpdateType, []ServerUpdateHandler](),
	}
}

func GetDefaultConfig() *SibylConfig {
	return &SibylConfig{
		HostUrl:    DefaultUrl,
		HttpClient: http.DefaultClient,
		Context:    context.Background(),
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
	if len(value) < 3 {
		return DefaultUrl
	}

	if value[len(value)-1] != '/' {
		value += "/"
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value
	}

	// animekaizoku's domains are mostly protected by cloudflare shit,
	// so we need to use https:// for them.
	if strings.Contains(value, "animekaizoku") {
		return "https://" + value
	}

	return "http://" + value
}
