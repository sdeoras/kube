package ns

import (
	"encoding/json"

	"k8s.io/api/core/v1"
)

// Config wraps v1.Pod so that it can be managed using configio
type Config struct {
	key       string
	Namespace *v1.Namespace
}

func (conf *Config) Init(key string) *Config {
	conf.Namespace = new(v1.Namespace)
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