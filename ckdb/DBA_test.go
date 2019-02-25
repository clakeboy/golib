package ckdb

import (
	"fmt"
	"github.com/clakeboy/golib/utils"
	"reflect"
	"testing"
	"time"
)

var cfg = &DBConfig{
	DBServer:   "168.168.0.10",
	DBPort:     "3306",
	DBName:     "test_db",
	DBUser:     "root",
	DBPassword: "kKie93jgUrn!k",
	DBPoolSize: 200,
	DBIdleSize: 100,
	DBDebug:    true,
}

func TestDBA_Insert(t *testing.T) {
	db, err := NewDBA(cfg)
	if err != nil {
		panic(err)
	}
	tab := db.Table("t_vehicle_info")
	list, err := tab.Limit(10, 1).Query().Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", list[0])
}

func TestDBA_InsertMulti(t *testing.T) {
	db, err := NewDBA(cfg)
	if err != nil {
		panic(err)
	}

	var dataList []interface{}

	for i := 0; i < 10; i++ {
		dataList = append(dataList, utils.M{
			"name": fmt.Sprintf("test_%d", i+1),
			"age":  i + 1,
		})
	}

	tab := db.Table("test")
	res, ok := tab.InsertMulti(dataList)

	fmt.Println(res, ok)
}

func TestDBA_WhereRecursion(t *testing.T) {

	dba, err := NewDBA(cfg)
	if err != nil {
		panic(err)
	}

	var params struct {
		TagName        string `json:"tag_name"`
		TagNum         int    `json:"tag_num"`
		TagCreatedDate int    `json:"tag_created_date"`
	}

	params.TagName = "我去"
	params.TagNum = 3
	params.TagCreatedDate = int(time.Now().Unix())

	data, err := dba.ConvertData(&params)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

//扫描数据到传入的类型
func scanType(scans []interface{}, columns []string, i interface{}) interface{} {
	if i == nil {
		return scanMap(scans, columns)
	}
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Ptr:
		return scanStruct(t.Elem(), scans, columns)
	case reflect.Struct:
		return scanStruct(t, scans, columns)
	case reflect.Map:
		fallthrough
	default:
		return scanMap(scans, columns)
	}
}

//扫描数据到结构体
func scanStruct(t reflect.Type, scans []interface{}, columns []string) interface{} {
	obj := reflect.New(t).Interface()
	objV := reflect.ValueOf(obj)
	for i, colName := range columns {
		idx := findTagOfStruct(t, colName)
		if idx != -1 {
			scans[i] = objV.Field(idx).Interface()
		}
	}
	return obj
}

//在结构体查找TAG值是否存在
func findTagOfStruct(t reflect.Type, colName string) int {
	for i := 0; i < t.NumField(); i++ {
		val, ok := t.Field(i).Tag.Lookup(colName)
		if ok && val == colName {
			return i
		}
	}
	return -1
}

//扫描数据到MAP 默认 utils.M
func scanMap(scans []interface{}, columns []string) interface{} {
	obj := utils.M{}

	return obj
}
