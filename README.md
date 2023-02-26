# Task Executor Sample

## Overview

This is a sample implementation of a task executor library.

## Quick start

### Build

```bash make build```
As a result you will have a binary `bin/manage_tasks` and `bin/task_executor_sample`

### Run

```bash ./bin/manage_tasks```



## Command line tool

`manage_tasks` is a sample command line tool that uses the task executor library.

### Usage

```text
env variable `TASK_STORE` not set, using default: /tmp/tasks.yamlNAME:
   manage_tasks - A new cli application

USAGE:
   manage_tasks [global options] command [command options] [arguments...]

COMMANDS:
   register    register a new task
   unregister  unregister an existing task
   list        list all registered tasks
   execute     execute a registered task
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### Register a task
To register a task:
- you need to provide a name and a location.
- the task must be executable.
- the task must accept `schema` argument and return a json scehma.
- the task must accept stdin and return a json output that fit protocol version `v1`.
-
```bash
./bin/manage_tasks register -name "bashTestTask" -location ./examples/bash.sh 
```

### Run a task
To run a task:
- you need to provide a name 
- the task must be registered.
- the task must accept stdin and return a json output that fit protocol version `v1`.

example:
```bash
./bin/manage_tasks execute -name pythonTestTask -taskArgs '{"cmd":"hostname","requestId":"214215161263","user":"tester"}'
```

Result
```bash
env variable `TASK_STORE` not set, using default: /tmp/tasks.yaml
2023/02/26 18:30:19 TASK::pythonTestTask stdin parsed {
2023/02/26 18:30:19 TASK::pythonTestTask    "meta": {
2023/02/26 18:30:19 TASK::pythonTestTask       "protocol": "v1",
2023/02/26 18:30:19 TASK::pythonTestTask       "caller": "self",
2023/02/26 18:30:19 TASK::pythonTestTask       "taskName": "pythonTestTask"
2023/02/26 18:30:19 TASK::pythonTestTask    },
2023/02/26 18:30:19 TASK::pythonTestTask    "data": {
2023/02/26 18:30:19 TASK::pythonTestTask       "cmd": "hostname",
2023/02/26 18:30:19 TASK::pythonTestTask       "requestId": "214215161263",
2023/02/26 18:30:19 TASK::pythonTestTask       "user": "tester"
2023/02/26 18:30:19 TASK::pythonTestTask    }
2023/02/26 18:30:19 TASK::pythonTestTask }
2023/02/26 18:30:19 TASK::pythonTestTask hello
2023/02/26 18:30:19 TASK::pythonTestTask Finised with status: 1, reason: ok
Task pythonTestTask executed successfully with 
result: 
{
  "Caller": {
    "name": "pythonTestTask",
    "location": "/Users/maxnikitenko/GolandProjects/task_executor_sample.go/examples/python.py",
    "schema": {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "properties": {
        "cmd": {
          "type": "string"
        },
        "requestId": {
          "type": "string"
        },
        "user": {
          "type": "string"
        }
      },
      "type": "object",
      "required": [
        "user",
        "cmd",
        "requestId"
      ]
    }
  },
  "status": 1,
  "reason": "ok",
  "payload": {
    "foo": "bar"
  }
}
```
## Protocol v1

### Json Schema

Task executor protocol `v1` is a json schema that defines the input and output of a task.

#### Input:
Script must accept following json input on stdin:
```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["meta", "data"],
    "properties": {
        "meta": {
            "type": "object",
            "required": ["protocol", "caller", "taskName"],
            "properties": {
                "protocol": {
                    "type": "string"
                },
                "caller": {
                    "type": "string"
                },
                "taskName": {
                    "type": "string"
                }
            }
        },
        "data": {
            "type": "object",
            "required": [],
            "properties": {}
            }
        }
    }
}
```
`data` object is a free form object that can be used to pass data to the task.
#### Output:
Script must return following json output on stdout:
```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["status", "reason", "payload"],
    "properties": {
        "status": {
            "type": "integer"
        },
        "reason": {
            "type": "string"
        },
        "payload": {}
    }
}
```
`payload` is a free form object that can be used to pass data to the caller.

You could also log any kind of data to stderr/stdout but last message from your task must be json formatted result.

#### AutoSchema:

Script should be able to return a json schema that describes the input and output of the task when called with `schema` argument.

Example:
```json
{
   "name": "python",
   "description": "Python example",
   "inputs": {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "type": "object",
      "properties": {
         "user": {
            "type": "string"
         },
         "cmd": {
            "type": "string"
         },
         "requestId": {
            "type": "string"
         }
      },
      "required": [
         "user",
         "cmd",
         "requestId"
      ]
   },
   "outputs": {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "type": "object",
      "properties": {
         "status": {
            "type": "integer"
         },
         "reason": {
            "type": "string"
         },
         "payload": {}
      },
      "required": [
         "status",
         "reason",
         "payload"
      ]
   }
}
```
