//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

// Versionable shows the API version in DTOs
type Versionable struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion" validate:"required"`
}
