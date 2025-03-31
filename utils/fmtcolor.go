package utils

import "fmt"

/**
  \033[0m   关闭所有属性
  \033[1m   设置高亮度
  \03[4m   下划线
  \033[5m   闪烁
  \033[7m   反显
  \033[8m   消隐
  \033[30m   --   \033[37m   设置前景色
  \033[40m   --   \033[47m   设置背景色
  \033[nA   光标上移n行
  \03[nB   光标下移n行
  \033[nC   光标右移n行
  \033[nD   光标左移n行
  \033[y;xH设置光标位置
  \033[2J   清屏
  \033[K   清除从光标到行尾的内容
  \033[s   保存光标位置
  \033[u   恢复光标位置
  \033[?25l   隐藏光标
  \33[?25h   显示光标
*/

const (
	//前景色
	FBLACK   = 30
	FRED     = 31
	FGREEN   = 32
	FYELLOW  = 33
	FBLUE    = 34
	FFUCHSIA = 35
	FCYAN    = 36
	FWHITE   = 37
	//背景色
	BBLACK   = 40
	BRED     = 41
	BGREEN   = 42
	BYELLOW  = 43
	BBLUE    = 44
	BFUCHSIA = 45
	BCYAN    = 46
	BWHITE   = 47
	//显示模式
	SDEF = 0 //终端默认
	SHGH = 1 //高亮显示
	SUNL = 4 //下划线
	SLL  = 5 //闪烁
	SB   = 7 //反白显示
	SNO  = 8 //不可见
)

func FmtColor(str string, colors ...int) string {
	plen := len(colors)
	switch plen {
	case 1:
		return formatStringColor(str, colors[0], 0, SDEF)
	case 2:
		return formatStringColor(str, colors[0], colors[1], SDEF)
	case 3:
		return formatStringColor(str, colors[0], colors[1], colors[2])
	default:
		return str
	}
}

func formatStringColor(str string, foregColor int, backColor int, showMode int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, showMode, backColor, foregColor, str, 0x1B)
}
