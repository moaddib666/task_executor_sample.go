package main

import (
	"scriptExecutor/tasks"
)

func main() {
	runner := tasks.NewScriptRunner()
	bashTask := tasks.NewTask("testBash", "./examples/bash.sh")
	pythonTask := tasks.NewTask("testPython", "./examples/python.py")
	args := &struct {
		User      string `json:"user"`
		Cmd       string `json:"cmd"`
		RequestId string `json:"requestId"`
	}{
		"tester",
		"hostname",
		"214215161263",
	}
	runner.ExecuteAsync(bashTask, args)
	runner.ExecuteAsync(pythonTask, args)
	runner.WaitUntilComplete()
}
