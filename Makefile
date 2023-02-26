build:
	go build -o bin/executor_sample main.go
	go build -o bin/manage_tasks cmd/manage_tasks/main.go

clean:
	rm -f task_executor_sample

run:
	go run main.go

test:
	go test -v ./...
