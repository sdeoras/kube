package jobs

import (
	"encoding/json"

	"k8s.io/api/batch/v1"
)

// Config wraps v1.Job so that it can be managed using configio
type Config struct {
	key string
	Job *v1.Job
}

func (conf *Config) Init(key string) *Config {
	conf.Job = new(v1.Job)
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
