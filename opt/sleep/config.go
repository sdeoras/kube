package sleep

import (
	"encoding/json"
	"time"
)

// Config wraps v1.Deployment so that it can be managed using configio
type Config struct {
	Name          string
	SleepDuration time.Duration
}

func (conf *Config) Init(name string, dur time.Duration) *Config {
	conf.Name = name
	conf.SleepDuration = dur
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
