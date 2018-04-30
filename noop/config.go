package noop

import (
	"encoding/json"
)

// Config wraps v1.Pod so that it can be managed using configio
type Config struct {
	Name  string
	Value int
}

func (conf *Config) Init(name string) *Config {
	conf.Name = name
	return conf
}

func (conf *Config) Key() string {
	return conf.Name
}

func (conf *Config) Marshal() ([]byte, error) {
	return json.MarshalIndent(conf, "", "  ")
}

func (conf *Config) Unmarshal(b []byte) error {
	return json.Unmarshal(b, conf)
}
