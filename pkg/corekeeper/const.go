//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

import "regexp"

// new constants relates to EdgeX Keeper service and will be added to go-mod-core-contracts in the future
const CoreKeeperServiceKey = "core-keeper"

// KeyAllowedCharsRegexString defined the characters allowed in the key name
const KeyAllowedCharsRegexString = "^[a-zA-Z0-9-_~;=./]+$"

var (
	KeyAllowedCharsRegex = regexp.MustCompile(KeyAllowedCharsRegexString)
)

// key delimiter for edgex keeper
const KeyDelimiter = "/"

// Constants related to defined routes in the v2 service APIs
const ApiBase = "/api/v2"
const ApiKVRoute = ApiBase + "/kv"

// Constants related to defined url path names and parameters in the v2 service APIs
const (
	Flatten = "flatten"
	Keys    = "keys"
)
