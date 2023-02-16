package gutil

import (
	"fmt"
	"sync"
)

// 任务结构体
type task struct {
	id int // 任务ID
}

// 线程池结构体
type threadPool struct {
	capacity int            // 线程池容量
	tasks    chan *task     // 任务通道
	wg       sync.WaitGroup // 等待组
}

// 线程函数，用于并发执行任务
func (p *threadPool) worker(id int) {
	defer p.wg.Done()
	for t := range p.tasks {
		fmt.Printf("worker %d processing task %d\n", id, t.id)
	}
}

// 初始化线程池
func newThreadPool(capacity int) *threadPool {
	p := &threadPool{
		capacity: capacity,
		tasks:    make(chan *task),
	}
	p.wg.Add(capacity)
	for i := 0; i < capacity; i++ {
		go p.worker(i)
	}
	return p
}

// 提交任务
func (p *threadPool) submit(t *task) {
	p.tasks <- t
}

// 等待所有任务执行完毕
func (p *threadPool) wait() {
	close(p.tasks)
	p.wg.Wait()
}
