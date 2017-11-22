package utils

import (
	"fmt"
	"testing"
)

func TestHttpClient_Post(t *testing.T) {
	data := M{
		"license_no": "渝B6S919",
	}
	res,_ := HttpPostJson("http://localhost:7908/serv/car/push_car_license", data)
	fmt.Println(res)
}

func TestHttpClient_Get(t *testing.T) {
	//da6ea0db35bd51ae0c2bab87ad1c9bf8 3ac179b21f884cf155d6e69b70b66d45
	//data := M{
	//	"call_id": "e678e293ce724a2ba11789b0c4cb15f9",
	//}
	//res,_ := HttpPostJson("http://localhost:7908/serv/car/query_car", data)
	//fmt.Println(res.ToJsonString())

	//content, err := HttpGet("http://m.cqtcxx.com/parkslist.aspx?areaid=522")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(content)

	client := NewHttpClient()
	_,err := client.Get("https://www.baidu.com")
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(res))
	resp := client.GetLastResponse()
	fmt.Printf("%+v",resp.Cookie.Cookies[0].Value)
}
