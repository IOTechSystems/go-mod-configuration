//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

import (
	"fmt"
	"net/url"
	"path"

	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/corekeeper/dtos"
	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/utils/http"
)

// KV is used to manipulate the K/V API
type KV struct {
	c *Client
}

// KV is used to return a handle to the K/V apis
func (c *Client) KV() *KV {
	return &KV{c}
}

// Get is used to lookup a single key. The returned pointer
// to the KVPair will be nil if the key does not exist.
func (k *KV) Get(key string) (res dtos.MultiKVResponse, err error) {
	url := fmt.Sprintf(ApiKVRoute+"/%s", key)
	err = http.GetRequest(&res, k.c.baseUrl, url, nil)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (k *KV) Keys(key string) (res dtos.MultiKeyResponse, err error) {
	pathParams := url.Values{}
	pathParams.Add(Keys, "true")

	url := fmt.Sprintf(ApiKVRoute+"/%s", key)
	err = http.GetRequest(&res, k.c.baseUrl, url, pathParams)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (k *KV) Put(key string, data interface{}) (err error) {
	keyPath := path.Join(ApiKVRoute, key)
	urlParams := url.Values{}
	urlParams.Add(Flatten, "true")
	request := dtos.AddKeysRequest{
		Value: data,
	}
	err = http.PutRequest(nil, k.c.baseUrl, keyPath, urlParams, request)
	if err != nil {
		return err
	}
	return
}
