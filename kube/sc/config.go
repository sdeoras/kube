package sc

import (
	"encoding/json"
	"k8s.io/api/storage/v1beta1"
)

// Config wraps v1.ClusterRoleBinding so that it can be managed using configio
type Config struct {
	key          string
	StorageClass *v1beta1.StorageClass
}

func (conf *Config) Init(key string) *Config {
	conf.StorageClass = new(v1beta1.StorageClass)
	conf.key = key
	conf.StorageClass.Kind = Kind
	conf.StorageClass.APIVersion = APIVersion
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
