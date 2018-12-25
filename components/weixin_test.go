package components

import (
	"../utils"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"
)

func TestWeixin_SignJsTicket(tx *testing.T) {
	data := map[string]interface{}{
		"url":          "123123",
		"noncestr":     utils.RandStr(16, nil),
		"timestamp":    time.Now().Unix(),
		"jsapi_ticket": "sadfsadfasdfasdfasdfasdfasdfasdo8392u4o23kjrnweeifu",
	}
	wx := NewWeixin("", "")
	sign := wx.SignJsTicket(data)
	fmt.Println(sign)
}

func TestWeixin_GetMedia(t *testing.T) {
	wx := NewWeixin("wx08ff51936e90b264", "4be28952f7d2af802ac0d07fe6e9ccd1")
	access_token, _ := wx.GetAccessToken()

	data, err := wx.GetMedia(access_token.AccessToken, "D3teZqcjy1yfoSCVbaQ08938xkYA3zv7J-F1WKS_QaN997CJAFAv6gcxCp951unS")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(data)
}

func getRegText() {
	dis := "attachment; filename=\"MEDIA_ID.jpg\""
	reg := regexp.MustCompile(`filename="(.+)"`)
	list := reg.FindStringSubmatch(dis)
	fmt.Println(list[1])
}
