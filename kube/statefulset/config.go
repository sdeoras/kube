package statefulset

import (
	"encoding/json"

	"k8s.io/api/apps/v1beta1"
)

// Config wraps v1.ClusterRoleBinding so that it can be managed using configio
type Config struct {
	key         string
	StatefulSet *v1beta1.StatefulSet
}

func (conf *Config) Init(key string) *Config {
	conf.StatefulSet = new(v1beta1.StatefulSet)
	conf.key = key
	conf.StatefulSet.Kind = Kind
	conf.StatefulSet.APIVersion = APIVersion
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
