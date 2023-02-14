package tasks

type Task struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func NewTask(name string, location string) *Task {
	return &Task{Name: name, Location: location}
}
