//
// Copyright (C) 2022-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/spf13/cast"

	"github.com/edgexfoundry/go-mod-configuration/v3/internal/pkg/keeper/api"
	"github.com/edgexfoundry/go-mod-configuration/v3/internal/pkg/keeper/dtos"
	"github.com/edgexfoundry/go-mod-configuration/v3/internal/pkg/keeper/utils/http"
	"github.com/edgexfoundry/go-mod-configuration/v3/pkg/types"

	"github.com/edgexfoundry/go-mod-messaging/v3/messaging"
	msgTypes "github.com/edgexfoundry/go-mod-messaging/v3/pkg/types"
)

const (
	keeperTopicPrefix = "edgex/configs"
)

type keeperClient struct {
	keeperUrl      string
	keeperClient   *api.Caller
	configBasePath string
	watchingDone   chan bool
}

// NewKeeperClient creates a new Keeper Client.
func NewKeeperClient(config types.ServiceConfig) *keeperClient {
	client := keeperClient{
		keeperUrl:      config.GetUrl(),
		configBasePath: config.BasePath,
		watchingDone:   make(chan bool, 1),
	}

	client.createKeeperClient(client.keeperUrl)
	return &client
}

func (client *keeperClient) fullPath(name string) string {
	return path.Join(client.configBasePath, name)
}

func (client *keeperClient) createKeeperClient(url string) {
	client.keeperClient = api.NewCaller(url)
}

// IsAlive simply checks if Core Keeper is up and running at the configured URL
func (client *keeperClient) IsAlive() bool {
	err := client.keeperClient.Ping()
	return err == nil
}

// HasConfiguration checks to see if Core Keeper contains the service's configuration.
func (client *keeperClient) HasConfiguration() (bool, error) {
	resp, err := client.keeperClient.KV().Keys(client.configBasePath)
	if err != nil {
		return false, fmt.Errorf("checking configuration existence from Core Keeper failed: %v", err)
	}
	if len(resp.Keys) == 0 {
		return false, nil
	}
	return true, nil
}

// HasSubConfiguration checks to see if the Configuration service contains the service's sub configuration.
func (client *keeperClient) HasSubConfiguration(name string) (bool, error) {
	keyPath := client.fullPath(name)
	resp, err := client.keeperClient.KV().Keys(keyPath)
	if err != nil {
		return false, fmt.Errorf("checking configuration existence from Core Keeper failed: %v", err)
	}
	if len(resp.Keys) == 0 {
		return false, nil
	}
	return true, nil
}

// PutConfigurationMap puts a full configuration map into Core Keeper.
// The sub-paths to where the values are to be stored in Core Keeper are generated from the map key.
func (client *keeperClient) PutConfigurationMap(configuration map[string]any, overwrite bool) error {

	keyValues := convertInterfaceToPairs("", configuration)

	// Put config properties into Core Keeper.
	for _, keyValue := range keyValues {
		exists, _ := client.ConfigurationValueExists(keyValue.Key)
		if !exists || overwrite {
			if err := client.PutConfigurationValue(keyValue.Key, []byte(keyValue.Value)); err != nil {
				return err
			}
		}
	}

	return nil
}

