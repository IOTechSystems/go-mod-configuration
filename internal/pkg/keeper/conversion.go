//
// Copyright (C) 2022 IOTech Ltd
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

func convertMapToKVPairs(path string, interfaceMap interface{}) []*pair {
	pairs := make([]*pair, 0)

	pathPre := ""
	if path != "" {
		pathPre = path + api.KeyDelimiter
	}

	switch value := interfaceMap.(type) {
	case []interface{}:
		for index, item := range value {
			nextPairs := convertMapToKVPairs(pathPre+strconv.Itoa(index), item)
			pairs = append(pairs, nextPairs...)
		}
	case map[string]interface{}:
		for index, item := range value {
			nextPairs := convertMapToKVPairs(pathPre+index, item)
			pairs = append(pairs, nextPairs...)
		}
	default:
		pairs = append(pairs, &pair{Key: path, Value: cast.ToString(value)})
	}

	return pairs
}

func convertInterfaceToPairs(path string, interfaceMap any) []*pair {
	pairs := make([]*pair, 0)

	pathPre := ""
	if path != "" {
		pathPre = path + "/"
	}

	switch value := interfaceMap.(type) {
	case []interface{}:
		for index, item := range value {
			nextPairs := convertInterfaceToPairs(pathPre+strconv.Itoa(index), item)
			pairs = append(pairs, nextPairs...)
		}

	case map[string]any:
		for index, item := range value {
			nextPairs := convertInterfaceToPairs(pathPre+index, item)
			pairs = append(pairs, nextPairs...)
		}

	case int:
		pairs = append(pairs, &pair{Key: path, Value: strconv.Itoa(value)})

	case int64:
		var value64 = int(value)
		pairs = append(pairs, &pair{Key: path, Value: strconv.Itoa(value64)})

	case float64:
		pairs = append(pairs, &pair{Key: path, Value: strconv.FormatFloat(value, 'f', -1, 64)})

	case bool:
		pairs = append(pairs, &pair{Key: path, Value: strconv.FormatBool(value)})

	case nil:
		pairs = append(pairs, &pair{Key: path, Value: ""})

	default:
		pairs = append(pairs, &pair{Key: path, Value: value.(string)})
	}

	return pairs
}
