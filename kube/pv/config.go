package pv

import (
	"encoding/json"

	"k8s.io/api/core/v1"
)

// Config wraps v1.PersistentVolumeClaim so that it can be managed using configio
type Config struct {
	key              string
	PersistentVolume *v1.PersistentVolume
}

func (conf *Config) Init(key string) *Config {
	conf.PersistentVolume = new(v1.PersistentVolume)
	conf.key = key
	conf.PersistentVolume.Kind = Kind
	conf.PersistentVolume.APIVersion = APIVersion
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
