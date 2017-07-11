package utils

import (
	"strings"
	"math/rand"
	"time"
	"reflect"
	"os"
	"fmt"
	"encoding/json"
	"strconv"
)

type M map[string]interface{}

func (m *M) ToJson() []byte {
	data ,err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return data
}

func (m *M) ToJsonString() string {
	data := m.ToJson()
	if data == nil {
		return ""
	}
	return string(data)
}

func (m *M) ParseJson(data []byte) error {
	err := json.Unmarshal(data,m)
	if err != nil{
		return err
	}
	return nil
}

func (m *M) ParseJsonString(data string) error {
	return m.ParseJson([]byte(data))
}

func Hump2Under(str string) {

}
//下划线转驼峰
func Under2Hump(str string) string {
	list := strings.Split(str, "_")

	for i, cap := range list {
		list[i] = UcFirst(cap)
	}

	return strings.Join(list,"")
}
//首字母大写
func UcFirst(str string) string{
	first := str[0:1]
	long := str[1:]
	return strings.ToUpper(first) + long
}

const randTable = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"
//生成随机字符串
func RandStr(number int,r_table interface{}) string {
	var table string
	if r_table != nil {
		table = r_table.(string)
	} else {
		table = randTable
	}
	time.Sleep(1)
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	str := []string{}
	rand_len := len(table)
	for i:=0;i<number;i++ {
		str = append(str,string(table[rand.Intn(rand_len)]))
	}

	return strings.Join(str,"")
}
//结构转map
func Struct2Map(obj interface{},fields []string) map[string]interface{} {
	var t reflect.Type
	var v reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		t = reflect.TypeOf(obj).Elem()
		v = reflect.ValueOf(obj).Elem()
	} else {
		t = reflect.TypeOf(obj)
		v = reflect.ValueOf(obj)
	}

	var data = make(map[string]interface{})

	for i:=0;i<t.NumField();i++ {
		if fields != nil {
			flag := StringIndexOf(fields,t.Field(i).Tag.Get("json"))
			if flag != -1 {
				data[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()
			}
		} else {
			data[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()
		}
	}

	return data
}
//map 转结构
func Map2Struct(obj interface{},stu interface{}) error {
	json_data,err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = json.Unmarshal(json_data,stu)
	if err != nil {
		return err
	}
	return nil
}
//创建UUID
func CreateUUID(step bool) string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	var uuid string
	if step {
		uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	} else {
		uuid = fmt.Sprintf("%x", b[:])
	}

	return uuid
}


//API输出
func ApiResult(status bool, msg string, i interface{}) *map[string]interface{} {
	res := &map[string]interface{}{"status": status, "msg": msg, "data": i}
	return res
}

// Convert uint to net.IP http://www.outofmemory.cn
func Inet_ntoa(ipnr int64) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return fmt.Sprintf(
		"%v.%v.%v.%v",
		bytes[0],bytes[1],bytes[2],bytes[3])
}

// Convert net.IP to int64 ,  http://www.outofmemory.cn
func Inet_aton(ipnr string) int64 {
	bits := strings.Split(ipnr, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}
//模拟二元操作符功能
func YN(condition bool,yes interface{},no interface{}) interface{} {
	if condition {
		return yes
	}
	return no
}