//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package keeper

import (
	"fmt"
	"strconv"

	"github.com/edgexfoundry/go-mod-configuration/v2/internal/pkg/keeper/api"
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
	case string:
		pairs = append(pairs, &pair{Key: path, Value: value})
	case int:
		pairs = append(pairs, &pair{Key: path, Value: strconv.Itoa(value)})
	case int8:
		value8 := int(value)
		pairs = append(pairs, &pair{Key: path, Value: strconv.Itoa(value8)})
	case int16:
		value16 := int(value)
		pairs = append(pairs, &pair{Key: path, Value: strconv.Itoa(value16)})
	case int32:
		value32 := int(value)
		pairs = append(pairs, &pair{Key: path, Value: strconv.Itoa(value32)})
	case int64:
		pairs = append(pairs, &pair{Key: path, Value: strconv.FormatInt(value, 10)})
	case float32:
		valueF64 := float64(value)
		pairs = append(pairs, &pair{Key: path, Value: strconv.FormatFloat(valueF64, 'g', -1, 32)})
	case float64:
		pairs = append(pairs, &pair{Key: path, Value: strconv.FormatFloat(value, 'g', -1, 64)})
	case bool:
		pairs = append(pairs, &pair{Key: path, Value: strconv.FormatBool(value)})
	case nil:
		pairs = append(pairs, &pair{Key: path, Value: ""})
	default:
		pairs = append(pairs, &pair{Key: path, Value: fmt.Sprintf("%v", value)})
	}

	return pairs
}
