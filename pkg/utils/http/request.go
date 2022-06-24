//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// GetRequest makes the get request and return the body
func GetRequest(returnValuePointer interface{}, baseUrl string, requestPath string, requestParams url.Values) error {
	req, err := createRequest(http.MethodGet, baseUrl, requestPath, requestParams)
	if err != nil {
		return err
	}

	res, err := sendRequest(req)
	if err != nil {
		return (err)
	}
	// Check the response content length to avoid json unmarshal error
	if len(res) == 0 {
		return nil
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.New("failed to parse the response body")
	}
	return nil
}

// PutRequest makes the put JSON request and return the body
func PutRequest(
	returnValuePointer interface{},
	baseUrl string, requestPath string,
	requestParams url.Values,
	data interface{}) error {

	req, err := createRequestWithRawData(http.MethodPut, baseUrl, requestPath, requestParams, data)
	if err != nil {
		return err
	}

	res, err := sendRequest(req)
	if err != nil {
		return err
	}
	// Get a 204 no content reply
	if len(res) == 0 {
		return nil
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.New("failed to parse the response body")
	}
	return nil
}
