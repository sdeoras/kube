package netpolicy

import (
	"encoding/json"

	v1 "k8s.io/api/networking/v1"

	"strings"
)

type Config struct {
	key           string
	NetworkPolicy *v1.NetworkPolicy
}

func (conf *Config) Init(key, namespace string) *Config {
	conf.NetworkPolicy = new(v1.NetworkPolicy)
	conf.NetworkPolicy.Namespace = namespace
	conf.NetworkPolicy.Name = strings.ToLower(key) + "-name"
	conf.NetworkPolicy.Kind = Kind
	conf.NetworkPolicy.APIVersion = ApiVersion
	conf.key = key
	return conf
}

func (conf *Config) Key() string {
	return conf.key
}

func (conf *Config) Marshal() ([]byte, error) {
	return json.MarshalIndent(conf, "", "  ")
}

func (conf *Config) Unmarshal(b []byte) error {
	return json.Unmarshal(b, conf)
}
