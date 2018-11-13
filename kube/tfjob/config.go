package tfjob

import (
	"encoding/json"

	"github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
)

// Config wraps v1.Svc so that it can be managed using configio
type Config struct {
	key string
	Job *v1alpha2.TFJob
}

func (conf *Config) Init(key string) *Config {
	conf.Job = new(v1alpha2.TFJob)
	conf.key = key
	conf.Job.Kind = Kind
	conf.Job.APIVersion = APIVersion
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
