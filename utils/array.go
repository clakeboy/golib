package utils

import (
	"reflect"
	"errors"
)

func StringIndexOf(arr []string, search string) int {
	for i,v := range arr {
		if search == v {
			return i
		}
	}

	return -1
}

func MapKeys(m map[string]interface{}) []string{
	var keys []string
	for i,_ := range m {
		keys = append(keys,i)
	}

	return keys
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
