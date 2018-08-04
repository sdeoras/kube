package deployment

import (
	"encoding/json"

	"k8s.io/api/apps/v1"
)

// Config wraps v1.Deployment so that it can be managed using configio
type Config struct {
	key        string
	Deployment *v1.Deployment
}

func (conf *Config) Init(key string) *Config {
	conf.Deployment = new(v1.Deployment)
	conf.key = key
	conf.Deployment.Kind = Kind
	conf.Deployment.APIVersion = APIVersion
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
