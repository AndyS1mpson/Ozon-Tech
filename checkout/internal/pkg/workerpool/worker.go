// Worker
package workerpool

import (
	"log"
)

// Describe worker
type Worker[I, O any] struct {
	ID       int
	taskChan chan *Task[I, O]
}

// Create new worker instance
func NewWorker[I, O any](channel chan *Task[I, O], ID int) *Worker[I, O] {
	return &Worker[I, O]{
		ID:       ID,
		taskChan: channel,
	}
}

// Run worker to listen and proccess tasks
func (wr *Worker[I, O]) Start() {
	log.Printf("Starting worker %d\n", wr.ID)

	go func() {
		for task := range wr.taskChan {
			process(wr.ID, task)
		}
	}()
}
