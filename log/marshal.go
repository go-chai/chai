package log

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)

func MarshalYAML(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
