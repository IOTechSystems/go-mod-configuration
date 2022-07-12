//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package models

// MessageBusInfo provides parameters related to connecting to a message bus as a publisher
// copied from go-mod-bootstrap/config/types.go
type MessageBusInfo struct {
	// Indicates the message bus implementation to use, i.e. zero, mqtt, redisstreams...
	Type string
	// Protocol indicates the protocol to use when accessing the message bus.
	Protocol string
	// Host is the hostname or IP address of the broker, if applicable.
	Host string
	// Port defines the port on which to access the message bus.
	Port int
	// PublishTopicPrefix indicates the topic prefix the data is published to.
	PublishTopicPrefix string
	// SubscribeTopic indicates the topic in which to subscribe.
	SubscribeTopic string
	// AuthMode specifies the type of secure connection to the message bus which are 'none', 'usernamepassword'
	// 'clientcert' or 'cacert'. Not all option supported by each implementation.
	// ZMQ doesn't support any Authmode beyond 'none', RedisStreams only supports 'none' & 'usernamepassword'
	// while MQTT supports all options.
	AuthMode string
	// SecretName is the name of the secret in the SecretStore that contains the Auth Credentials. The credential are
	// dynamically loaded using this name and store the Option property below where the implementation expected to
	// find them.
	SecretName string
	// Provides additional configuration properties which do not fit within the existing field.
	// Typically the key is the name of the configuration property and the value is a string representation of the
	// desired value for the configuration property.
	Optional map[string]string
	// SubscribeEnabled indicates whether enable the subscription to the Message Queue
	SubscribeEnabled bool
}

// InsecureSecrets is used to hold the secrets stored in the configuration
// copied from go-mod-bootstrap/config/types.go
type InsecureSecrets map[string]InsecureSecretsInfo

// InsecureSecretsInfo encapsulates info used to retrieve insecure secrets
// copied from go-mod-bootstrap/config/types.go
type InsecureSecretsInfo struct {
	Path    string
	Secrets map[string]string
}

// TelemetryInfo contains the configuration for a service's metrics collection
// copied from go-mod-bootstrap/config/types.go
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

type WritableInfo struct {
	PersistData     bool
	LogLevel        string
	InsecureSecrets InsecureSecrets
	Telemetry       TelemetryInfo
}

type ConfigurationStruct struct {
	Writable     WritableInfo
	MessageQueue MessageBusInfo
}
