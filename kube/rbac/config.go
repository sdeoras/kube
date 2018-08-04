package rbac

import (
	"encoding/json"
	"k8s.io/api/rbac/v1beta1"
)

// Config wraps v1.ClusterRoleBinding so that it can be managed using configio
type Config struct {
	key                string
	ClusterRoleBinding *v1beta1.ClusterRoleBinding
}

func (conf *Config) Init(key string) *Config {
	conf.ClusterRoleBinding = new(v1beta1.ClusterRoleBinding)
	conf.key = key
	conf.ClusterRoleBinding.Kind = Kind
	conf.ClusterRoleBinding.APIVersion = APIVersion
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
