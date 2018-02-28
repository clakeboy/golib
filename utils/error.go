package utils

import (
	"fmt"
	"errors"
)
//添加新的错误信息
func Error(msg string,err error) error {
	if err == nil {
		return errors.New(msg)
	}
	return errors.New(fmt.Sprintf("%s(org:%s)",msg,err.Error()))
}