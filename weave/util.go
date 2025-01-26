package weave

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
)

// parseData parses the input bytes as either JSON or YAML into a slice of maps
func parseData(dataBytes []byte) ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	// Try parsing as JSON
	if json.Unmarshal(dataBytes, &data) == nil {
		return data, nil
	}

	// Try parsing as YAML
	if yaml.Unmarshal(dataBytes, &data) == nil {
		return data, nil
	}

	return nil, errors.New("input is neither valid JSON nor YAML")
}
