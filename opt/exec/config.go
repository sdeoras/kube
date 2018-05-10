package exec

import (
	"encoding/json"
	"io"
)

// Config wraps a command description so that it can be managed using configio
type Config struct {
	key    string
	Cmd    string
	Args   []string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (conf *Config) Init(key string, stdin io.Reader, stdout, stderr io.Writer, cmd string, args ...string) *Config {
	conf.key = key
	conf.stdin = stdin
	conf.stdout = stdout
	conf.stderr = stderr
	conf.Cmd = cmd
	conf.Args = args
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
