package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func StringIndexOf(arr []string, search string) int {
	for i, v := range arr {
		if search == v {
			return i
		}
	}

	return -1
}

func MapKeys(m map[string]interface{}) []string {
	var keys []string
	for i, _ := range m {
		keys = append(keys, i)
	}

	return keys
}

func MapDump(m M) ([]string, []interface{}) {
	keys := make([]string, len(m))
	values := make([]interface{}, len(m))
	idx := 0
	for k, v := range m {
		keys[idx] = k
		values[idx] = v
		idx++
	}
	return keys, values
}

func Contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in")
}

func PrintMap(obj map[string]interface{}, step string) {
	fmt.Println(step, "{")
	for k, v := range obj {
		switch v.(type) {
		case map[string]interface{}, map[string]string:
			fmt.Print(step, k, ": ")
			PrintMap(v.(map[string]interface{}), step+"     ")
		default:
			fmt.Println(step, k, ": ", v)
		}
	}
	fmt.Println(step, "}")
}
