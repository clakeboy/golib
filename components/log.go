package components

import (
	"os"
	"ck_go_lib/utils"
	"time"
	"fmt"
	"io"
)

type SysLog struct {
	Prefix string
}

func NewSysLog(prefix string) *SysLog {
	return &SysLog{Prefix:prefix}
}

func (l *SysLog) Write(p []byte) (n int, err error) {
	err = os.MkdirAll("./logs/",0755)
	if err != nil {
		panic(err)
	}
	file_name := fmt.Sprintf("%s%v.log",l.Prefix,utils.FormatTime("YY-MM-DD",time.Now().Unix()))

	err = l.WriteFile("./logs/"+file_name,p,0755)
	if err != nil {
		return 0,err
	}
	return len(p),nil
}

func (l *SysLog) WriteFile(filename string,data []byte,perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func (l *SysLog) checkFileIsExist(filename string) (bool) {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}