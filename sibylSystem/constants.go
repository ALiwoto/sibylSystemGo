// sibylSystemGo library Project
// Copyright (C) 2021-2022 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package sibylSystem

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

const (
	// EntityTypeUser represents a normal user while being scanned.
	// please notice that "being normal", doesn't necessarily mean
	// not being criminal.
	EntityTypeUser EntityType = iota
	// EntityTypeBot represents an account which is considered as a bot.
	// as API has no idea what is a "bot account", the value "is_bot"
	// should be set by the enforcer/inspector while sending requests
	// to sibyl.
	EntityTypeBot
	// EntityTypeAdmin represents an account which is considered as an admin
	// in a psychohazard event. it's completely up to the person who is scanning
	// to decide what is an admin account.
	EntityTypeAdmin
	// EntityTypeOwner represents an account which is considered as an owner
	// in a psychohazard event. it's completely up to the person who is scanning
	// to decide what is an owner account.
	EntityTypeOwner
	// EntityTypeChannel represents an entity which is considered as a channel.
	EntityTypeChannel
	// EntityTypeGroup represents an entity which is considered as a group.
	EntityTypeGroup
)

// flags constants
const (
	BanFlagTrolling     = "TROLLING"
	BanFlagSpam         = "SPAM"
	BanFlagEvade        = "EVADE"
	BanFlagCustom       = "CUSTOM"
	BanFlagPsychoHazard = "PSYCHOHAZARD"
	BanFlagMalImp       = "MALIMP"
	BanFlagNSFW         = "NSFW"
	BanFlagRaid         = "RAID"
	BanFlagSpamBot      = "SPAMBOT"
	BanFlagMassAdd      = "MASSADD"
)
