package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestMain(t *testing.T) {
	InitCommand()
	savePath := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(CmdOut, "/"), CmdName, "cache")
	fmt.Println("开始获取 golang 模板...")
	fetchGolangFiles(savePath)
	if CmdFront {
		fmt.Printf("开始获取 react-%s 模板...\n", CmdFrontType)
		fetchFrontFiles(savePath)
	}
	fmt.Println("清理下载文件...")
	os.RemoveAll(savePath)
	fmt.Printf("完成项目 [%s] 初始化\n", CmdName)
	fmt.Println("进入项目目录执行 go mod tidy 完成 golang 项目依赖安装")
	if CmdFront {
		fmt.Println("进入项目目录 frontend 执行 npm i 完成前端项目依赖安装")
	}
}

type model[T any] struct {
	packages []T
	names    []string
	index    int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
	callback func(item T, idx int) tea.Msg
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func newModel[T any](total []T, callback func(item T, idx int) tea.Msg) model[T] {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return model[T]{
		packages: total,
		names:    make([]string, len(total)),
		spinner:  s,
		progress: p,
		callback: callback,
	}
}

func (m model[T]) Init() tea.Cmd {
	msgStr := m.callback(m.packages[m.index], m.index)
	m.names[m.index] = msgStr.(string)
	return tea.Batch((func() tea.Msg {
		return msgStr
	}), m.spinner.Tick)
}

func (m model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case PkgMsg:
		if m.index >= len(m.packages)-1 {
			// Everything's been installed. We're done!
			m.done = true
			return m, tea.Quit
		}

		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.packages)-1))

		m.index++
		msgStr := m.callback(m.packages[m.index], m.index)
		m.names[m.index] = msgStr.(string)
		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, msgStr), // print success message above our program
			func() tea.Msg {
				return msgStr
			}, // download the next package
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m model[T]) View() string {
	n := len(m.packages)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render("完成所有文件初始化.\n")
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n-1)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	pkgName := currentPkgNameStyle.Render(m.names[m.index])
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("解压文件 " + pkgName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + pkgCount
}

type PkgMsg string

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
