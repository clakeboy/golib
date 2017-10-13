package utils

import (
	"testing"
	"fmt"
	"time"
)

func TestM_ToJson(t *testing.T) {
	fmt.Println(time.Now().Nanosecond())
}

func TestInet_ntoa(t *testing.T) {
	run := NewExecTime()
	run.Start()
	ip := "168.168.0.10"
	a := 'a'
	ip_int := Inet_aton(ip)
	fmt.Println(ip,"=",ip_int,a)

	ip_str := Inet_ntoa(ip_int)

	fmt.Println(ip_str)
	run.End(true)
}

func TestBinaryStringToBytes(t *testing.T) {
	a := byte(1)
	fmt.Println(ByteToBinaryString(a))
}
