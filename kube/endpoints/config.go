package endpoints

import (
	"encoding/json"

	"k8s.io/api/core/v1"
)

// Config wraps v1.Endpoints so that it can be managed using configio
type Config struct {
	key       string
	Endpoints *v1.Endpoints
}

func (conf *Config) Init(key string) *Config {
	conf.Endpoints = new(v1.Endpoints)
	conf.key = key
	conf.Endpoints.Kind = Kind
	conf.Endpoints.APIVersion = APIVersion
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
