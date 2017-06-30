package ckdb

import (
	"fmt"
	"testing"
	"time"
)

var cfg = &DBConfig{
	DBServer:   "168.168.0.10",
	DBPort:     "3306",
	DBName:     "pcbx_ddb",
	DBUser:     "root",
	DBPassword: "kKie93jgUrn!k",
	DBPoolSize: 200,
	DBIdleSize: 100,
	DBDebug:    true,
}

func TestDBA_Insert(t *testing.T) {

}

func TestDBA_WhereRecursion(t *testing.T) {

	dba, err := NewDBA(cfg)
	if err != nil {
		panic(err)
	}

	var params struct{
		TagName string `json:"tag_name"`
		TagNum   int `json:"tag_num"`
		TagCreatedDate int `json:"tag_created_date"`
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
