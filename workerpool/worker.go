package workerpool

type Worker struct {
	pool    chan *Worker
	channel chan Job
	stop    chan bool
}

// start 工人就绪
func (w *Worker) start() {
	go func() {
		var job Job
		for {
			w.pool <- w                  //循环，初次和任务完了注册自己到工人池里，以便等待分派任务
			select {
			case job = <-w.channel:
				job()
			case stop := <-w.stop:
				if stop {
					w.stop <- true
					return
				}
			}
		}
	}()
}

// 生成新的worker对象
func newWorker(p chan *Worker) *Worker {
	return &Worker{
		pool:    p,
		channel: make(chan Job),
		stop:    make(chan bool),
	}
}
