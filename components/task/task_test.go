package task

import (
	"testing"
	"time"
	"fmt"
	"ck_go_lib/utils"
)

func TestNewManagement(t *testing.T) {
	lastDate := time.Now()
	time.Sleep(time.Second * 2)
	currentData := time.Now()

	fmt.Println(currentData.Sub(lastDate).Seconds())
	fmt.Println(int(currentData.Sub(lastDate).Minutes()))
	fmt.Println(currentData.Sub(lastDate).Hours())

	var ss time.Time
	var bol bool
	fmt.Println(ss.IsZero(),bol)
}

func TestManagement_Start(t *testing.T) {
	taskService := NewManagement()
	//every second execute func
	taskService.AddTaskString("*/1 * * * * *", func(item *Item) bool {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"Every second execute")
		return true
	},nil)
	//10 second execute func
	taskService.AddTaskString("10 * * * * *", func(item *Item) bool {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"10 second execute")
		return true
	}, func(item *Item) {
		fmt.Println(utils.FmtColor("10 second callback function",utils.FYELLOW))
	})
	taskService.AddTaskString("* */1 * * * *", func(item *Item) bool {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"Every one minute execute")
		return true
	},nil)
	taskService.AddTaskString("1 18 11 * * *", func(item *Item) bool {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"11:18:01 execute")
		return true
	},nil)
	taskService.Start()
	fmt.Println("start")
	out := make(chan bool,1)
	<-out
}