package workerpool

import (
	"context"
	"sync"
)

type Task func(context.Context, int) error

type TaskError struct {
	WorkerId int
	Err      error
}

func (e TaskError) Error() string {
	return e.Err.Error()
}

func NewTaskError(workerId int, err error) TaskError {
	return TaskError{
		WorkerId: workerId,
		Err:      err,
	}
}


func worker(ctx context.Context, id int, tasks <-chan Task, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-tasks:
			err := t(ctx, id)
			if err != nil {
				errors <- NewTaskError(id, err)
			}
		}
	}
}

type WorkerPool struct {
	workerCount int
	tasks       chan Task
	errors      chan<- error
	wg          *sync.WaitGroup
}

func NewWorkerPool(workerCount int, workerBuffer int) *WorkerPool {
	cSizes := workerCount * workerBuffer
	tasks := make(chan Task, cSizes)
	errors := make(chan error, cSizes)

	return &WorkerPool{
		workerCount: workerCount,
		tasks:       tasks,
		errors:      errors,
		wg:          &sync.WaitGroup{},
	}
}

func (p *WorkerPool) Run(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go worker(ctx, i, p.tasks, p.errors, p.wg)
	}
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

func (p *WorkerPool) Close() {
	close(p.tasks)
	close(p.errors)
	p.Wait()
}

func (p *WorkerPool) Submit(t Task) {
	p.tasks <- t
}
