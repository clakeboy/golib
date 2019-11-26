package components

import (
	"bytes"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

type SysLog struct {
	Prefix string
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func NewSysLog(prefix string) *SysLog {
	return &SysLog{Prefix: prefix}
}

func (l *SysLog) Write(p []byte) (n int, err error) {
	err = os.MkdirAll("./logs/", 0755)
	if err != nil {
		panic(err)
	}
	file_name := fmt.Sprintf("%s%v.log", l.Prefix, utils.FormatTime("YY-MM-DD", time.Now().Unix()))

	err = l.WriteFile("./logs/"+file_name, p, 0755)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (l *SysLog) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func (l *SysLog) checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//记录错误
func (l *SysLog) Error(err interface{}) {
	stack_str := stack(4)
	msg := fmt.Sprintf("[ERROR][%s] panic recovered:\n%s\n%s\n", time.Now().Format("2006-01-02 15:04:05"), err, stack_str)
	l.Write([]byte(msg))
}

func (l *SysLog) Info(msg string) {
	msg_str := fmt.Sprintf("[INFO][%s]:\n%s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	l.Write([]byte(msg_str))
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
