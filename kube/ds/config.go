package ds

import (
	"encoding/json"

	apps_v1beta2 "k8s.io/api/apps/v1beta2"
)

// Config wraps v1.DaemonSet so that it can be managed using configio
type Config struct {
	key       string
	DaemonSet *apps_v1beta2.DaemonSet
}

func (conf *Config) Init(key string) *Config {
	conf.DaemonSet = new(apps_v1beta2.DaemonSet)
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
