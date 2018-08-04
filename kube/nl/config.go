package nl

import (
	"encoding/json"

	"k8s.io/api/core/v1"
)

// Config wraps v1.ConfigMap so that it can be managed using configio
type Config struct {
	key      string
	NodeList *v1.NodeList
}

func (conf *Config) Init(key string) *Config {
	conf.NodeList = new(v1.NodeList)
	conf.key = key
	conf.NodeList.Kind = Kind
	conf.NodeList.APIVersion = APIVersion
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
