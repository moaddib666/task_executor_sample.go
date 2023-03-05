package tasks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type ScriptRunner struct {
	wg           *sync.WaitGroup
	asyncResults chan *TaskResult
	logger       *logrus.Logger
}

func NewScriptRunner() Runner {
	return &ScriptRunner{
		wg:           &sync.WaitGroup{},
		asyncResults: make(chan *TaskResult),
		logger:       logrus.New(),
	}
}

func (s *ScriptRunner) SetLogger(logger *logrus.Logger) {
	s.logger = logger
}

func (s *ScriptRunner) Execute(task *Task, execArgs interface{}) *TaskResult {
	var err error
	result := &TaskResult{
		Caller: task,
		Status: -1,
		Reason: "No task result from " + task.Location,
	}
	args, err := task.ParseArgs(execArgs)
	if err != nil {
		result.SetFail(err.Error())
		return result
	}
	caller, _ := filepath.Abs(os.Args[0])
	header := &Header{
		Meta: struct {
			Protocol string `json:"protocol"`
			Caller   string `json:"caller"`
			TaskName string `json:"taskName"`
		}{
			"v1",
			caller,
			task.Name,
		},
		Data: args,
	}
	raw, err := json.Marshal(header)
	reader := bytes.NewReader(raw)
	if err != nil {
		s.logger.Fatal(err)
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
		s.logger.Infof("TASK::%s %s", task.Name, scanner.Text())
	}

	if err = cmd.Wait(); err != nil {
		result.SetFail("Program ends with " + err.Error())
	}
	_ = pr.Close()
	return result
}

func (s *ScriptRunner) ExecuteAsync(task *Task, execArgs interface{}) {
	s.wg.Add(1)
	go func() {
		r := s.Execute(task, execArgs)
		s.logger.Infof("TASK::%s Finised with status: %d, reason: %s", r.Caller.Name, r.Status, r.Reason)
		s.wg.Done()
		s.saveTaskResult(r)
	}()
}

func (s *ScriptRunner) saveTaskResult(result *TaskResult) {
	s.asyncResults <- result
}

func (s *ScriptRunner) WaitUntilComplete() {
	s.wg.Wait()
}

func (s *ScriptRunner) GetAsyncResults() <-chan *TaskResult {
	return s.asyncResults
}

func (s *ScriptRunner) GetAsyncResultsCount() int {
	return len(s.asyncResults)
}
