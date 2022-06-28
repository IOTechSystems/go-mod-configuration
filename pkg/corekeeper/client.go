//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

import (
	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/corekeeper/dtos"
	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/utils/http"
)

type Client struct {
	baseUrl string
}

// NewClient creates an instance of Client
func NewClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
	}
}

func (c *Client) Ping() (res dtos.PingResponse, err error) {
	err = http.GetRequest(&res, c.baseUrl, ApiPingRoute, nil)
	if err != nil {
		return res, err
	}
	return res, nil
}
