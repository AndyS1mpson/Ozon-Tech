// Worker Task
package workerpool

import (
	"context"
	"log"
)

// Describe the result of a task or an error
type Either[T any] struct {
	Value *T
	Err   error
}

// Describe the channel through which the results of the work will be transmitted
type Promise[T any] struct {
	Out chan Either[T]
}

// Describe the task for worker,
// that contains a data, a context and a result promise
type Task[I, O any] struct {
	data   I
	ctx    context.Context
	step   func(ctx context.Context, i I) (O, error)
	result Promise[O]
}

// Create new task for worker
func NewTask[I, O any](ctx context.Context, data I, f func(ctx context.Context, i I) (O, error), result Promise[O]) *Task[I, O] {
	return &Task[I, O]{
		data:   data,
		ctx:    ctx,
		step:   f,
		result: result,
	}
}

// Run task and send the result to the task result channel
func process[I, O any](workerID int, task *Task[I, O]) {
	log.Printf("Worker %d processes task", workerID)
	v, err := task.step(task.ctx, task.data)
	task.result.Out <- Either[O]{
		Value: &v,
		Err:   err,
	}
}
