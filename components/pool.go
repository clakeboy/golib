package components

import "sync"

//协程池
type GoroutinePool struct {
	Queue chan interface{} //队列池
	Number int //并发协程数
	Total int  //处理数据量
	Worker func(obj... interface{}) bool
	finishCallback func()
	wait *sync.WaitGroup
}

//新建一个协程池
func NewPoll(number int,worker func(obj... interface{}) bool) *GoroutinePool{
	p := &GoroutinePool{
		Number:number,
		Worker:worker,
		wait:&sync.WaitGroup{},
	}
	return p
}

func (this *GoroutinePool) Start() {
	for i := 0; i < this.Number; i++ {
		go func(idx int) {
			isDone := true
			for isDone {
				select {
				case task := <-this.Queue:
					this.Worker(task,idx)
				default:
					isDone = false
				}
			}
			this.wait.Done()
		}(i)
		this.wait.Add(1)
	}

	this.wait.Wait()

	if this.finishCallback != nil {
		this.finishCallback()
	}
}

func (this *GoroutinePool) AddTaskStrings(tasks []string) {
	total := len(tasks)
	this.Total = total
	this.Queue = make(chan interface{},total)
	for _,obj := range tasks {
		this.Queue <- obj
	}
}

func (this *GoroutinePool) AddTaskInterface(tasks []interface{}) {
	total := len(tasks)
	this.Total = total
	this.Queue = make(chan interface{},total)
	for _,obj := range tasks {
		this.Queue <- obj
	}
}

func (this *GoroutinePool) AddTask(task interface{}) {
	this.Queue <- task
}

func (this *GoroutinePool) Stop() {
	close(this.Queue)
}

func (this *GoroutinePool) SetFinishCallback(callback func()) {
	this.finishCallback = callback
}