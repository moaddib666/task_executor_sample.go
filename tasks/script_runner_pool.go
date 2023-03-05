package tasks

import (
	"context"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
)

type ScriptRunnerPool struct {
	*ScriptRunner
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
				p.ScriptRunner.logger.Infof("Worker pool stopped")
				// close channel
				close(p.workQueue)
				return
			}
		}
	}()
}

func (p *ScriptRunnerPool) runPoolWorker(workerId int) {
	p.ScriptRunner.logger.Debugf("Worker %d started", workerId)
	for {
		request := <-p.workQueue
		// check if channel is closed
		if request.task == nil {
			p.ScriptRunner.logger.Debugf("Worker %d stopped", workerId)
			return
		}
		p.ScriptRunner.logger.Debugf("Worker %d executing task %s", workerId, request.task.Name)
		p.ScriptRunner.ExecuteAsync(request.task, request.execArgs)
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
		ScriptRunner: &ScriptRunner{
			wg:           &sync.WaitGroup{},
			asyncResults: make(chan *TaskResult),
			logger:       logrus.New(),
		},
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
