package task

import (
	"testing"
	"time"
	"fmt"
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