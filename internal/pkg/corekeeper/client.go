//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package corekeeper

import (
	"context"
	"fmt"
	"path"
	"sync"

	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/corekeeper"
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
	_, err := client.keeperClient.Ping()
	if err != nil {
		return false
	}

	return true
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

// PutConfigurationToml puts a full toml configuration into Core Keeper
func (client *coreKeeperClient) PutConfigurationToml(configuration *toml.Tree, overwrite bool) error {
	configurationMap := configuration.ToMap()
	err := client.PutConfiguration(configurationMap, overwrite)
	if err != nil {
		return err
	}
	return nil
}

func (client *coreKeeperClient) PutConfiguration(config interface{}, overwrite bool) error {
	var err error
	if overwrite {
		err = client.keeperClient.KV().Put(client.configBasePath, config)
	} else {
		kvPairs := convertMapToKVPairs("", config)
		for _, kv := range kvPairs {
			exists, err := client.ConfigurationValueExists(kv.Key)
			if err != nil {
				return err
			}
			if !exists {
				// Only create the key if not exists in core keeper
				if err = client.PutConfigurationValue(kv.Key, []byte(kv.Value)); err != nil {
					return err
				}
			}
		}
	}
	if err != nil {
		return fmt.Errorf("error occurred while creating/updating configuration, error: %v", err)
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

	err = decode(client.configBasePath+corekeeper.KeyDelimiter, resp.KV, configStruct)
	if err != nil {
		return nil, err
	}
	return configStruct, nil
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
	keyPath := client.fullPath(name)
	resp, err := client.keeperClient.KV().Get(keyPath)
	if err != nil {
		return nil, err
	}
	if len(resp.KV) == 0 {
		return nil, fmt.Errorf("%s configuration not found", name)
	}
	value := resp.KV[0].Value.(string)
	return []byte(value), nil
}

func (client *coreKeeperClient) PutConfigurationValue(name string, value []byte) error {
	keyPath := client.fullPath(name)
	err := client.keeperClient.KV().Put(keyPath, value)
	if err != nil {
		return fmt.Errorf("unable to JSON marshal configStruct, err: %v", err)
	}
	return nil
}
