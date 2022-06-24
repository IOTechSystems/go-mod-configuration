//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

// BaseRequest defines the base content for request DTOs (data transfer objects).
// This object and its properties correspond to the BaseRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.1.0#/BaseRequest
type BaseRequest struct {
	Versionable `json:",inline"`
	RequestId   string `json:"requestId" validate:"len=0|uuid"`
}

// AddKeysRequest defines the Request Content for POST Key DTO.
type AddKeysRequest struct {
	BaseRequest `json:",inline"`
	Value       interface{} `json:"value,omitempty"`
}
