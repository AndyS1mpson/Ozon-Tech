// Implementation of a pool of workers
package workerpool

import (
	"context"
)

// Implement worker pool
type Pool[I, O any] struct {
	ctx   context.Context
	limit int
	tasks chan *Task[I, O]
}

// Create a new worker pool
func NewPool[I, O any](ctx context.Context, limit int) *Pool[I, O] {
	p := &Pool[I, O]{
		ctx:   ctx,
		limit: limit,
		tasks: make(chan *Task[I, O], limit),
	}

	go func() {
		<-ctx.Done()
		close(p.tasks)
	}()

	for i := 0; i < p.limit; i++ {
		worker := NewWorker(p.tasks, i)
		worker.Start()
	}

	return p
}

// Send a task to the worker pool tasks channel and execute it
func (p *Pool[I, O]) Exec(data I, step func(ctx context.Context, i I) (O, error)) Promise[O] {
	result := Promise[O]{Out: make(chan Either[O], 1)}
	select {
	case <-p.ctx.Done():
		result.Out <- Either[O]{Err: p.ctx.Err()}
	case p.tasks <- NewTask[I, O](p.ctx, data, step, result):
	}

	return result
}
