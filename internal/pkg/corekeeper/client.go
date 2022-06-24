//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

import (
	"context"
	"errors"
	"fmt"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/corekeeper"
	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/corekeeper/dtos"
	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/types"

	"github.com/pelletier/go-toml"
)

type coreKeeperClient struct {
	keeperUrl       string
	keeperClient    *corekeeper.Client
	configBasePath  string
	watchingDoneCtx context.Context
	watchingDone    context.CancelFunc
	watchingWait    sync.WaitGroup
}

// NewCoreKeeperClient creates a new Core Keeper Client.
func NewCoreKeeperClient(config types.ServiceConfig) *coreKeeperClient {

	client := coreKeeperClient{
		keeperUrl:      config.GetUrl(),
		configBasePath: config.BasePath,
	}

	client.watchingDoneCtx, client.watchingDone = context.WithCancel(context.Background())

	client.createKeeperClient(client.keeperUrl)
	return &client
}

func (client *coreKeeperClient) fullPath(name string) string {
	return path.Join(client.configBasePath, name)
}

func (client *coreKeeperClient) createKeeperClient(url string) {
	client.keeperClient = corekeeper.NewClient(url)
}

// IsAlive simply checks if Core Keeper is up and running at the configured URL
func (client *coreKeeperClient) IsAlive() bool {
	return false
}

// HasConfiguration checks to see if Consul contains the service's configuration.
func (client *coreKeeperClient) HasConfiguration() (bool, error) {
	resp, err := client.keeperClient.KV().Keys(client.configBasePath)
	if err != nil {
		return false, fmt.Errorf("checking configuration existence from Core Keeper failed: %v", err)
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	return true, nil
}

func (client *coreKeeperClient) HasSubConfiguration(name string) (bool, error) {
	keyPath := client.fullPath(name)
	resp, err := client.keeperClient.KV().Keys(keyPath)
	if err != nil {
		return false, fmt.Errorf("checking configuration existence from Core Keeper failed: %v", err)
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	return true, nil
}

func (client *coreKeeperClient) PutConfigurationToml(configuration *toml.Tree, overwrite bool) error {
	return nil
}

func (client *coreKeeperClient) PutConfiguration(configStruct interface{}, overwrite bool) error {
	err := client.keeperClient.KV().Put(client.configBasePath, configStruct)
	if err != nil {
		return fmt.Errorf("unable to JSON marshal configStruct, err: %v", err)
	}
	return nil
}

func (client *coreKeeperClient) GetConfiguration(configStruct interface{}) (interface{}, error) {
	exists, err := client.HasConfiguration()
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("the Configuration service (EdgeX Keeper) doesn't contain configuration for %s", client.configBasePath)
	}

	resp, err := client.keeperClient.KV().Get(client.configBasePath)
	if err != nil {
		return nil, err
	}

	err = setStructFields(resp.KV, configStruct)
	if err != nil {
		return nil, err
	}
	return configStruct, nil
}

func setStructFields(kv []dtos.KV, configStruct interface{}) error {
	config := reflect.ValueOf(configStruct)

	var rootStruct, innerStructEle reflect.Value
	if config.Kind() == reflect.Ptr {
		rootStruct = config.Elem()
	} else {
		rootStruct = config
	}
	if rootStruct.Kind() != reflect.Struct {
		return fmt.Errorf("configuration value must be a pointer to a struct interface")
	}

	for _, item := range kv {
		paths := strings.Split(item.Key, corekeeper.KeyDelimiter)
		// reset the structElement to te top level
		innerStructEle = rootStruct
		for index, path := range paths[1:] {
			innerStructEle = innerStructEle.FieldByName(path)
			if index+1 == len(paths)-1 {
				if value, ok := item.Value.(string); ok {
					innerStructEle.SetString(value)
				} else if boolValue, ok := item.Value.(bool); ok {
					innerStructEle.SetBool(boolValue)
				} else if intValue, ok := item.Value.(int64); ok {
					innerStructEle.SetInt(intValue)
				} else if floatValue, ok := item.Value.(float64); ok {
					innerStructEle.SetFloat(floatValue)
				} else {
					return errors.New("unable to set config struct field value")
				}
			}
		}
	}
	return nil
}

func (client *coreKeeperClient) WatchForChanges(updateChannel chan<- interface{}, errorChannel chan<- error, configuration interface{}, waitKey string) {
}

func (client *coreKeeperClient) StopWatching() {}

func (client *coreKeeperClient) ConfigurationValueExists(name string) (bool, error) {
	keyPath := client.fullPath(name)
	resp, err := client.keeperClient.KV().Get(keyPath)
	if err != nil {
		return false, fmt.Errorf("checking configuration existence from Core Keeper failed: %v", err)
	}
	if resp.StatusCode == 404 {
		return false, nil
	}

	exists := false
	// traverse each key value pair from the response and get the key matched the fullPath
	for _, kvPair := range resp.KV {
		if kvPair.Key == keyPath {
			exists = kvPair.Value != nil
			break
		}
	}
	return exists, nil
}

func (client *coreKeeperClient) GetConfigurationValue(name string) ([]byte, error) {
	return []byte{}, nil
}

func (client *coreKeeperClient) PutConfigurationValue(name string, value []byte) error {
	err := client.keeperClient.KV().Put(client.configBasePath, value)
	if err != nil {
		return fmt.Errorf("unable to JSON marshal configStruct, err: %v", err)
	}
	return nil
}
