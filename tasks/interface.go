package tasks

type Runner interface {
	Execute(task *Task, execArgs interface{}) *TaskResult
	ExecuteAsync(task *Task, execArgs interface{})
	WaitUntilComplete()
}
