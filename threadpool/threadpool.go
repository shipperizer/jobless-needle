package threadpool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

const (
	ContextDoneMsg = "ThreadPool context done"
)

type TaskResult[T any] struct {
	key   string
	value T
}

func NewTaskCompletion[T any](taskID string) *TaskResult[T] {
	return &TaskResult[T]{
		key: taskID,
	}
}

func NewTaskResult[T any](taskID string, result T) *TaskResult[T] {
	return &TaskResult[T]{
		key:   taskID,
		value: result,
	}
}

type ThreadPool struct {
	cancel    context.CancelFunc
	done      <-chan struct{}
	stopped   *atomic.Bool
	mainGroup *sync.WaitGroup
}

func (t *ThreadPool) Stop() {
	t.stopped.Swap(true)
	t.cancel()
}

func (t *ThreadPool) Await() {
	t.mainGroup.Wait()
}

func (t *ThreadPool) Done() <-chan struct{} {
	return t.done
}

func (t *ThreadPool) Submit(taskID string, command any, results chan<- *TaskResult[any], wg *sync.WaitGroup) error {
	if t.stopped.Load() {
		return fmt.Errorf("ThreadPool is stopped")
	}

	go t.worker(taskID, command, results, wg)()
	return nil
}

func (t *ThreadPool) worker(taskID string, command any, results chan<- *TaskResult[any], wg *sync.WaitGroup) func() {
	return func() {
		t.executeFunc(taskID, command, results, wg)
	}
}

func (t *ThreadPool) executeFunc(taskID string, processor any, results chan<- *TaskResult[any], wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	select {
	case <-t.done:
		log.Println(ContextDoneMsg)
	default:
		t.mainGroup.Add(1)
		defer t.mainGroup.Done()

		switch task := processor.(type) {
		case func():
			task()
			results <- NewTaskCompletion[any](taskID)
		case func() any:
			value := task()
			results <- NewTaskResult[any](taskID, value)
		}
	}
}

func NewThreadPool() *ThreadPool {
	ctx, cancelFunc := context.WithCancel(context.Background())
	stopped := atomic.Bool{}
	stopped.Store(false)
	return &ThreadPool{
		done:      ctx.Done(),
		cancel:    cancelFunc,
		stopped:   &stopped,
		mainGroup: &sync.WaitGroup{},
	}
}

func NewThreadPoolWithContext(ctx context.Context) *ThreadPool {
	ctx, cancelFunc := context.WithCancel(ctx)
	stopped := atomic.Bool{}
	stopped.Store(false)
	return &ThreadPool{
		done:      ctx.Done(),
		cancel:    cancelFunc,
		stopped:   &stopped,
		mainGroup: &sync.WaitGroup{},
	}
}

func Take[T any](results chan *TaskResult[any], n int) []TaskResult[T] {
	ret := make([]TaskResult[T], 0)

	for i := 0; i < n; i++ {
		r := <-results
		result := NewTaskResult[T](r.key, r.value.(T))
		ret = append(ret, *result)
	}

	return ret
}
