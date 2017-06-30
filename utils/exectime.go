package utils

import (
	"time"
	"fmt"
)

type ExecTime struct {
	start time.Time
	end time.Time
}

func NewExecTime() *ExecTime {
	return &ExecTime{}
}

func (this *ExecTime) Start() {
	this.start = time.Now()
}

func (this *ExecTime) End(print bool) time.Duration {
	this.end = time.Now()
	diff := this.end.Sub(this.start)
	if print {
		fmt.Println("exec time：",diff)
	}
	return diff
}
