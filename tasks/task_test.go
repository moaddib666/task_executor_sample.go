package tasks

import "testing"

func TestNewTask(t *testing.T) {
	task := NewTask("test", "test")
	task.Accept = nil
	if task.Name != "test" {
		t.Error("task name not correct")
	}
	if task.Location != "test" {
		t.Error("task location not correct")
	}
	if task.Accept != nil {
		t.Error("task schema should be nil")
	}
}

func TestTask_ParseArgs(t *testing.T) {
	schema := &JsonSchema{Schema: "http://json-schema.org/draft-04/schema#", Type: "object", Properties: map[string]*JsonSchema{}}
	schema.Properties["test"] = &JsonSchema{Type: "string"}
	task := NewTask("test", "test")
	task.Accept = schema
	args := map[string]interface{}{"test": "test"}
	_, err := task.ParseArgs(args)
	if err != nil {
		t.Error(err)
	}
}

func TestTask_ParseArgs_Invalid(t *testing.T) {
	schema := &JsonSchema{Schema: "http://json-schema.org/draft-04/schema#", Type: "object", Properties: map[string]*JsonSchema{}}
	schema.Properties["test"] = &JsonSchema{Type: "string"}
	task := NewTask("test", "test")
	task.Accept = schema
	args := map[string]interface{}{"test": 1}
	_, err := task.ParseArgs(args)
	if err == nil {
		t.Error("args should be invalid")
	}
}
