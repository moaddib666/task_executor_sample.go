package tasks

import (
	"encoding/json"
	"fmt"
)

type Task struct {
	Name     string      `json:"name" yaml:"name"`
	Location string      `json:"location" yaml:"location"`
	Accept   *JsonSchema `json:"schema" yaml:"schema"`
}

func (t *Task) JsonSchema() string {
	schema, err := json.MarshalIndent(t.Accept, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return string(schema)
}

func (t *Task) SetSchemaFromBites(schema []byte) error {
	// Parse the schema from JSON into a map[string]interface{}
	if err := json.Unmarshal(schema, &t.Accept); err != nil {
		return fmt.Errorf("failed to parse schema: %s", err)
	}
	return nil
}

func (t *Task) ParseArgs(args interface{}) (map[string]interface{}, error) {
	var argsMap map[string]interface{}
	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("args must be a map[string]interface{}")
	}

	err := validateJSONAgainstSchema(argsMap, t.Accept)
	if err != nil {
		return nil, err
	}

	return argsMap, nil
}

func NewTask(name string, location string) *Task {
	return &Task{Name: name, Location: location, Accept: &JsonSchema{}}
}
