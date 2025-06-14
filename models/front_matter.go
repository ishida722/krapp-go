package models

import (
	"gopkg.in/yaml.v3"
)

type FrontMatter map[string]any

func (fm FrontMatter) ToYAML() (string, error) {
	bytes, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}
	return "---\n" + string(bytes) + "---\n", nil
}

func (fm FrontMatter) Exists() bool {
	return len(fm) > 0
}