// PutConfiguration puts a full configuration struct into the Configuration provider
func (client *keeperClient) PutConfiguration(config interface{}, overwrite bool) error {
	var err error
	if overwrite {
		err = client.keeperClient.KV().PutKeys(client.configBasePath, config)
	} else {
		kvPairs := convertInterfaceToPairs("", config)
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

func (client *keeperClient) GetConfiguration(configStruct interface{}) (interface{}, error) {
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

	err = decode(client.configBasePath+api.KeyDelimiter, resp.KVs, configStruct)
	if err != nil {
		return nil, err
	}
	return configStruct, nil
}

func (client *keeperClient) WatchForChanges(updateChannel chan<- interface{}, errorChannel chan<- error, configuration interface{}, waitKey string, messageBus messaging.MessageClient) {
	if messageBus == nil {
		configErr := errors.New("unable to use MessageClient to watch for configuration changes")
		errorChannel <- configErr
		return
	}

	messages := make(chan msgTypes.MessageEnvelope)
	topic := path.Join(keeperTopicPrefix, client.configBasePath, waitKey, "#")
	topics := []msgTypes.TopicChannel{
		{
			Topic:    topic,
			Messages: messages,
		},
	}

	watchErrors := make(chan error)
	err := messageBus.Subscribe(topics, watchErrors)
	if err != nil {
		_ = messageBus.Disconnect()
		errorChannel <- err
		return
	}

	go func() {
		defer func() {
			_ = messageBus.Disconnect()
		}()

		// send a nil value to updateChannel once the watcher connection is established
		// for go-mod-bootstrap to ignore the first change event
		// refer to the isFirstUpdate variable declared in https://github.com/edgexfoundry/go-mod-bootstrap/blob/main/bootstrap/config/config.go
		updateChannel <- nil

	outerLoop:
		for {
			select {
			case <-client.watchingDone:
				return
			case e := <-watchErrors:
				errorChannel <- e
			case msgEnvelope := <-messages:
				if msgEnvelope.ContentType != http.ContentTypeJSON {
					continue
				}
				var updatedConfig dtos.KV
				// unmarshal the updated config to KV DTO
				err := json.Unmarshal(msgEnvelope.Payload, &updatedConfig)
				if err != nil {
					continue
				}
				keyPrefix := path.Join(client.configBasePath, waitKey)

				// get the whole configs KV DTO array from Keeper with the same keyPrefix
				kvConfigs, err := client.keeperClient.KV().Get(keyPrefix)
				if err != nil {
					continue
				}

				// if the updated key not equal to keyPrefix, need to check the updated key and value from the message payload are valid
				// e.g. keyPrefix = "edgex/core/2.0/core-data/Writable" which is the root level of Writable configuration
				if updatedConfig.Key != keyPrefix {
					foundUpdatedKey := false
					for _, c := range kvConfigs.KVs {
						if c.Key == updatedConfig.Key {
							// the updated key from the message payload has been found in Keeper
							foundUpdatedKey = true
							// if the updated value in the message payload is different from the one obtained by Keeper
							// skip this subscribed message payload and continue the outer loop
							if c.Value != updatedConfig.Value {
								continue outerLoop
							}
							break
						}
					}
					// if the updated key from the message payload hasn't been found in Keeper
					// skip this subscribed message payload
					if !foundUpdatedKey {
						continue
					}
				}

				// decode KV DTO array to configuration struct
				err = decode(keyPrefix, kvConfigs.KVs, configuration)
				if err != nil {
					continue
				}
				updateChannel <- configuration
			}
		}
	}()
}

func (client *keeperClient) StopWatching() {
	client.watchingDone <- true
}

func (client *keeperClient) ConfigurationValueExists(name string) (bool, error) {
	keyPath := client.fullPath(name)
	res, err := client.keeperClient.KV().Keys(keyPath)
	if err != nil {
		return false, fmt.Errorf("checking configuration existence from Core Keeper failed: %v", err)
	}
	if len(res.Keys) == 0 {
		return false, nil
	}
	return true, nil
}

func (client *keeperClient) GetConfigurationValue(name string) ([]byte, error) {
	keyPath := client.fullPath(name)
	return client.GetConfigurationValueByFullPath(keyPath)
}

// GetConfigurationValueByFullPath gets a specific configuration value given the full path from Core Keeper
func (client *keeperClient) GetConfigurationValueByFullPath(name string) ([]byte, error) {
	resp, err := client.keeperClient.KV().Get(name)
	if err != nil {
		return nil, err
	}
	if len(resp.KVs) == 0 {
		return nil, fmt.Errorf("%s configuration not found", name)
	}

	valueStr := cast.ToString(resp.KVs[0].Value)

	return []byte(valueStr), nil
}

func (client *keeperClient) PutConfigurationValue(name string, value []byte) error {
	keyPath := client.fullPath(name)
	err := client.keeperClient.KV().Put(keyPath, value)
	if err != nil {
		return fmt.Errorf("unable to JSON marshal configStruct, err: %v", err)
	}
	return nil
}

func (client *keeperClient) GetConfigurationKeys(name string) ([]string, error) {
	keyPath := client.fullPath(name)
	resp, err := client.keeperClient.KV().Keys(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to get list of keys for %s from Core Keeper: %v", keyPath, err)
	}

	var list []string
	for _, v := range resp.Keys {
		list = append(list, string(v))
	}
	return list, nil
}
