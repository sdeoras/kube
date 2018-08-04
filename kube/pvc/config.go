package pvc

import (
	"encoding/json"

	"k8s.io/api/core/v1"
)

// Config wraps v1.PersistentVolumeClaim so that it can be managed using configio
type Config struct {
	key                   string
	PersistentVolumeClaim *v1.PersistentVolumeClaim
}

func (conf *Config) Init(key string) *Config {
	conf.PersistentVolumeClaim = new(v1.PersistentVolumeClaim)
	conf.PersistentVolumeClaim.Kind = Kind
	conf.PersistentVolumeClaim.APIVersion = APIVersion
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
