package utils

import (
	"fmt"
	"testing"
)

type TestUser struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	*JsonParse `json:"-"`
}

func TestJsonParse_ParseJson(t *testing.T) {
	u := &TestUser{
		Name: "clake",
		Age:  18,
	}
	fmt.Println(u)
	u.JsonParse = NewJsonParse(u)
	fmt.Println(u)
	fmt.Println(u.ToJsonString())
	fmt.Println(u.Name, u.Age)
}
