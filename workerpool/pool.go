package workerpool

import "sync"

// Job 定义任务，任务为回调函数
type Job func(args ...interface{})

// Pool 池用来管理和分配任务
type Pool struct {
	Jobs       chan Job
	dispatcher *Dispatcher
	wg         sync.WaitGroup
}

// NewPool 创建指定工人数量和队列长度的池
func NewPool(workersCount int, queueSize int) *Pool {
	queueChannel := make(chan Job, queueSize)
	workersChannel := make(chan *Worker, workersCount)

	pool := &Pool{
		Jobs:       queueChannel,
		dispatcher: newDispatcher(workersChannel, queueChannel),
	}

	return pool
}

// JobDone 在任务完成时需要被调用
func (p *Pool) JobDone() {
	p.wg.Done()
}

// WaitCount 配置选项，设置我们当前有多少个工作需要完成
// we should wait when calling WaitAll.
func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

// 阻塞等待，直到完成所有工作
func (p *Pool) WaitAll() {
	p.wg.Wait()
}

// 释放池资源.
func (p *Pool) Release() {
	p.dispatcher.stop <- true
	<-p.dispatcher.stop
}

// 将回调函数工作放到池中,池中调度器将会分配worker执行工作
func (p *Pool) Enqueue(job Job, args ...interface{}) {
	arg := args
	p.Jobs <- func(args ...interface{}) {
		defer p.JobDone()
		job(arg...)
	}
}
