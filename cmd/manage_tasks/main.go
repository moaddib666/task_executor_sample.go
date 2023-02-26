package main

import (
	"encoding/json"
	"fmt"
	"github.com/moaddib666/task_executor_sample.go/tasks"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path/filepath"
)

func getStore() (*tasks.TaskStore, error) {
	filename := os.Getenv("TASK_STORE")
	if filename == "" {
		filename = "/tmp/tasks.yaml"
		fmt.Printf("env variable `TASK_STORE` not set, using default: %s\n\n", filename)
	}
	store, err := tasks.NewTaskStore(filename)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func main() {
	store, err := getStore()
	if err != nil {
		fmt.Printf("Error getting task store: %s", err.Error())
		return
	}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "register",
				Usage: "register a new task",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "task name",
					},
					&cli.StringFlag{
						Name:  "location",
						Usage: "task location",
					},
				},
				Action: func(c *cli.Context) error {
					name := c.String("name")
					location := c.String("location")
					if name == "" {
						return fmt.Errorf("name is required")
					}
					if location == "" {
						return fmt.Errorf("location is required")
					}
					// resolve to absolute path
					location, err = filepath.Abs(location)
					if err != nil {
						return err
					}
					cmd := exec.Command(location, "schema")
					out, err := cmd.Output()
					if err != nil {
						fmt.Printf("Error executing script with `schema` arg: %s", err.Error())
						return err
					}
					schema, err := tasks.NewScriptSchemaFromJSON(string(out))
					if err != nil {
						fmt.Printf("Error parsing schema: %s", err.Error())
						return err
					}
					task := tasks.NewTask(name, location)
					task.Accept = schema.Inputs
					err = store.AddTask(task)
					if err == nil {
						fmt.Printf("Task %s registered successfully", name)
					}
					return err
				},
			},
			{
				Name:  "unregister",
				Usage: "unregister an existing task",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "task name",
					},
				},
				Action: func(c *cli.Context) error {
					name := c.String("name")
					return store.DeleteTask(name)
				},
			},
			{
				Name:  "list",
				Usage: "list all registered tasks",
				Action: func(c *cli.Context) error {
					taskList := store.GetTasks()
					if len(taskList) == 0 {
						fmt.Println("No tasks registered")
						return nil
					}
					for _, task := range taskList {
						fmt.Printf("---\nname: %s\nlocation: %s\nschema: %+v\n", task.Name, task.Location, task.JsonSchema())
					}
					return nil
				},
			},
			{
				Name:  "execute",
				Usage: "execute a registered task",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "task name",
					},
					&cli.StringFlag{
						Name:  "taskArgs",
						Usage: "task arguments",
					},
				},
				Action: func(c *cli.Context) error {
					t, err := store.GetTask(c.String("name"))
					if err != nil {
						fmt.Printf("Error getting task: %s", err.Error())
						return err
					}
					runner := tasks.NewScriptRunner()
					args := c.String("taskArgs")
					if args == "" {
						return fmt.Errorf("can't run task %s with no arguments", t.Name)
					}
					var argsMap = map[string]interface{}{}
					err = json.Unmarshal([]byte(args), &argsMap)
					if err != nil {
						fmt.Printf("Error parsing task arguments: %s", err.Error())
						return err
					}

					result := runner.Execute(t, argsMap)
					jsonResult, err := json.MarshalIndent(result, "", "  ")
					fmt.Printf("Task: `%s`\n", t.Name)
					fmt.Printf("Result:\n%s\n", jsonResult)
					return err
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
