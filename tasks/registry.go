package tasks

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TaskStore struct {
	filename string
	tasks    map[string]*Task
}

func NewTaskStore(filename string) (*TaskStore, error) {
	store := &TaskStore{
		filename: filename,
		tasks:    make(map[string]*Task),
	}

	err := store.load()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *TaskStore) AddTask(task *Task) error {
	if _, ok := s.tasks[task.Name]; ok {
		return fmt.Errorf("task with name '%s' already exists", task.Name)
	}

	s.tasks[task.Name] = task
	return s.save()
}

func (s *TaskStore) GetTask(name string) (*Task, error) {
	if task, ok := s.tasks[name]; ok {
		return task, nil
	}

	return nil, fmt.Errorf("task with name '%s' not found", name)
}

func (s *TaskStore) GetTasks() []*Task {
	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *TaskStore) DeleteTask(name string) error {
	if _, ok := s.tasks[name]; !ok {
		return fmt.Errorf("task with name '%s' not found", name)
	}

	delete(s.tasks, name)
	return s.save()
}

func (s *TaskStore) Clear() error {
	s.clear()
	return s.save()
}

func (s *TaskStore) load() error {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking if task file exists %s", err.Error())
	}

	bytes, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return fmt.Errorf("error reading task file %s", err.Error())
	}

	err = yaml.Unmarshal(bytes, &s.tasks)
	if err != nil {
		return fmt.Errorf("error unmarshaling task data %s", err.Error())
	}

	return nil
}

func (s *TaskStore) save() error {
	bytes, err := yaml.Marshal(s.tasks)
	if err != nil {
		return fmt.Errorf("error marshaling task data %s", err.Error())
	}

	dir := filepath.Dir(s.filename)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("error creating task directory %s", err.Error())
	}

	err = ioutil.WriteFile(s.filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing task file %s", err.Error())
	}
	return nil
}

func (s *TaskStore) clear() {
	s.tasks = make(map[string]*Task)
}
