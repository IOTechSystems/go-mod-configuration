//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

import (
	"strconv"
	"testing"
	"time"

	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/types"

	"github.com/pelletier/go-toml"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serviceName = "coreKeeperUnitTest"
	testHost    = "localhost"
	port        = 59883
	dummyConfig = "dummy"
)

type LoggingInfo struct {
	EnableRemote bool
	File         string
}

type TestConfig struct {
	Logging  LoggingInfo
	Port     int
	Host     string
	LogLevel string
}

func makeCoreKeeperClient(serviceName string) *coreKeeperClient {
	config := types.ServiceConfig{
		Host:     testHost,
		Port:     port,
		BasePath: serviceName,
	}

	client := NewCoreKeeperClient(config)
	return client
}

func getUniqueServiceName() string {
	return serviceName + strconv.Itoa(time.Now().Nanosecond())
}

func configValueExists(key string, client *coreKeeperClient) bool {
	exists, _ := client.ConfigurationValueExists(key)
	return exists
}

func TestHasConfigurationFalse(t *testing.T) {
	serviceName := getUniqueServiceName()
	client := makeCoreKeeperClient(serviceName)

	// Make sure the configuration doesn't already exists
	//reset(t, client)

	actual, err := client.HasConfiguration()
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	assert.False(t, actual)
}

func TestHasConfigurationTrue(t *testing.T) {
	serviceName := getUniqueServiceName()
	client := makeCoreKeeperClient(serviceName)

	// Make sure the configuration doesn't already exists
	//reset(t, client)

	actual, err := client.HasConfiguration()
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	assert.True(t, actual)
}

func TestHasSubConfigurationFalse(t *testing.T) {
	serviceName := getUniqueServiceName()
	client := makeCoreKeeperClient(serviceName)

	// Make sure the configuration doesn't already exists
	//reset(t, client)

	actual, err := client.HasSubConfiguration(dummyConfig)
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	assert.False(t, actual)
}

func TestHasSubConfigurationTrue(t *testing.T) {
	serviceName := getUniqueServiceName()
	client := makeCoreKeeperClient(serviceName)

	// Make sure the configuration doesn't already exists
	//reset(t, client)

	actual, err := client.HasSubConfiguration(dummyConfig)
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	assert.True(t, actual)
}

func createConfigMap() map[string]interface{} {
	configMap := make(map[string]interface{})

	configMap["int"] = 1
	configMap["int64"] = int64(64)
	configMap["float64"] = float64(1.4)
	configMap["string"] = "hello"
	configMap["bool"] = true
	configMap["nestedNode"] = map[string]interface{}{"field1": "value1", "field2": "value2"}

	return configMap
}

func TestPutConfigurationToml(t *testing.T) {
	client := makeCoreKeeperClient(getUniqueServiceName())
	configMap := createConfigMap()

	configToml, err := toml.TreeFromMap(configMap)
	if err != nil {
		t.Fatalf("unable to create TOML Tree from map: %v", err)
	}

	err = client.PutConfigurationToml(configToml, false)
	if !assert.NoError(t, err) {
		t.Fatal()
	}
}

func TestPutConfigurationTomlWithoutOverwrite(t *testing.T) {
	client := makeCoreKeeperClient("coreKeeperUnitTest321867000")
	configMap := createConfigMap()
	originConfig := configMap
	// overwrite the configMap fields
	configMap["nestedNode"] = map[string]interface{}{"field1": "xxx", "field2": "yyy"}
	configToml, err := toml.TreeFromMap(configMap)
	if err != nil {
		t.Fatalf("unable to create TOML Tree from map: %v", err)
	}

	err = client.PutConfigurationToml(configToml, false)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	actual, err := client.GetConfigurationValue("nestedNode/field1")
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	expected := []byte(originConfig["nestedNode"].(map[string]interface{})["field1"].(string))
	if !assert.Equal(t, expected, actual, "Values for %s are not equal, expected equal", "nestedNode/field1") {
		t.Fatal()
	}
}

func TestPutConfigurationTomlWithOverwrite(t *testing.T) {
	client := makeCoreKeeperClient("coreKeeperUnitTest321867000")
	configMap := createConfigMap()
	originConfig := configMap
	// overwrite the configMap fields
	configMap["nestedNode"] = map[string]interface{}{"field1": "value1", "field2": "value2"}
	configToml, err := toml.TreeFromMap(configMap)
	if err != nil {
		t.Fatalf("unable to create TOML Tree from map: %v", err)
	}

	err = client.PutConfigurationToml(configToml, true)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	actual, err := client.GetConfigurationValue("nestedNode/field1")
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	expected := []byte(originConfig["nestedNode"].(map[string]interface{})["field1"].(string))
	if !assert.NotEqual(t, expected, actual, "Values for %s are equal, expected not equal", "nestedNode/field1") {
		t.Fatal()
	}
}

