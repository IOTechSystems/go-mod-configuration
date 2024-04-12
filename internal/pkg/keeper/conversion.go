//
// Copyright (C) 2022-2024 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package keeper

import (
	"strconv"

	"github.com/edgexfoundry/go-mod-configuration/v3/internal/pkg/keeper/api"

	"github.com/spf13/cast"
)

type pair struct {
	Key   string
	Value string
}

func convertInterfaceToPairs(path string, interfaceMap interface{}) []*pair {
	pairs := make([]*pair, 0)

	pathPre := ""
	if path != "" {
		pathPre = path + api.KeyDelimiter
	}

	switch value := interfaceMap.(type) {
	case []interface{}:
		for index, item := range value {
			nextPairs := convertInterfaceToPairs(pathPre+strconv.Itoa(index), item)
			pairs = append(pairs, nextPairs...)
		}
	case map[string]interface{}:
		for index, item := range value {
			nextPairs := convertInterfaceToPairs(pathPre+index, item)
			pairs = append(pairs, nextPairs...)
		}
		if len(value) == 0 {
			pairs = append(pairs, &pair{Key: pathPre + "Placeholder", Value: ""})
		}
	default:
		pairs = append(pairs, &pair{Key: path, Value: cast.ToString(value)})
	}

	return pairs
}
