package task

import (
	"github.com/clakeboy/golib/components"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TimeType 时间类型
type TimeType int

const (
	Second TimeType = iota
	Minute
	Hour
	DayOfMonth
	Month
	DayOfWeek
)

// Rule 任务规则
type Rule struct {
	Raw    string   //原始值
	Type   TimeType //时间类型
	IsLoop bool     //是否循环
	Value  int      //循环值
}

// Item 任务项
type Item struct {
	RuleList     []*Rule               //规则列表
	ExecFunc     func(item *Item) bool //任务执行方法
	CallbackFunc func(item *Item)      //任务回调方法
	LastExecDate time.Time             //最后执行任务的时间
	Args         []interface{}         //任务数据
	Lock         sync.RWMutex
}

// Management 任务管理
type Management struct {
	stop     bool    //是否停止
	list     []*Item //任务列表
	listLock sync.Mutex
}

// NewManagement 创建管理任务工厂方法
func NewManagement() *Management {
	return &Management{
		stop: true,
		//listLock: sync.Mutex{},
	}
}

// Add 添加一个任务项
func (m *Management) Add(item *Item) {
	m.listLock.Lock()
	m.list = append(m.list, item)
	m.listLock.Unlock()
}

// RemoveTask 删除列表中一个任务
func (m *Management) RemoveTask(item *Item) {
	m.listLock.Lock()
	for i, v := range m.list {
		if v == item {
			m.list = append(m.list[:i], m.list[i+1:]...)
			break
		}
	}
	m.listLock.Unlock()
}

// RemoveForeach foreach 回调方法删除一个任务
func (m *Management) RemoveForeach(f func(*Item) bool) {
	m.listLock.Lock()
	for i, v := range m.list {
		ok := f(v)
		if ok {
			m.list = append(m.list[:i], m.list[i+1:]...)
			break
		}
	}
	m.listLock.Unlock()
}

// ClearTask 清空任务列表
func (m *Management) ClearTask() {
	m.listLock.Lock()
	m.list = nil
	m.listLock.Unlock()
}

// AddTask 使用选项添加一个任务项
func (m *Management) AddTask(
	second string,
	minute string,
	hour string,
	dayOfMonth string,
	month string,
	dayOfWeek string,
	exec func(item *Item) bool,
	callback func(item *Item),
	args ...interface{}) {

	var rules []*Rule
	rules = append(rules, m.explainString2Type(second, Second))
	rules = append(rules, m.explainString2Type(minute, Minute))
	rules = append(rules, m.explainString2Type(hour, Hour))
	rules = append(rules, m.explainString2Type(dayOfMonth, DayOfMonth))
	rules = append(rules, m.explainString2Type(month, Month))
	rules = append(rules, m.explainString2Type(dayOfWeek, DayOfWeek))

	item := &Item{
		RuleList:     rules,
		ExecFunc:     exec,
		CallbackFunc: callback,
		Args:         args,
		//Lock:         sync.RWMutex{}Mutex{},
		LastExecDate: time.Now(),
	}

	m.Add(item)
}

// AddTaskString 使用选项字符串添加一个任务项
func (m *Management) AddTaskString(taskStr string, exec func(item *Item) bool, callback func(item *Item), args ...interface{}) {
	typeList := strings.Split(taskStr, " ")
	if len(typeList) < 6 {
		return
	}
	m.AddTask(
		typeList[0],
		typeList[1],
		typeList[2],
		typeList[3],
		typeList[4],
		typeList[5],
		exec,
		callback,
		args...,
	)
}

// 把原始时间值转为任务规则
func (m *Management) explainString2Type(str string, timeType TimeType) *Rule {
	if str == "*" {
		return nil
	}
	rule := &Rule{
		Raw:  str,
		Type: timeType,
	}
	val := strings.Split(str, "/")
	var err error
	if len(val) > 1 {
		rule.IsLoop = true
		rule.Value, err = strconv.Atoi(val[1])
		if err != nil {
			return nil
		}
	} else {
		rule.Value, err = strconv.Atoi(val[0])
		if err != nil {
			return nil
		}
	}

	return rule
}

// Start 开始
func (m *Management) Start() {
	m.stop = false
	go m.run()
}

func (m *Management) Stop() {
	m.stop = true
}

func (m *Management) Status() string {
	if m.stop {
		return "Stop"
	} else {
		return "Running"
	}
}

// 运行
func (m *Management) run() {
	for {
		//fmt.Println(time.Now().Format("15:04:05"),time.Now().Nanosecond() % 1000000,time.Now().Nanosecond()/1000000)
		if m.stop {
			break
		}
		go m.execute(time.Now())
		//fmt.Println(time.Second - time.Duration(time.Now().Nanosecond()))
		time.Sleep(time.Second - time.Duration(time.Now().Nanosecond()))
	}
}

// 执行任务
func (m *Management) execute(currentDate time.Time) {
	poll := components.NewPoll(runtime.NumCPU(), func(obj ...interface{}) bool {
		item := obj[0].(*Item)
		m.executeTask(item, currentDate)
		return true
	})

	var taskList []interface{}
	for _, v := range m.list {
		taskList = append(taskList, v)
	}
	poll.AddTaskInterface(taskList)
	poll.Start()
}

// 执行任务项
func (m *Management) executeTask(item *Item, currentDate time.Time) {
	if item.ExecFunc == nil {
		return
	}
	//解释时间是否已可执行
	item.Lock.RLock()
	for _, rule := range item.RuleList {
		if !m.explainType(rule, currentDate, item.LastExecDate) {
			item.Lock.RUnlock()
			return
		}
	}
	item.Lock.RUnlock()
	ok := item.ExecFunc(item)
	item.Lock.Lock()
	item.LastExecDate = currentDate
	item.Lock.Unlock()
	if ok && item.CallbackFunc != nil {
		item.CallbackFunc(item)
	}
}

// 解释类型
func (m *Management) explainType(rule *Rule, currentDate time.Time, lastDate time.Time) bool {
	if rule == nil {
		return true
	}
	if rule.IsLoop && lastDate.IsZero() {
		return true
	}
	var (
		timeValue int
	)
	switch rule.Type {
	case Second:
		if rule.IsLoop {
			timeValue = int(currentDate.Sub(lastDate).Seconds())
		} else {
			timeValue = currentDate.Second()
		}
	case Minute:
		if rule.IsLoop {
			timeValue = int(currentDate.Sub(lastDate).Minutes())
		} else {
			timeValue = currentDate.Minute()
		}
	case Hour:
		if rule.IsLoop {
			timeValue = int(currentDate.Sub(lastDate).Hours())
		} else {
			timeValue = currentDate.Hour()
		}
	case DayOfMonth:
		if rule.IsLoop {
			timeValue = int(currentDate.Sub(lastDate).Hours() / 24)
		} else {
			timeValue = currentDate.Day()
		}
	case Month:
		if rule.IsLoop {
			timeValue = int(currentDate.Month() - lastDate.Month())
		} else {
			timeValue = int(currentDate.Month())
		}
	case DayOfWeek:
		if rule.IsLoop {
			timeValue = int(currentDate.Weekday() - lastDate.Weekday())
		} else {
			timeValue = int(currentDate.Weekday())
		}
	default:
		return false
	}
	if rule.IsLoop {
		return timeValue >= rule.Value
	}
	return timeValue == rule.Value
}
