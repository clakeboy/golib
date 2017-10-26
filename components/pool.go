package components

//协程池
type GoroutinePool struct {
	Queue chan interface{} //队列池
	Number int //并发协程数
	Total int  //处理数据量
	Worker func(obj... interface{}) bool

	result chan bool
	finishCallback func()
}

//新建一个协程池
func NewPoll(number int,worker func(obj... interface{}) bool) *GoroutinePool{
	p := &GoroutinePool{
		Number:number,
		Worker:worker,
	}
	return p
}

func (this *GoroutinePool) Start() {

	for i := 0; i < this.Number; i++ {
		go func() {
			for {
				task, ok := <-this.Queue
				if !ok {
					break
				}

				flag := this.Worker(task,i)
				this.result <- flag
			}
		}()
	}

	for j:=0;j<this.Total;j++ {
		res, ok := <-this.result
		if !ok {
			break
		}

		if !res {

		}
	}

	if this.finishCallback != nil {
		this.finishCallback()
	}
}

func (this *GoroutinePool) AddTaskStrings(tasks []string) {
	total := len(tasks)
	this.Total = total
	this.Queue = make(chan interface{},total)
	this.result = make(chan bool,total)
	for _,obj := range tasks {
		this.Queue <- obj
	}
}

func (this *GoroutinePool) AddTaskInterface(tasks []interface{}) {
	total := len(tasks)
	this.Total = total
	this.Queue = make(chan interface{},total)
	this.result = make(chan bool,total)
	for _,obj := range tasks {
		this.Queue <- obj
	}
}

func (this *GoroutinePool) AddTask(task interface{}) {
	this.Queue <- task
}

func (this *GoroutinePool) Stop() {
	close(this.Queue)
	close(this.result)
}

func (this *GoroutinePool) SetFinishCallback(callback func()) {
	this.finishCallback = callback
}