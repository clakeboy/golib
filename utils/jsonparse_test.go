package utils

import (
	"testing"
	"fmt"
)

type Test struct {
	Name string `json:"name"`
	Age int `json:"age"`
	*JsonParse
}

func TestJsonParse_ParseJson(t *testing.T) {
	ts := &Test{
		Name:"clake",
		Age:11,
		JsonParse:&JsonParse{},
	}

	fmt.Println(ts.ToJsonString())
}