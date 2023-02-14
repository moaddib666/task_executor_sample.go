package tasks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"sync"
)

type ScriptRunner struct {
	wg *sync.WaitGroup
}

func NewScriptRunner() *ScriptRunner {
	return &ScriptRunner{
		wg: &sync.WaitGroup{},
	}
}

func (s *ScriptRunner) Execute(task *Task, execArgs interface{}) *TaskResult {
	var err error
	result := &TaskResult{
		Caller: task,
		Status: -1,
		Reason: "No task result from " + task.Location,
	}

	header := &Header{
		Meta: struct {
			Protocol string `json:"protocol"`
			Caller   string `json:"caller"`
			TaskName string `json:"taskName"`
		}{
			"v1",
			"self",
			task.Name,
		},
		Data: execArgs,
	}
	raw, err := json.Marshal(header)
	reader := bytes.NewReader(raw)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(task.Location)
	//
	pr, pw, err := os.Pipe()
	cmd.Stderr = pw
	cmd.Stdout = pw
	scanner := bufio.NewScanner(pr)
	cmd.Stdin = reader

	if err != nil {
		result.SetFail("Can't create pipe" + err.Error())
	}

	if err = cmd.Start(); err != nil {
		result.SetFail("can't start program" + err.Error())
	}
	_ = pw.Close()
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), result)
		if err == nil {
			break
		}
		log.Printf("TASK::%s %s", task.Name, scanner.Text())
	}

	if err = cmd.Wait(); err != nil {
		result.SetFail("Program ends with " + err.Error())
	}
	_ = pr.Close()
	log.Printf("TASK::%s Finised with status: %d, reason: %s", result.Caller.Name, result.Status, result.Reason)
	return result
}

func (s *ScriptRunner) ExecuteAsync(task *Task, execArgs interface{}) {
	s.wg.Add(1)
	go func() {
		s.Execute(task, execArgs)
		s.wg.Done()
	}()
}

func (s *ScriptRunner) WaitUntilComplete() {
	s.wg.Wait()
}
