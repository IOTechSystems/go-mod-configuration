//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetRequest makes the get request and return the body
func GetRequest(returnValuePointer interface{}, baseUrl string, requestPath string, requestParams url.Values) (ErrorResponse, error) {
	req, err := createRequest(http.MethodGet, baseUrl, requestPath, requestParams)
	if err != nil {
		return ErrorResponse{}, err
	}

	res, errResp, err := sendRequest(req)
	if err != nil {
		return ErrorResponse{}, err
	}
	if errResp.StatusCode != 0 {
		return errResp, nil
	}
	// Check the response content length to avoid json unmarshal error
	if returnValuePointer == nil || len(res) == 0 {
		return ErrorResponse{}, nil
	}

	if unmarshalErr := json.Unmarshal(res, &returnValuePointer); unmarshalErr != nil {
		return ErrorResponse{}, fmt.Errorf("failed to parse the response body, err: %s", unmarshalErr.Error())
	}
	return ErrorResponse{}, nil
}

// PutRequest makes the put JSON request and return the body
func PutRequest(
	returnValuePointer interface{},
	baseUrl string, requestPath string,
	requestParams url.Values,
	data interface{}) (ErrorResponse, error) {

	req, err := createRequestWithRawData(http.MethodPut, baseUrl, requestPath, requestParams, data)
	if err != nil {
		return ErrorResponse{}, err
	}

	res, errResp, err := sendRequest(req)
	if err != nil {
		return ErrorResponse{}, err
	}
	if errResp.StatusCode != 0 {
		return errResp, nil
	}
	// no need to unmarshal the response if returnValuePointer is nil
	if returnValuePointer == nil {
		return ErrorResponse{}, nil
	}
	if unmarshalErr := json.Unmarshal(res, returnValuePointer); unmarshalErr != nil {
		return ErrorResponse{}, fmt.Errorf("failed to parse the response body, err: %s", unmarshalErr.Error())
	}
	return ErrorResponse{}, nil
}

// DeleteRequest makes the delete request and return the body
func DeleteRequest(returnValuePointer interface{}, baseUrl string, requestPath string, requestParams url.Values) (ErrorResponse, error) {
	req, err := createRequest(http.MethodDelete, baseUrl, requestPath, requestParams)
	if err != nil {
		return ErrorResponse{}, err
	}

	res, errResp, err := sendRequest(req)
	if err != nil {
		return ErrorResponse{}, err
	}
	if errResp.StatusCode != 0 {
		return errResp, nil
	}
	// Check the response content length to avoid json unmarshal error
	if returnValuePointer == nil || len(res) == 0 {
		return errResp, nil
	}

	if unmarshalErr := json.Unmarshal(res, &returnValuePointer); unmarshalErr != nil {
		return ErrorResponse{}, fmt.Errorf("failed to parse the response body, err: %s", unmarshalErr.Error())
	}
	return errResp, nil
}
