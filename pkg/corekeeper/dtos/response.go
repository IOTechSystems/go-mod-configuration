//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

type KV struct {
	Key string `json:"key,omitempty"`
	StoredData
}

type StoredData struct {
	DBTimestamp
	Value interface{} `json:"value,omitempty"`
}

type KeyOnly string

func (kv *KV) SetKey(newKey string) {
	kv.Key = newKey
}

func (key *KeyOnly) SetKey(newKey string) {
	*key = KeyOnly(newKey)
}

type DBTimestamp struct {
	Created  int64 `json:"created,omitempty"`
	Modified int64 `json:"modified,omitempty"`
}

type MultiKVResponse struct {
	BaseResponse `json:",inline"`
	KV           []KV `json:"response"`
}

type MultiKeyResponse struct {
	BaseResponse `json:",inline"`
	Key          []KeyOnly `json:"response"`
}

// BaseResponse defines the base content for response DTOs (data transfer objects).
// This object and its properties correspond to the BaseResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.1.0#/BaseResponse
type BaseResponse struct {
	Versionable `json:",inline"`
	RequestId   string `json:"requestId,omitempty"`
	Message     string `json:"message,omitempty"`
	StatusCode  int    `json:"statusCode"`
}
