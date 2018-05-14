package task

import (
	"sync"
	"time"
	"ck_go_lib/components"
	"strings"
	"strconv"
)

//时间类型
type TimeType int

const (
	Second TimeType = iota
	Minute
	Hour
	DayOfMonth
	Month
	DayOfWeek
)

//任务规则
type Rule struct {
	Raw    string   //原始值
	Type   TimeType //时间类型
	IsLoop bool     //是否循环
	Value  int      //循环值
}

//任务项
type Item struct {
	RuleList     []*Rule //规则列表
	ExecFunc     func(item *Item) bool  //任务执行方法
	CallbackFunc func(item *Item)  //任务回调方法
	LastExecDate time.Time //最后执行任务的时间
}

//任务管理
type Management struct {
	list     []*Item //任务列表
	listLock sync.Mutex
}

//创建管理任务工厂方法
func NewManagement() *Management {
	return &Management{}
}

//添加一个任务项
func (m *Management) Add(item *Item) {
	m.listLock.Lock()
	m.list = append(m.list, item)
	m.listLock.Unlock()
}

//使用选项添加一个任务项
func (m *Management) AddTask(
	second string,
	minute string,
	hour string,
	dayOfMonth string,
	month string,
	dayOfWeek string,
	exec func(item *Item) bool,
	callback func(item *Item)) {

	var rules []*Rule
	rules = append(rules, m.explainString2Type(second,Second))
	rules = append(rules, m.explainString2Type(minute,Minute))
	rules = append(rules, m.explainString2Type(hour,Hour))
	rules = append(rules, m.explainString2Type(dayOfMonth,DayOfMonth))
	rules = append(rules, m.explainString2Type(month,Month))
	rules = append(rules, m.explainString2Type(dayOfWeek,DayOfWeek))

	item := &Item{
		RuleList:rules,
		ExecFunc:exec,
		CallbackFunc:callback,
	}

	m.Add(item)
}
//使用选项字符串添加一个任务项
func (m *Management) AddTaskString(taskStr string,exec func(item *Item) bool, callback func(item *Item)) {
	typeList := strings.Split(taskStr," ")
	m.AddTask(
		typeList[0],
		typeList[1],
		typeList[2],
		typeList[3],
		typeList[4],
		typeList[5],
		exec,
		callback,
	)
}
//把原始时间值转为任务规则
func (m *Management) explainString2Type(str string,timeType TimeType) *Rule {
	if str == "*" {
		return nil
	}
	rule := &Rule{
		Raw:str,
		Type:timeType,
	}
	val := strings.Split(str,"/")
	var err error
	if len(val) > 1 {
		rule.IsLoop = true
		rule.Value,err = strconv.Atoi(val[1])
		if err != nil {
			return nil
		}
	} else {
		rule.Value,err = strconv.Atoi(val[0])
		if err != nil {
			return nil
		}
	}

	return rule
}

//开始
func (m *Management) Start() {
	go m.run()
}

//运行
func (m *Management) run() {
	for {
		go m.execute(time.Now())
		time.Sleep(time.Second)
	}
}

//执行任务
func (m *Management) execute(currentDate time.Time) {
	poll := components.NewPoll(8, func(obj ...interface{}) bool {
		m.executeTask(obj[0].(*Item),currentDate)
		return true
	})

	var taskList []interface{}
	for _,v := range m.list {
		taskList = append(taskList,v)
	}

	poll.AddTaskInterface(taskList)
	poll.Start()
}

//执行任务项
func (m *Management) executeTask(item *Item,currentDate time.Time) {
	item.LastExecDate = currentDate
	if item.ExecFunc == nil {
		return
	}
	//解释时间是否已可执行
	for _,rule := range item.RuleList {
		if !m.explainType(rule,currentDate,item.LastExecDate) {
			return
		}
	}
	ok := item.ExecFunc(item)
	if ok && item.CallbackFunc != nil {
		item.CallbackFunc(item)
	}
}

//解释类型
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
			timeValue = int(currentDate.Sub(lastDate).Hours()/24)
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
