//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"errors"

	"github.com/edgexfoundry/go-mod-configuration/v2/internal/pkg/keeper/utils/http"
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

func (c *Client) Ping() error {
	errResp := http.GetRequest(nil, c.baseUrl, ApiPingRoute, nil)
	if errResp.StatusCode != 0 {
		return errors.New(errResp.Message)
	}
	return nil
}
