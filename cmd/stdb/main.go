package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/asdine/storm/v3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clakeboy/golib/utils"
	"go.etcd.io/bbolt"
)

var cmdTimeout int
var cmdFile string
var cmdRead bool
var db *storm.DB

func main() {
	flag.IntVar(&cmdTimeout, "t", 30, "默认30秒执行超时")
	flag.BoolVar(&cmdRead, "read", true, "只读模式打开")
	flag.Parse()
	// fmt.Println(flag.Args())
	cmdFile = flag.Arg(0)

	if !utils.Exist(cmdFile) {
		fmt.Println("db file not exist", cmdFile)
		os.Exit(1)
	}
	var err error
	db, err = storm.Open(cmdFile, storm.BoltOptions(0, &bbolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: cmdRead,
	}))
	if err != nil {
		fmt.Println("open db file error:", err)
		os.Exit(1)
	}
	p := tea.NewProgram(NewStormCmdLine())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
