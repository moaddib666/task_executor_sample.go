package tasks

type TaskResult struct {
	Caller  *Task
	Status  int         `json:"status"`
	Reason  string      `json:"reason"`
	Payload interface{} `json:"payload"`
}

func (tr *TaskResult) SetFail(reason string) {
	tr.Reason = reason
	tr.Status = -1
}
