package tasks

import (
	"context"
	"log"
	"runtime"
)

type ScriptRunnerPool struct {
	Runner
	poolSize  int
	workQueue chan struct {
		task     *Task
		execArgs interface{}
	}
	ctx context.Context
}

func (p *ScriptRunnerPool) startWorkerPool() {
	for i := 0; i < p.poolSize; i++ {
		go p.runPoolWorker(i)
	}
	// wait for context to be cancelled
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				log.Printf("Worker pool stopped")
				// close channel
				close(p.workQueue)
				return
			}
		}
	}()
}

func (p *ScriptRunnerPool) runPoolWorker(workerId int) {
	log.Printf("Worker %d started", workerId)
	for {
		request := <-p.workQueue
		// check if channel is closed
		if request.task == nil {
			log.Printf("Worker %d stopped", workerId)
			return
		}
		log.Printf("Worker %d executing task %s", workerId, request.task.Name)
		result := p.Execute(request.task, request.execArgs)
		log.Printf("Worker %d finished task %s with status %d details %s", workerId, request.task.Name, result.Status, result.Reason)
	}

}

func (p *ScriptRunnerPool) ExecuteAsync(task *Task, execArgs interface{}) {
	p.workQueue <- struct {
		task     *Task
		execArgs interface{}
	}{task: task, execArgs: execArgs}
}

func NewScriptRunnerPool(poolSize int, ctx context.Context) Runner {
	if poolSize <= 0 {
		poolSize = runtime.NumCPU()
	}
	runner := &ScriptRunnerPool{
		Runner:   NewScriptRunner(),
		poolSize: poolSize,
		workQueue: make(chan struct {
			task     *Task
			execArgs interface{}
		}),
		ctx: ctx,
	}
	runner.startWorkerPool()
	return runner
}
