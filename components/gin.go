package components

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"fmt"
	"reflect"
	"ck_go_lib/utils"
	"github.com/DeanThompson/ginpprof"
)
//得到POST原始数据
func GetProperty(c *gin.Context) []byte {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	sid := c.Request.Header.Get("CK-Pro-S")
	if sid != "" {
		enc := utils.NewAes(sid)
		data, err = enc.Decrypt(data)
		if err != nil {
			panic(err)
		}
	}
	return data
}

//调用Controller 的 Action 方法
func CallAction(i interface{}, c *gin.Context) {
	t := reflect.TypeOf(i)
	method, ok := t.MethodByName("Action" + utils.Under2Hump(c.Param("action")))
	if ok {
		var params []reflect.Value
		if method.Type.NumIn() == 2 {
			params = make([]reflect.Value,1)
			res := GetProperty(c)
			params[0] = reflect.ValueOf(res)
		}
		v := reflect.ValueOf(i)
		list := v.MethodByName(method.Name).Call(params)
		rnflag := true
		msg := "ok"

		args_len := len(list)

		buildOutput(list[args_len-1],&rnflag,&msg)

		c.JSON(200, utils.ApiResult(rnflag, msg, list[0].Interface()))
	} else {
		fmt.Println("not found action")
	}
}

func buildOutput(v reflect.Value,rnflag *bool, msg *string) {
	switch v.Type().String() {
	case "bool":
		if !v.Interface().(bool) {
			*rnflag = false
			*msg = "error"
		}
	case "error":
		if v.Interface() != nil {
			*rnflag = false
			*msg = v.Interface().(error).Error()
		}
	}
}

func InitPprof(server *gin.Engine) {
	ginpprof.Wrapper(server)
}
//是否可以跨域调用
func Cross(c *gin.Context,is_cross bool,org string) {
	if is_cross {
		c.Header("Access-Control-Allow-Origin", org)
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Accept, Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, CK-Pro-S")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
}