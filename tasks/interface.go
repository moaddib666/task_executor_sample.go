package tasks

type Runner interface {
	Execute(task *Task) *TaskResult
}
