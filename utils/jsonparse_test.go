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
	f := 2.3482922

	d := Round(f,4)
	fmt.Println(d)
}