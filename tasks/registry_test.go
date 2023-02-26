package tasks

import (
	"os"
	"testing"
)

const testStoreLocation = "/tmp/task_executor/"
const testStoreName = "tasks.yaml"

// bootstrap the test store
func setup() {
	if _, err := os.Stat(testStoreLocation); os.IsNotExist(err) {
		err = os.Mkdir(testStoreLocation, 0755)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(testStoreLocation + testStoreName); !os.IsNotExist(err) {
		os.Remove(testStoreLocation + testStoreName)
	}
	store, err := NewTaskStore(testStoreLocation + testStoreName)
	if err != nil {
		panic(err)
	}
	store.clear()
	// empty json schema
	schema := &JsonSchema{
		Schema:     "http://json-schema.org/draft-04/schema#",
		Type:       "object",
		Properties: map[string]*JsonSchema{},
		Required:   []string{},
	}

	pythonTask := NewTask("python", "/usr/bin/python")
	pythonTask.Accept = schema
	bashTask := NewTask("bash", "/bin/bash")
	bashTask.Accept = schema
	err = store.AddTask(pythonTask)
	if err != nil {
		panic(err)
	}

	err = store.AddTask(bashTask)
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	os.Remove(testStoreLocation)
}

func TestMain(m *testing.M) {
	// perform setup
	setup()
	// run tests
	code := m.Run()

	// perform cleanup
	cleanup()

	// indicate whether tests passed or failed
	os.Exit(code)
}

func TestNewTaskStore(t *testing.T) {
	store, err := NewTaskStore(testStoreLocation + testStoreName)
	if err != nil {
		t.Error(err)
	}

	if store == nil {
		t.Error("store is nil")
	}

	if store.filename != testStoreLocation+testStoreName {
		t.Error("filename is not correct")
	}

	if len(store.tasks) != 2 {
		t.Error("tasks not loaded")
	}
}

func TestTaskStore_AddTask(t *testing.T) {
	store, err := NewTaskStore(testStoreLocation + testStoreName)
	if err != nil {
		t.Error(err)
	}
	store.clear() // clean up from previous tests
	pythonTask := NewTask("python", "/usr/bin/python")
	bashTask := NewTask("bash", "/bin/bash")
	err = store.AddTask(pythonTask)
	if err != nil {
		t.Error(err)
	}
	if len(store.tasks) != 1 {
		t.Error("task not added")
	}

	err = store.AddTask(bashTask)
	if err != nil {
		t.Error(err)
	}

	if len(store.tasks) != 2 {
		t.Error("task not added")
	}

	err = store.AddTask(pythonTask)
	if err == nil {
		t.Error("task added twice")
	}

	if len(store.tasks) != 2 {
		t.Error("task added twice")
	}

}

func TestTaskStore_GetTask(t *testing.T) {
	store, err := NewTaskStore(testStoreLocation + testStoreName)
	if err != nil {
		t.Error(err)
	}

	task, err := store.GetTask("python")
	if err != nil {
		t.Error(err)
	}

	if task.Name != "python" {
		t.Error("task name not correct")
	}

	if task.Location != "/usr/bin/python" {
		t.Error("task location not correct")
	}

	task, err = store.GetTask("bash")
	if err != nil {
		t.Error(err)
	}

	if task.Name != "bash" {
		t.Error("task name not correct")
	}

	if task.Location != "/bin/bash" {
		t.Error("task location not correct")
	}

	task, err = store.GetTask("does not exist")
	if err == nil {
		t.Error("task should not exist")
	}
}

func TestTaskStore_GetTasks(t *testing.T) {
	store, err := NewTaskStore(testStoreLocation + testStoreName)
	if err != nil {
		t.Error(err)
	}

	tasks := store.GetTasks()
	if len(tasks) != 2 {
		t.Error("tasks not returned")
	}
}

func TestTaskStore_DeleteTask(t *testing.T) {
	store, err := NewTaskStore(testStoreLocation + testStoreName)
	if err != nil {
		t.Error(err)
	}

	err = store.DeleteTask("python")
	if err != nil {
		t.Error(err)
	}

	if len(store.tasks) != 1 {
		t.Error("task not deleted")
	}

	err = store.DeleteTask("does not exist")
	if err == nil {
		t.Error("task should not exist")
	}

	if len(store.tasks) != 1 {
		t.Error("task should not be deleted")
	}

	err = store.DeleteTask("bash")
	if err != nil {
		t.Error(err)
	}
}
