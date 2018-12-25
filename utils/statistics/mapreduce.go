package statistics

import (
	"fmt"
	"github.com/clakeboy/golib/utils"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type Stat interface {
	Map(data interface{}) interface{}
	Reduce(data interface{}) interface{}
}

//计算结果值
type ReduceValue struct {
	Count int     //计算个数
	Value float64 //计算合值
}

//MAP带时间值
type MapValue struct {
	Value    interface{}
	Datetime *StatDate
}

//统计时间
type StatDate struct {
	Year    int //年
	Month   int //月
	Week    int //年中第几周
	Day     int //月份中的第几天
	YearDay int //年份中的第几天
}

//新创建一个统计时间
func NewStatDate(unix_time interface{}) *StatDate {
	var date time.Time
	v := reflect.ValueOf(unix_time)
	switch v.Kind() {
	case reflect.String:
		i, err := strconv.ParseInt(v.String(), 10, 64)
		if err != nil {
			date = time.Now()
		}
		date = time.Unix(i, 0)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		date = time.Unix(v.Int(), 0)
	}

	year, week := date.ISOWeek()
	return &StatDate{
		Year:    year,
		Month:   int(date.Month()),
		Week:    week,
		Day:     date.Day(),
		YearDay: date.YearDay(),
	}
}

//reduce相加
func (r *ReduceValue) Add(val *ReduceValue) {
	r.Count += val.Count
	r.Value += val.Value
}

//计算结果集
type ReduceResult map[string]*ReduceValue

//计算条件结果集
type CondReduceResult map[string]ReduceResult

//计算条件结果集
type CondDateReduceResult map[string]DateReduceResult

//计算条件时间结果集
type DateReduceResult map[int]map[string]map[int]ReduceResult

//数据集合
type MapData map[string][]*MapValue

//条件数据集合
type CondMapData map[*MapCond]MapData

//条件字条 =,!,>,<
type CondColumn struct {
	Field string
	Icon  string
}

//条件类型
type MapCond struct {
	Name      string  //条件名称
	Condition utils.M //条件数据
	IsAll     bool    //是否全部数据统计
}

//创建一个统计条件类
func NewMapCond(name string, cond_data utils.M, is_all bool) *MapCond {
	return &MapCond{Name: name, Condition: cond_data, IsAll: is_all}
}

//验证数据是否满足条件
func (mc *MapCond) Valid(data utils.M) bool {
	if mc.IsAll {
		return true
	}
	and, ok := mc.Condition["AND"]
	if ok {
		return mc.conditionRecursion(and.(utils.M), data, "AND")
	}
	or, ok := mc.Condition["OR"]
	if ok {
		return mc.conditionRecursion(or.(utils.M), data, "OR")
	}
	return mc.conditionRecursion(mc.Condition, data, "AND")
}

//递归条件判断
func (mc *MapCond) conditionRecursion(cond utils.M, data utils.M, icon string) bool {
	var cond_icon *CondColumn
	if icon == "AND" {
		for k, v := range cond {
			cond_icon = mc.explainColumn(k)
			var ok bool
			if k == "AND" || k == "OR" {
				ok = mc.conditionRecursion(v.(utils.M), data, k)
			} else {
				ok = mc.Compared(v, data[cond_icon.Field], cond_icon.Icon)
			}
			if !ok {
				return false
			}
		}
		return true
	} else if icon == "OR" {
		for k, v := range cond {
			cond_icon = mc.explainColumn(k)
			var ok bool
			if k == "AND" || k == "OR" {
				ok = mc.conditionRecursion(v.(utils.M), data, k)
			} else {
				ok = mc.Compared(v, data[cond_icon.Field], cond_icon.Icon)
			}
			//fmt.Printf("%s %s %t %t %t %t\n",k,ok,utils.ConvertFloat(data[cond_icon.Field]),utils.ConvertFloat(v),data[cond_icon.Field],v)
			if ok {
				return true
			}
		}
		return false
	}
	return false
}

//以传入的条件对比两个值
func (mc *MapCond) Compared(v1, v2 interface{}, icon string) bool {
	switch icon {
	case ">":
		return utils.ConvertFloat(v1) > utils.ConvertFloat(v2)
	case "<":
		return utils.ConvertFloat(v1) < utils.ConvertFloat(v2)
	case "!":
		t := reflect.TypeOf(v1)
		if t.Kind() == reflect.String {
			return v1 != v2
		}
		return utils.ConvertFloat(v1) != utils.ConvertFloat(v2)
	default:
		t := reflect.TypeOf(v1)
		if t.Kind() == reflect.String {
			return v1 == v2
		}
		if t.String() == "*regexp.Regexp" {
			return v1.(*regexp.Regexp).MatchString(fmt.Sprintf("%v", v2))
		}
		return utils.ConvertFloat(v1) == utils.ConvertFloat(v2)
	}
}

//解释字段名
func (mc *MapCond) explainColumn(column string) *CondColumn {
	reg := regexp.MustCompile(`(.+?)\[(!|>|<)\]`)
	match := reg.FindStringSubmatch(column)
	field := &CondColumn{}
	if len(match) > 0 {
		field.Field = match[1]
		field.Icon = match[2]
	} else {
		field.Field = column
		field.Icon = "="
	}
	return field
}

//统计功能类
type MapReduce struct {
	maps       CondMapData //产生的MAP数据集
	keys       []string    //需要统计的KEY集合
	timeColumn string      //统计时间字段,用于生成统计日期段
}

//创建统计功能实体类
func NewMapReduce() *MapReduce {
	return &MapReduce{
		maps: CondMapData{},
	}
}

//设置记录时间字段
func (u *MapReduce) SetDateColumn(key string) {
	u.timeColumn = key
}

//开始运行统计
func (u *MapReduce) Run(keys []string, data []utils.M, conditions []*MapCond) {
	allCond := NewMapCond("ALL", nil, true)
	if conditions == nil {
		conditions = []*MapCond{allCond}
	} else {
		conditions = append(conditions, allCond)
	}
	u.keys = keys
	for _, v := range data {
		if u.keys == nil {
			u.keys = utils.MapKeys(v)
		}

		for _, k := range u.keys {
			u.Map(k, v, conditions)
		}

	}
}

//分类数据
func (u *MapReduce) Map(key string, data utils.M, conditions []*MapCond) {
	//以条件分类数据
	for _, mc := range conditions {
		if mc.Valid(data) {
			map_data, ok := u.maps[mc]
			if !ok {
				map_data = MapData{}
			}

			list, ok := map_data[key]
			if !ok {
				list = []*MapValue{}
			}
			map_item := &MapValue{
				Value:    data[key],
				Datetime: NewStatDate(data[u.timeColumn]),
			}
			list = append(list, map_item)
			map_data[key] = list
			u.maps[mc] = map_data
		}
	}
}

//最后返回MAP后的计算结果
func (u *MapReduce) Reduce() CondReduceResult {
	condReduceResult := CondReduceResult{}
	//循环条件数据
	for cond, map_data := range u.maps {
		reduce := ReduceResult{}
		//循环当前条件数据的数据汇总字段,并累计结果
		for key, v := range map_data {
			count := 0
			value := 0.0
			for _, c := range v {
				count += 1
				value += utils.ConvertFloat(c.Value)
			}
			reduce[key] = &ReduceValue{
				count,
				value,
			}
		}
		condReduceResult[cond.Name] = reduce
	}
	return condReduceResult
}

//返回带时间周期的计算结果集 以年,月,周,日的统计结果
func (u *MapReduce) ReduceDate() CondDateReduceResult {
	condReduceResult := CondDateReduceResult{}
	//循环条件数据
	for cond, mapData := range u.maps {
		reduce := DateReduceResult{}
		//循环当前条件数据的数据汇总字段,并累计结果
		for key, v := range mapData {
			for _, c := range v {
				u.reduceDateValue(key, c, reduce)
			}
		}
		condReduceResult[cond.Name] = reduce
	}
	return condReduceResult
}

//计算所有时间段数据
func (u *MapReduce) reduceDateValue(key string, val *MapValue, dateVal DateReduceResult) {
	dataColl, ok := dateVal[val.Datetime.Year]
	if !ok {
		dataColl = map[string]map[int]ReduceResult{}
		dateVal[val.Datetime.Year] = dataColl
	}
	//年汇总
	yearly, ok := dataColl["Yearly"]
	if !ok {
		yearly = map[int]ReduceResult{}
		dataColl["Yearly"] = yearly
	}

	yearlyResult, ok := yearly[val.Datetime.Year]
	if !ok {
		yearlyResult = ReduceResult{}
		yearly[val.Datetime.Year] = yearlyResult
	}

	if k, ok := yearlyResult[key]; ok {
		k.Add(&ReduceValue{1, utils.ConvertFloat(val.Value)})
	} else {
		yearlyResult[key] = &ReduceValue{1, utils.ConvertFloat(val.Value)}
	}

	//月汇总
	monthly, ok := dataColl["monthly"]
	if !ok {
		monthly = map[int]ReduceResult{}
		dataColl["monthly"] = monthly
	}

	monthlyResult, ok := monthly[val.Datetime.Month]
	if !ok {
		monthlyResult = ReduceResult{}
		monthly[val.Datetime.Month] = monthlyResult
	}

	if k, ok := monthlyResult[key]; ok {
		k.Add(&ReduceValue{1, utils.ConvertFloat(val.Value)})
	} else {
		monthlyResult[key] = &ReduceValue{1, utils.ConvertFloat(val.Value)}
	}

	//周汇总
	weekly, ok := dataColl["weekly"]
	if !ok {
		weekly = map[int]ReduceResult{}
		dataColl["weekly"] = weekly
	}

	weeklyResult, ok := weekly[val.Datetime.Week]
	if !ok {
		weeklyResult = ReduceResult{}
		weekly[val.Datetime.Week] = weeklyResult
	}

	if k, ok := weeklyResult[key]; ok {
		k.Add(&ReduceValue{1, utils.ConvertFloat(val.Value)})
	} else {
		weeklyResult[key] = &ReduceValue{1, utils.ConvertFloat(val.Value)}
	}

	//年份天数汇总

	daily, ok := dataColl["daily"]
	if !ok {
		daily = map[int]ReduceResult{}
		dataColl["daily"] = daily
	}

	dailyResult, ok := daily[val.Datetime.YearDay]
	if !ok {
		dailyResult = ReduceResult{}
		daily[val.Datetime.YearDay] = dailyResult
	}

	if k, ok := dailyResult[key]; ok {
		k.Add(&ReduceValue{1, utils.ConvertFloat(val.Value)})
	} else {
		dailyResult[key] = &ReduceValue{1, utils.ConvertFloat(val.Value)}
	}
}
