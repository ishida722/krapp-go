package models

import (
	"fmt"
	"time"

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

func (fm FrontMatter) SetCreatedNow() error {
	fm.SetCreated(time.Now())
	return nil
}

func (fm FrontMatter) SetCreated(t time.Time) error {
	if t.IsZero() {
		return fmt.Errorf("invalid time: zero value")
	}
	fm["created"] = t.Format("2006-01-02")
	return nil
}

func (fm FrontMatter) Created() (time.Time, error) {
	createdValue, ok := fm["created"]
	if !ok {
		return time.Time{}, fmt.Errorf("created field not found")
	}
	switch v := createdValue.(type) {
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse created date: %w", err)
		}
		return t, nil
	case time.Time:
		return v, nil
	default:
		return time.Time{}, fmt.Errorf("created field is not a valid type, expected string or time.Time")
	}
}

func (fm FrontMatter) Label() (string, error) {
	labelValue, ok := fm["label"]
	if !ok {
		return "", fmt.Errorf("label field not found")
	}
	switch v := labelValue.(type) {
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("label field is not a string")
	}
}

func (fm FrontMatter) Exists() bool {
	return len(fm) > 0
}
