package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clakeboy/golib/utils"
	"github.com/xwb1989/sqlparser"
	"go.etcd.io/bbolt"
)

type (
	errMsg error
)

var timeout = 30 * time.Second

type StormCmdLine struct {
	cmdInput textinput.Model //命令输入框
	err      error
	currCmd  string //当前命令
	timer    timer.Model
	exeCmd   bool     //执行中
	exeDone  bool     //执行完成
	history  *History //执行结果记录
	keyIdx   int      //当前key
}

func NewStormCmdLine() *StormCmdLine {
	ti := textinput.New()
	ti.Placeholder = "SQL查询语句"
	ti.Focus()
	ti.CharLimit = 512
	ti.ShowSuggestions = true
	ti.SetSuggestions([]string{"help", "select ", "from ", "limit ", "show tables"})
	return &StormCmdLine{
		cmdInput: ti,
		history:  NewHistory(50),
	}
}

func (s *StormCmdLine) Init() tea.Cmd {
	// s.timer.Init()
	return textinput.Blink
}

func (s *StormCmdLine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		s.timer, cmd = s.timer.Update(msg)
		return s, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		s.timer, cmd = s.timer.Update(msg)
		return s, cmd
	case timer.TimeoutMsg:
		s.err = fmt.Errorf("执行超时,超过30秒")
		s.exeCmd = false
		return s, s.cmdInput.Focus()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return s, tea.Quit
		case tea.KeyCtrlC:
			s.exeCmd = false
			s.err = fmt.Errorf("手动停止执行")
			return s, s.timer.Stop()
		case tea.KeyUp:
			if s.history.Len()-1 >= s.keyIdx {
				s.cmdInput.SetValue(s.history.Index(s.history.Len() - 1 - s.keyIdx).Command)
			}
			s.keyIdx++
			return s, nil
		case tea.KeyEnter:
			s.currCmd = s.cmdInput.Value()
			if s.currCmd == "" {
				return s, nil
			}
			s.exeCmd = true
			s.exeDone = false
			s.cmdInput.SetValue("")
			tti := timer.NewWithInterval(timeout, time.Millisecond)
			s.timer = tti
			go s.exeSql()
			return s, tti.Init()
		}

	// We handle errors just like any other message
	case errMsg:
		s.err = msg
		return s, nil
	}

	if s.exeDone {
		s.exeDone = false
		s.exeCmd = false
		s.err = nil
		return s, tea.Batch(s.timer.Stop(), s.cmdInput.Focus())
	}
	s.cmdInput, cmd = s.cmdInput.Update(msg)
	return s, cmd
}

func (s *StormCmdLine) View() string {
	var view string
	if s.exeCmd {
		view = s.ViewExecute()
	} else {
		view = s.cmdInput.View()
	}

	return fmt.Sprintf(
		"%s%s%s\n\n%s",
		s.history.View(),
		s.ViewError(),
		view,
		"(esc to quit)",
	) + "\n"
}

func (s *StormCmdLine) ViewError() string {
	if s.err != nil {
		return s.err.Error() + "\n\n"
	}
	return ""
}

func (s *StormCmdLine) ViewExecute() string {
	return fmt.Sprintf("正在执行中... %s", (timeout-s.timer.Timeout)) + "\n"
}

func (s *StormCmdLine) exeSql() {
	defer func() {
		s.keyIdx = 0
	}()

	if s.currCmd == "show tables" {
		var list []string
		db.Bolt.View(func(tx *bbolt.Tx) error {
			return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				list = append(list, string(name))
				return nil
			})
		})
		res := &HistoryResult{
			Command:     s.currCmd,
			Result:      strings.Join(list, "\n"),
			ExecuteTime: (timeout - s.timer.Timeout),
		}
		s.history.Set(res)
		s.exeDone = true
		return
	}

	stmt, err := sqlparser.Parse(s.currCmd)
	if err != nil {
		// s.err = fmt.Errorf("sql 语句解释出错: %v", err)
		res := &HistoryResult{
			Command:     s.currCmd,
			Result:      fmt.Sprintf("sql 语句解释出错: %v", err),
			ExecuteTime: (timeout - s.timer.Timeout),
		}
		s.history.Set(res)
		s.exeDone = true
		return
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		// _ = stmt.From
		var offset int
		var limit int
		tableName := sqlparser.String(stmt.From)
		// columns := sqlparser.String(stmt.SelectExprs)
		if stmt.Limit != nil {
			offset, _ = strconv.Atoi(sqlparser.String(stmt.Limit.Offset))
			limit, _ = strconv.Atoi(sqlparser.String(stmt.Limit.Rowcount))
		}
		list, _ := queryStormDb(tableName, offset, limit)
		res := &HistoryResult{
			Command:     s.currCmd,
			Result:      strings.Join(list, "\n"),
			ExecuteTime: (timeout - s.timer.Timeout),
		}
		s.history.Set(res)
	case *sqlparser.Insert:
		res := &HistoryResult{
			Command:     s.currCmd,
			Result:      fmt.Sprintf("%s %s", stmt.Columns, stmt.Rows),
			ExecuteTime: (timeout - s.timer.Timeout),
		}
		s.history.Set(res)
	}
	s.exeDone = true
}

func queryStormDb(tableName string, offset int, limit int) ([]string, error) {
	buckets := strings.Split(tableName, ".")
	var list []string
	err := db.Bolt.View(func(tx *bbolt.Tx) error {
		var b *bbolt.Bucket
		// b := tx.Bucket([]byte("Account"))
		// b = b.Bucket([]byte("AccountData"))
		for _, name := range buckets {
			if b == nil {
				b = tx.Bucket([]byte(name))
			} else {
				b = b.Bucket([]byte(name))
			}
		}

		// b = b.Bucket([]byte("__storm_index_Id"))
		b.ForEach(func(k, v []byte) error {
			_, err := utils.BytesToInt64(k)
			if err != nil {
				return fmt.Errorf("stop")
			} else {
				// fmt.Printf("key: %d, value: %s\n", idx, string(v))
				list = append(list, string(v))
			}
			return nil
		})
		return nil
	})
	return list, err
}
