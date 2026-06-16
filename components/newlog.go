package components

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// SlogFileHandler 是一个 slog.Handler 代理，
// 将日志写入 <当前运行目录>/logs/<subDir>/<YYYY-MM-DD>.log。
// 每次写入时检查日期，跨天自动切换到新的日志文件。
// 输出格式：
//
//	<2006-01-02 15:04:05> <LEVEL> <source>
//	<message> <key1>=<val1> <key2>=<val2> ...
type SlogFileHandler struct {
	mu   sync.Mutex
	dir  string
	file *os.File
	day  string
}

// NewSlogFile 创建一个 *slog.Logger，其 handler 写入 ./logs/<subDir>/<YYYY-MM-DD>.log。
// subDir 是 logs 下的子目录名；传入空字符串则直接使用 logs 目录。
func NewSlogFile(subDir string) (*slog.Logger, error) {
	h, err := newSlogFileHandler(subDir)
	if err != nil {
		return nil, err
	}
	return slog.New(h), nil
}

func newSlogFileHandler(subDir string) (*SlogFileHandler, error) {
	h := &SlogFileHandler{dir: subDir}
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.rotateLocked(); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *SlogFileHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }

func (h *SlogFileHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.maybeRotateLocked(); err != nil {
		return err
	}
	_, err := h.file.Write(formatRecord(r, nil))
	return err
}

func (h *SlogFileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogDerived{parent: h, attrs: attrs}
}

func (h *SlogFileHandler) WithGroup(name string) slog.Handler {
	return &slogDerived{parent: h, groups: []string{name}}
}

func (h *SlogFileHandler) maybeRotateLocked() error {
	today := time.Now().Format("2006-01-02")
	if today == h.day {
		return nil
	}
	return h.rotateToLocked(today)
}

func (h *SlogFileHandler) rotateLocked() error {
	return h.rotateToLocked(time.Now().Format("2006-01-02"))
}

func (h *SlogFileHandler) rotateToLocked(today string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir := filepath.Join(cwd, "logs", h.dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if h.file != nil {
		_ = h.file.Close()
	}
	f, err := os.OpenFile(filepath.Join(dir, today+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	h.file = f
	h.day = today
	return nil
}

// formatRecord 格式化单条日志。
// groups 是累积的 group 路径（来自 logger.WithGroup(...)），attr key 会自动加上 "<g1>.<g2>." 前缀。
func formatRecord(r slog.Record, groups []string) []byte {
	var src string
	if r.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{r.PC})
		if frame, _ := frames.Next(); frame.File != "" {
			src = filepath.Base(frame.File)
		}
	}

	var buf []byte
	buf = append(buf, r.Time.Format("2006-01-02 15:04:05")...)
	buf = append(buf, ' ')
	buf = append(buf, r.Level.String()...)
	if src != "" {
		buf = append(buf, ' ')
		buf = append(buf, src...)
	}
	buf = append(buf, '\n')
	buf = append(buf, r.Message...)

	r.Attrs(func(a slog.Attr) bool {
		buf = append(buf, ' ')
		for _, g := range groups {
			buf = append(buf, g...)
			buf = append(buf, '.')
		}
		buf = appendAttr(buf, a)
		return true
	})
	buf = append(buf, '\n')
	return buf
}

func appendAttr(buf []byte, a slog.Attr) []byte {
	buf = append(buf, a.Key...)
	buf = append(buf, '=')
	return appendValue(buf, a.Value)
}

func appendValue(buf []byte, v slog.Value) []byte {
	switch v.Kind() {
	case slog.KindString:
		return appendString(buf, v.String())
	case slog.KindInt64:
		return strconv.AppendInt(buf, v.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(buf, v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.AppendFloat(buf, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		return strconv.AppendBool(buf, v.Bool())
	case slog.KindDuration:
		return append(buf, v.Duration().String()...)
	case slog.KindTime:
		return append(buf, v.Time().Format(time.RFC3339)...)
	case slog.KindGroup:
		for _, av := range v.Group() {
			buf = appendAttr(buf, av)
			buf = append(buf, ' ')
		}
		return buf
	case slog.KindLogValuer:
		return appendValue(buf, v.LogValuer().LogValue())
	case slog.KindAny:
		if e, ok := v.Any().(error); ok {
			return append(buf, e.Error()...)
		}
		return append(buf, fmt.Sprint(v.Any())...)
	}
	return buf
}

func appendString(buf []byte, s string) []byte {
	if needsQuoting(s) {
		return strconv.AppendQuote(buf, s)
	}
	return append(buf, s...)
}

func needsQuoting(s string) bool {
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			return true
		}
	}
	return false
}

// slogDerived 包装父 SlogFileHandler，累积 attrs 和 group 路径。
type slogDerived struct {
	parent *SlogFileHandler
	attrs  []slog.Attr
	groups []string
}

func (d *slogDerived) Enabled(_ context.Context, _ slog.Level) bool { return true }

func (d *slogDerived) Handle(_ context.Context, r slog.Record) error {
	d.parent.mu.Lock()
	defer d.parent.mu.Unlock()
	if err := d.parent.maybeRotateLocked(); err != nil {
		return err
	}
	if len(d.attrs) > 0 {
		r.AddAttrs(d.attrs...)
	}
	_, err := d.parent.file.Write(formatRecord(r, d.groups))
	return err
}

func (d *slogDerived) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := make([]slog.Attr, 0, len(d.attrs)+len(attrs))
	merged = append(merged, d.attrs...)
	merged = append(merged, attrs...)
	return &slogDerived{parent: d.parent, attrs: merged, groups: d.groups}
}

func (d *slogDerived) WithGroup(name string) slog.Handler {
	groups := make([]string, 0, len(d.groups)+1)
	groups = append(groups, d.groups...)
	groups = append(groups, name)
	return &slogDerived{parent: d.parent, attrs: d.attrs, groups: groups}
}
