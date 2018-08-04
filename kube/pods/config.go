package pods

import (
	"encoding/json"

	"k8s.io/api/core/v1"
)

// Config wraps v1.Deployment so that it can be managed using configio
type Config struct {
	key string
	Pod *v1.Pod
}

func (conf *Config) Init(key string) *Config {
	conf.Pod = new(v1.Pod)
	conf.key = key
	conf.Pod.Kind = Kind
	conf.Pod.APIVersion = APIVersion
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
