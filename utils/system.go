package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

//生成PID文件
func WritePid(pidName string) {
	if osName := runtime.GOOS; osName != "windows" {
		pid := strconv.Itoa(os.Getpid())
		ioutil.WriteFile(pidName, []byte(pid), 0755)
	}
}

//收到退出信号时处理程序关闭
func ExitApp(out chan os.Signal, callback func(os.Signal)) {
	signal.Notify(out, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case s := <-out:
		if callback != nil {
			callback(s)
		}
		fmt.Println(s)
		os.Exit(0)
	}
}
