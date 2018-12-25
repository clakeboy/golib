package utils

import (
	"fmt"
	"testing"
)

type TestUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	JsonParse
}

func TestJsonParse_ParseJson(t *testing.T) {
	u := &TestUser{
		Name: "clake",
		Age:  18,
	}

	fmt.Println(u.ToJsonString())
	fmt.Println(u.Name, u.Age)
}
