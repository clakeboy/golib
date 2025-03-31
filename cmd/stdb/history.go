package main

import (
	"fmt"
	"strings"
	"time"
)

type HistoryResult struct {
	Command     string        //执行命令
	ExecuteTime time.Duration //执行时间
	Result      string        //执行结果
	CreatedDate time.Time     //开始时间
}

type History struct {
	history []*HistoryResult //执行记录集
	limit   int              //记录存放总量
}

func NewHistory(limit int) *History {
	return &History{
		limit: limit,
	}
}

func (h *History) Len() int { return len(h.history) }

func (h *History) Index(idx int) *HistoryResult {
	if idx >= h.Len() || idx < 0 {
		return nil
	}

	return h.history[idx]
}

func (h *History) Set(s *HistoryResult) {
	s.CreatedDate = time.Now()
	if len(h.history) >= h.limit {
		h.history = h.history[1:]
	}
	h.history = append(h.history, s)
}

func (h *History) View() string {
	if h.Len() <= 0 {
		return ""
	}
	var list []string
	for _, s := range h.history {
		str := fmt.Sprintf("> %s\n%s\n\n%s",
			s.Command,
			fmt.Sprintf("> 执行时间 %s, 完成时间: %s", s.ExecuteTime, s.CreatedDate.Format("2006-01-02 15:04:05")),
			s.Result)
		list = append(list, str)
	}
	return strings.Join(list, "\n\n") + "\n\n"
}
