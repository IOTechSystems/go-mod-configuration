//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/corekeeper"
	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/corekeeper/dtos"

	"github.com/mitchellh/mapstructure"
)

// decode converts the key-value pairs from core keeper to the target configuration data type
func decode(prefix string, pairs []dtos.KV, configTarget interface{}) error {
	raw := make(map[string]interface{})
	for _, p := range pairs {
		// Trim the prefix off our key first
		key := strings.TrimPrefix(p.Key, prefix)

		// Determine what map we're writing the value to. We split by '/'
		// to determine any sub-maps that need to be created.
		m := raw
		children := strings.Split(key, corekeeper.KeyDelimiter)
		if len(children) > 0 {
			key = children[len(children)-1]
			children = children[:len(children)-1]
			for _, child := range children {
				if m[child] == nil {
					m[child] = make(map[string]interface{})
				}

				subm, ok := m[child].(map[string]interface{})
				if !ok {
					return fmt.Errorf("child is both a data item and dir: %s", child)
				}

				m = subm
			}
		}

		switch p.Value.(type) {
		case bool:
			m[key] = p.Value.(bool)
		case int:
			m[key] = p.Value.(int)
		case int64:
			m[key] = p.Value.(int64)
		case float64:
			m[key] = p.Value.(float64)
		case string:
			m[key] = p.Value.(string)
		default:
			return errors.New("unknown data type of the stored value")
		}

	}

	// Now decode into it
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   configTarget,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(raw); err != nil {
		return err
	}

	return nil
}
