package tasks

type Header struct {
	Meta struct {
		Protocol string `json:"protocol"`
		Caller   string `json:"caller"`
		TaskName string `json:"taskName"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}
