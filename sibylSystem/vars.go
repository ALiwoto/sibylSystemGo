// sibylSystemGo library Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystemGo

import "errors"

// error variables
var (
	ErrInvalidHostUrl = errors.New("invalid host url")
	ErrInvalidToken   = errors.New("token length should be more than 20")
	ErrNoReason       = errors.New("reason is required for this action")
)
