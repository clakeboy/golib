package components

import "sync"

//协程池
type GoroutinePool struct {
	Queue          chan interface{} //队列池
	Number         int              //并发协程数
	Total          int              //处理数据量
	Worker         func(obj ...interface{}) bool
	finishCallback func()
	wait           sync.WaitGroup
	stop           bool //关闭协程池信号
}

//新建一个协程池
func NewPoll(number int, worker func(obj ...interface{}) bool) *GoroutinePool {
	p := &GoroutinePool{
		Number: number,
		Worker: worker,
		wait:   sync.WaitGroup{},
	}
	return p
}

func (g *GoroutinePool) Start() {
	g.stop = false
	for i := 0; i < g.Number; i++ {
		g.wait.Add(1)
		go func(idx int) {
			isDone := false
			for !isDone {
				select {
				case task, ok := <-g.Queue:
					if !ok {
						isDone = true
					}
					g.Worker(task, idx, g)
				default:
					isDone = true
				}
				if g.stop {
					break
				}
			}
			g.wait.Done()
		}(i)
	}

	g.wait.Wait()

	if g.finishCallback != nil {
		g.finishCallback()
	}
	g.Stop()
}

func (g *GoroutinePool) AddTaskStrings(tasks []string) {
	total := len(tasks)
	g.Total = total
	g.Queue = make(chan interface{}, total)
	for _, obj := range tasks {
		g.Queue <- obj
	}
}

func (g *GoroutinePool) AddTaskInterface(tasks []interface{}) {
	total := len(tasks)
	g.Total = total
	g.Queue = make(chan interface{}, total)
	for _, obj := range tasks {
		g.Queue <- obj
	}
}

func (g *GoroutinePool) AddTask(task interface{}) {
	g.Queue <- task
}

func (g *GoroutinePool) Stop() {
	g.stop = true
	close(g.Queue)
}

func (g *GoroutinePool) SetFinishCallback(callback func()) {
	g.finishCallback = callback
}
