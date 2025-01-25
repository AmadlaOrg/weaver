package weave

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"text/template"
)

type IWeave interface {
	Do(data string) string
}

type SWeave struct{}

func (s *SWeave) Do(data string, tmplPaths ...string) (string, error) {
	tmpl, err := template.ParseFiles(tmplPaths...)
	if err != nil {
		panic(err)
	}

}

func parse(input string) (map[string]interface{}, error) {
	// Try parsing as JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err == nil {
		return data, nil
	}

	// If JSON parsing fails, try parsing as YAML
	if err := yaml.Unmarshal([]byte(input), &data); err == nil {
		return data, nil
	}

	// Return an error if neither JSON nor YAML parsing succeeds
	return nil, fmt.Errorf("input is neither valid JSON nor YAML")
}
