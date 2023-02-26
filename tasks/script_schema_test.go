package tasks

import (
	"testing"
)

func TestNewScriptSchema(t *testing.T) {
	inputSchema := &JsonSchema{
		Type:       "object",
		Properties: map[string]*JsonSchema{},
	}
	outputSchema := &JsonSchema{
		Type:       "object",
		Properties: map[string]*JsonSchema{},
	}
	inputSchema.Properties["test"] = &JsonSchema{Type: "string"}
	schema := NewScriptSchema("test", "test", inputSchema, outputSchema)
	if schema.Name != "test" {
		t.Error("schema name not correct")
	}
	if schema.Description != "test" {
		t.Error("schema description not correct")
	}
	if schema.Inputs.Type != "object" {
		t.Error("schema inputs not correct")
	}
	if schema.Outputs.Type != "object" {
		t.Error("schema outputs not correct")
	}
}

func TestNewScriptSchemaFromJSON(t *testing.T) {
	jsonInputsSchema := `{"type": "object", "properties": {"test": {"type": "string"}}}`
	jsonOutputsSchema := `{"type": "object","properties": {"status": {"type": "integer"},"reason": {"type": "string"},"payload": {}},"required": ["status", "reason", "payload"]}`
	jsonSchema := `{
		"name": "test",
		"description": "test",
		"inputs": ` + jsonInputsSchema + `,
		"outputs": ` + jsonOutputsSchema + `	
	}`
	schema, err := NewScriptSchemaFromJSON(jsonSchema)
	if err != nil {
		t.Error(err)
	}
	if schema.Name != "test" {
		t.Error("schema name not correct")
	}
	if schema.Description != "test" {
		t.Error("schema description not correct")
	}
	if schema.Inputs.Type != "object" {
		t.Error("schema inputs not correct")
	}
	if schema.Outputs.Type != "object" {
		t.Error("schema outputs not correct")
	}
}
