package components

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SlogFileHandler 是一个 slog.Handler 代理，
// 将日志写入 <当前运行目录>/logs/<subDir>/<YYYY-MM-DD>.log。
// 每次写入时检查日期，跨天自动切换到新的日志文件。
type SlogFileHandler struct {
	mu   sync.Mutex
	dir  string
	file *os.File
	base slog.Handler
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

func (h *SlogFileHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.maybeRotateLocked(); err != nil {
		return err
	}
	return h.base.Handle(ctx, r)
}

func (h *SlogFileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogDerived{parent: h, attrs: attrs}
}

func (h *SlogFileHandler) WithGroup(name string) slog.Handler {
	return &slogDerived{parent: h, group: name}
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
	h.base = slog.NewTextHandler(f, &slog.HandlerOptions{AddSource: true})
	return nil
}

// slogDerived 包装父 SlogFileHandler，累积 attrs/group，
// 每次 Handle 时用父级最新的 base 重建临时 handler。
type slogDerived struct {
	parent *SlogFileHandler
	attrs  []slog.Attr
	group  string
}

func (d *slogDerived) Enabled(ctx context.Context, l slog.Level) bool { return true }

func (d *slogDerived) Handle(ctx context.Context, r slog.Record) error {
	d.parent.mu.Lock()
	defer d.parent.mu.Unlock()
	if err := d.parent.maybeRotateLocked(); err != nil {
		return err
	}
	h := d.parent.base
	if len(d.attrs) > 0 {
		h = h.WithAttrs(d.attrs)
	}
	if d.group != "" {
		h = h.WithGroup(d.group)
	}
	return h.Handle(ctx, r)
}

func (d *slogDerived) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := make([]slog.Attr, 0, len(d.attrs)+len(attrs))
	merged = append(merged, d.attrs...)
	merged = append(merged, attrs...)
	return &slogDerived{parent: d.parent, attrs: merged, group: d.group}
}

func (d *slogDerived) WithGroup(name string) slog.Handler {
	return &slogDerived{parent: d.parent, attrs: d.attrs, group: name}
}
