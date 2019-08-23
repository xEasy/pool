package workerpool

// Dispatcher 调度器或调派员，负责将工作分配给空闲的工人
type Dispatcher struct {
	pool  chan *Worker // 工人channel
	queue chan Job     // 工作队列channel
	stop  chan bool
}

// dispatch 获取工作将把它分配给第一个可用的工人
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.queue:
			worker := <-d.pool
			worker.channel <- job

		case stop := <-d.stop:
			if stop {
				for i := 0; i < cap(d.pool); i++ {
					worker := <-d.pool
					worker.stop <- true
					<-worker.stop
				}
				d.stop <- true
				return
			}
		}
	}
}

// newDispatcher is construction function for Dispatcher object.
func newDispatcher(p chan *Worker, q chan Job) *Dispatcher {
	d := &Dispatcher{
		pool:  p,
		queue: q,
		stop:  make(chan bool),
	}

	for i := 0; i < cap(d.pool); i++ {
		worker := newWorker(d.pool)
		worker.start()
	}

	go d.dispatch()
	return d
}
