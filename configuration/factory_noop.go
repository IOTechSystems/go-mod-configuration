//go:build no_consul

package configuration

import (
	"github.com/pelletier/go-toml"

	"github.com/edgexfoundry/go-mod-configuration/v2/pkg/types"
)

type configClient struct{}

func NewConfigurationClient(config types.ServiceConfig) (Client, error) {
	return configClient{}, nil
}

func (c configClient) HasConfiguration() (bool, error) {
	return false, nil
}

func (c configClient) HasSubConfiguration(name string) (bool, error) {
	return false, nil
}

func (c configClient) PutConfigurationToml(configuration *toml.Tree, overwrite bool) error {
	return nil
}

func (c configClient) PutConfiguration(configStruct interface{}, overwrite bool) error {
	return nil
}

func (c configClient) GetConfiguration(configStruct interface{}) (interface{}, error) {
	return nil, nil
}

func (c configClient) WatchForChanges(updateChannel chan<- interface{}, errorChannel chan<- error, configuration interface{}, waitKey string) {
}

func (c configClient) StopWatching() {
}

func (c configClient) IsAlive() bool {
	return false
}

func (c configClient) ConfigurationValueExists(name string) (bool, error) {
	return false, nil
}

func (c configClient) GetConfigurationValue(name string) ([]byte, error) {
	return nil, nil
}

func (c configClient) PutConfigurationValue(name string, value []byte) error {
	return nil
}
