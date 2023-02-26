package tasks

import "encoding/json"

type ScriptSchema struct {
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`               // optional
	Description string      `json:"description,omitempty" yaml:"description,omitempty"` // optional
	Inputs      *JsonSchema `json:"inputs" yaml:"inputs"`
	Outputs     *JsonSchema `json:"outputs,omitempty" yaml:"outputs,omitempty"` // optional
}

type JsonSchema struct {
	Schema     string                 `json:"$schema,omitempty" yaml:"$schema,omitempty"`       // optional
	ID         string                 `json:"$id,omitempty" yaml:"$id,omitempty"`               // optional
	Properties map[string]*JsonSchema `json:"properties,omitempty" yaml:"properties,omitempty"` // optional
	Type       string                 `json:"type,omitempty" yaml:"type,omitempty"`             // optional
	Required   []string               `json:"required,omitempty" yaml:"required,omitempty"`     // optional
}

func NewScriptSchema(name string, description string, inputs *JsonSchema, outputs *JsonSchema) *ScriptSchema {
	return &ScriptSchema{Name: name, Description: description, Inputs: inputs, Outputs: outputs}
}

func NewScriptSchemaFromJSON(jsonSchema string) (*ScriptSchema, error) {
	var schema ScriptSchema
	if err := json.Unmarshal([]byte(jsonSchema), &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
