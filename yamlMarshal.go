package kube

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

func YAMLMarshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nil, nil
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return yaml.Marshal(m)
}
