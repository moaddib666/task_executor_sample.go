package tasks

import (
	"context"
	"runtime"
	"testing"
)

func TestNewScriptRunnerPool(t *testing.T) {
	// setup
	ctx, cancel := context.WithCancel(context.Background())
	// run
	pool := NewScriptRunnerPool(0, ctx).(*ScriptRunnerPool)

	// assert
	if pool.poolSize != runtime.NumCPU() {
		t.Errorf("Expected pool size to be %d, got %d", runtime.NumCPU(), pool.poolSize)
	}
	if pool.ctx != ctx {
		t.Errorf("Expected context to be %v, got %v", ctx, pool.ctx)
	}
	if pool.workQueue == nil {
		t.Error("Expected work queue to be initialized")
	}
	// teardown
	cancel()
}

func TestScriptRunnerPool_ExecuteAsync(t *testing.T) {
	// setup
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewScriptRunnerPool(5, ctx).(*ScriptRunnerPool)

	task := NewTask("test", "id")

	// run
	for i := 0; i < 50; i++ {
		pool.ExecuteAsync(task, map[string]interface{}{"message": "hello world"})
	}
	// ensure workers count is correct
	if pool.poolSize != 5 {
		t.Errorf("Expected pool size to be %d, got %d", 2, pool.poolSize)
	}

	// workers are asynchronous, so we need to wait for them to finish
	cancel()
	if len(pool.workQueue) != 0 {
		t.Errorf("Expected work queue to be empty, got %d", len(pool.workQueue))
	}
}
