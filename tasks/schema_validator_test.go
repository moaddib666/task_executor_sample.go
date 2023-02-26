package tasks

import "testing"

func TestSchemaValidator(t *testing.T) {

	schema := &JsonSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
		Properties: map[string]*JsonSchema{
			"test": {Type: "string"},
			"id":   {Type: "integer"},
		},
		Required: []string{"test"},
	}

	validJsonData := map[string]interface{}{"test": "test"}

	err := validateJSONAgainstSchema(validJsonData, schema)

	if err != nil {
		t.Error(err)
	}

	invalidJsonData := map[string]interface{}{"test1": "test"}

	err = validateJSONAgainstSchema(invalidJsonData, schema)

	if err == nil {
		t.Error("json data should be invalid")
	}
}

func TestSchemaValidator_InvalidSchema(t *testing.T) {

	schema := &JsonSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
		Properties: map[string]*JsonSchema{
			"test": {Type: "string"},
			"id":   {Type: "integer"},
		},
		Required: []string{"test"},
	}

	validJsonData := map[string]interface{}{"test": "test", "id": "test"}

	err := validateJSONAgainstSchema(validJsonData, schema)

	if err == nil {
		t.Error("json data should be invalid")
	}
}

func TestSchemaValidator_InvalidSchema2(t *testing.T) {

	schema := &JsonSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
		Properties: map[string]*JsonSchema{
			"test": {Type: "string"},
			"id":   {Type: "integer"},
		},
		Required: []string{"test"},
	}

	validJsonData := map[string]interface{}{"test": "test", "id": 1}

	err := validateJSONAgainstSchema(validJsonData, schema)

	if err != nil {
		t.Error(err)
	}
}

func TestSchemaValidator_InvalidSchema3(t *testing.T) {

	schema := &JsonSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
		Properties: map[string]*JsonSchema{
			"test": {Type: "string"},
			"id":   {Type: "integer"},
		},
		Required: []string{"test"},
	}

	validJsonData := map[string]interface{}{"id": 1}

	err := validateJSONAgainstSchema(validJsonData, schema)

	if err == nil {
		t.Error("json data should be invalid")
	}
}
