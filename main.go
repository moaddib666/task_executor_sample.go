package main

import (
	"context"
	"github.com/moaddib666/task_executor_sample.go/tasks"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
	ctx := context.Background()
	runner := tasks.NewScriptRunnerPool(0, ctx) // 0 means system cpu count
	runner.SetLogger(logger)
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

	// Output: async results
	for {
		select {
		case result := <-runner.GetAsyncResults():
			logger.Infof("Task executed: %s, status: %d, details: `%s`\n", result.Caller.Name, result.Status, result.Reason)
		default:
			return
		}
	}

}
