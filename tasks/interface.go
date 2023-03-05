package tasks

import (
	"github.com/sirupsen/logrus"
)

type Runner interface {
	Execute(task *Task, execArgs interface{}) *TaskResult
	ExecuteAsync(task *Task, execArgs interface{})
	WaitUntilComplete()
	GetAsyncResults() <-chan *TaskResult
	SetLogger(logger *logrus.Logger)
	saveTaskResult(result *TaskResult)
}
