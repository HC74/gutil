package gutil

import (
	"fmt"
	"sync"
)

type Job interface {
	Do()
}

type Worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				job.Do()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Pool struct {
	workers     []Worker
	jobQueue    chan Job
	workerQueue chan chan Job
	wg          sync.WaitGroup
}

func NewPool(numWorkers, jobQueueLen int) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerQueue := make(chan chan Job, numWorkers)

	pool := &Pool{
		workers:     make([]Worker, numWorkers),
		jobQueue:    jobQueue,
		workerQueue: workerQueue,
	}

	for i := 0; i < numWorkers; i++ {
		pool.workers[i] = NewWorker(workerQueue)
	}

	return pool
}

func (p *Pool) Start() {
	for _, worker := range p.workers {
		worker.Start()
	}

	go p.dispatch()
}

func (p *Pool) Stop() {
	for _, worker := range p.workers {
		worker.Stop()
	}

	p.wg.Wait()
}

func (p *Pool) dispatch() {
	for {
		job := <-p.jobQueue

		p.wg.Add(1)
		go func(job Job) {
			worker := <-p.workerQueue
			worker <- job
			p.wg.Done()
		}(job)
	}
}

func (p *Pool) Submit(job Job) {
	p.jobQueue <- job
}

type ExampleJob struct {
	ID int
}

func (j ExampleJob) Do() {
	fmt.Printf("ExampleJob %d is running\n", j.ID)
}
