//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

type Client struct {
	baseUrl string
}

// NewClient creates an instance of Client
func NewClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
	}
}