func TestPutConfiguration(t *testing.T) {
	expected := TestConfig{
		Logging: LoggingInfo{
			EnableRemote: true,
			File:         "NONE",
		},
		Port:     8000,
		Host:     "localhost",
		LogLevel: "debug",
	}

	client := makeCoreKeeperClient(getUniqueServiceName())

	err := client.PutConfiguration(expected, true)
	if !assert.NoErrorf(t, err, "unable to put configuration: %v", err) {
		t.Fatal()
	}

	actual, err := client.HasConfiguration()
	require.NoError(t, err)
	if !assert.True(t, actual, "Failed to put configuration") {
		t.Fail()
	}

	assert.True(t, configValueExists("Logging/EnableRemote", client))
	assert.True(t, configValueExists("Logging/File", client))
	assert.True(t, configValueExists("Port", client))
	assert.True(t, configValueExists("Host", client))
	assert.True(t, configValueExists("LogLevel", client))
}

type ConfigurationStruct struct {
	Writable WritableInfo
}
type WritableInfo struct {
	LogLevel        string
	InsecureSecrets InsecureSecrets
	Telemetry       TelemetryInfo
}

// InsecureSecrets is used to hold the secrets stored in the configuration
type InsecureSecrets map[string]InsecureSecretsInfo

// InsecureSecretsInfo encapsulates info used to retrieve insecure secrets
type InsecureSecretsInfo struct {
	Path    string
	Secrets map[string]string
}

// TelemetryInfo contains the configuration for a service's metrics collection
type TelemetryInfo struct {
	// Interval is the time duration in which to collect and report the service's metrics
	Interval string
	// PublishTopicPrefix is the base topic in which to publish (report) the service's metrics to the EdgeX MessageBus
	// The service name and the metric name are appended to this base topic. i.e. <prefix>/<service-name>/<metric-name>
	PublishTopicPrefix string
	// Metrics is the list of service's metrics that can be collected. Each of the service's metrics must be in the list
	// and set to true if enable or false if disabled.
	Metrics map[string]bool
	// Tags is a list of service level tags that are attached to every metric reported for the service
	// Example: Gateway = "Gateway123"
	Tags map[string]string
}

func TestGetConfiguration(t *testing.T) {
	mockServiceName := getUniqueServiceName()
	client := makeCoreKeeperClient(mockServiceName)

	mockConf := ConfigurationStruct{
		Writable: WritableInfo{
			LogLevel: "INFO",
			InsecureSecrets: map[string]InsecureSecretsInfo{
				"DB": {
					Path:    "redisdb",
					Secrets: map[string]string{"username": "xxx", "password": "yyy"},
				},
			},
		},
	}

	err := client.PutConfiguration(mockConf, true)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	result, err := client.GetConfiguration(&ConfigurationStruct{})

	if !assert.NoError(t, err) {
		t.Fatal()
	}

	configuration := (result).(*ConfigurationStruct)

	if !assert.NotNil(t, configuration) {
		t.Fatal()
	}
}

func TestGetConfigurationValue(t *testing.T) {
	key := "Foo"
	expected := []byte("bar")

	uniqueServiceName := getUniqueServiceName()
	client := makeCoreKeeperClient(uniqueServiceName)

	err := client.PutConfigurationValue(key, expected)
	assert.NoError(t, err)

	actual, err := client.GetConfigurationValue(key)
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	if !assert.Equal(t, expected, actual) {
		t.Fatal()
	}
}

func TestPutConfigurationValue(t *testing.T) {
	key := "Foo"
	expected := []byte("bar")

	uniqueServiceName := getUniqueServiceName()

	client := makeCoreKeeperClient(uniqueServiceName)

	// Make sure the configuration doesn't already exists
	// reset(t, client)

	err := client.PutConfigurationValue(key, expected)
	assert.NoError(t, err)

	resp, err := client.keeperClient.KV().Get(client.fullPath(key))
	if !assert.NoError(t, err) {
		t.Fatal()
	}
	if !assert.NotNil(t, resp, "%s value not found", key) {
		t.Fatal()
	}

	actual := []byte(resp.KV[0].Value.(string))

	assert.Equal(t, expected, actual)
}
