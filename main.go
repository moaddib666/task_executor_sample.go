package main

import (
	"github.com/moaddib666/task_executor_sample.go/tasks"
)

func main() {
	runner := tasks.NewScriptRunner()
	schema :=
		`{
		  "$schema": "http://json-schema.org/draft-07/schema#",
		  "$id": "task.Default",
		  "type": "object",
		  "properties": {
			"user": {
			  "type": "string"
			},
			"cmd": {
			  "type": "string"
			},
			"requestId": {
			  "type": "string"
			}
		  },
		  "required": [
			"user",
			"cmd",
			"requestId"
		  ]
		}
	`

	bashTask := tasks.NewTask("testBash", "./examples/bash.sh")
	err := bashTask.SetSchemaFromBytes([]byte(schema))
	if err != nil {
		panic(err)
	}
	pythonTask := tasks.NewTask("testPython", "./examples/python.py")
	err = pythonTask.SetSchemaFromBytes([]byte(schema))
	if err != nil {
		panic(err)
	}
	argsMap := make(map[string]interface{})
	argsMap["user"] = "tester"
	argsMap["cmd"] = "hostname"
	argsMap["requestId"] = "214215161263"
	runner.ExecuteAsync(bashTask, argsMap)
	runner.ExecuteAsync(pythonTask, argsMap)
	runner.WaitUntilComplete()
}
