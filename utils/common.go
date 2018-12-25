package utils

import (
	"../utils/uuid"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type M map[string]interface{}

func (m *M) ToJson() []byte {
	data, err := json.Marshal(m)
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
	err := json.Unmarshal(data, m)
	if err != nil {
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

	for i, caps := range list {
		list[i] = UcFirst(caps)
	}

	return strings.Join(list, "")
}

//首字母大写
func UcFirst(str string) string {
	first := str[0:1]
	long := str[1:]
	return strings.ToUpper(first) + long
}

const randTable = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"

//生成随机字符串
func RandStr(number int, r_table interface{}) string {
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
	for i := 0; i < number; i++ {
		str = append(str, string(table[rand.Intn(rand_len)]))
	}

	return strings.Join(str, "")
}

//结构转map
func Struct2Map(obj interface{}, fields []string) map[string]interface{} {
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

	for i := 0; i < t.NumField(); i++ {
		if fields != nil {
			flag := StringIndexOf(fields, t.Field(i).Tag.Get("json"))
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
func Map2Struct(obj interface{}, stu interface{}) error {
	json_data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = json.Unmarshal(json_data, stu)
	if err != nil {
		return err
	}
	return nil
}

//创建UUID
func CreateUUID(step bool) string {
	if runtime.GOOS == "windows" {
		ui := uuid.Must(uuid.NewV4())
		if step {
			return ui.String()
		} else {
			return fmt.Sprintf("%x", ui.Bytes())
		}

	} else {
		f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
		b := make([]byte, 16)
		f.Read(b)
		f.Close()
		var ui string
		if step {
			ui = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		} else {
			ui = fmt.Sprintf("%x", b[:])
		}

		return ui
	}
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
		bytes[3], bytes[2], bytes[1], bytes[0])
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
func YN(condition bool, yes interface{}, no interface{}) interface{} {
	if condition {
		return yes
	}
	return no
}

//转换任意数据为 float64
func ConvertFloat(c interface{}) float64 {
	rv := reflect.ValueOf(c)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return float64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return rv.Float()
	case reflect.String:
		fl, err := strconv.ParseFloat(rv.String(), 64)
		if err != nil {
			return 0
		}
		return fl
	}
	return 0
}

//设置时间回调
func SetTimeout(step time.Duration, callback func()) {
	go func() {
		time.Sleep(step)
		callback()
	}()
}

var commands = map[string]string{
	"windows": "cmd /c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

func OpenBrowse(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmd := exec.Command(run, uri)
	return cmd.Start()
}
