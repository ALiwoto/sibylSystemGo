// sibylSystemGo library Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystemGo

const (
	DefaultUrl = "https://psychopass.animekaizoku.com/"
)

const (
	// NormalUser Can read from the Sibyl System.
	NormalUser UserPermission = iota
	// Enforcer Can only report to the Sibyl System.
	Enforcer
	// Inspector Can read/write directly to the Sibyl System.
	Inspector
	// Owner Can create/revoke tokens.
	Owner
)
