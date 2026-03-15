package jufmt

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

// ANSI 颜色码
const (
	reset = "\033[0m"
)

// Color 类型
type Color struct {
	code string
}

// Printf 输出带颜色的格式化文本
func (c Color) Printf(format string, a ...interface{}) {
	trace := GetTrace(2)
	fmt.Fprintf(os.Stdout, "%s %s %s", GetNowTimeMs(), trace, c.Sprintf(format, a...))
}

// Println 输出带颜色的文本并换行
func (c Color) Println(a ...interface{}) {
	trace := GetTrace(2)
	fmt.Fprintf(os.Stdout, "%s %s %s", GetNowTimeMs(), trace, c.Sprintln(a...))
}
func (c Color) printf(skip int, format string, a ...interface{}) {
	trace := GetTrace(skip + 2)
	fmt.Fprintf(os.Stdout, "%s %s %s", GetNowTimeMs(), trace, c.Sprintf(format, a...))
}

// println 以默认格式打印参数，输出带颜色的文本并换行
func (c Color) println(skip int, a ...interface{}) {
	trace := GetTrace(skip + 2)
	fmt.Fprintf(os.Stdout, "%s %s %s", GetNowTimeMs(), trace, c.Sprintln(a...))
}

// print 以默认格式打印参数，输出带颜色的文本
func (c Color) print(skip int, a ...interface{}) {
	trace := GetTrace(skip + 2)
	fmt.Fprintf(os.Stdout, "%s %s %s", GetNowTimeMs(), trace, c.Sprint(a...))
}

// Sprintf 返回带颜色的格式化字符串
func (c Color) Sprintf(format string, a ...interface{}) string {
	return c.code + fmt.Sprintf(format, a...) + reset
}
func (c Color) Sprintln(a ...interface{}) string {
	return c.code + fmt.Sprintln(a...) + reset
}
func (c Color) Sprint(a ...interface{}) string {
	return c.code + fmt.Sprint(a...) + reset
}

// 16 个标准颜色变量
var (
	// 普通色

	Black   = Color{"\033[30m"}
	Red     = Color{"\033[31m"}
	Green   = Color{"\033[32m"}
	Yellow  = Color{"\033[33m"}
	Blue    = Color{"\033[34m"}
	Magenta = Color{"\033[35m"}
	Cyan    = Color{"\033[36m"}
	White   = Color{"\033[37m"}

	// 亮色（高亮）

	BrightBlack   = Color{"\033[90m"} // 即 Gray
	BrightRed     = Color{"\033[91m"}
	BrightGreen   = Color{"\033[92m"}
	BrightYellow  = Color{"\033[93m"}
	BrightBlue    = Color{"\033[94m"}
	BrightMagenta = Color{"\033[95m"}
	BrightCyan    = Color{"\033[96m"}
	BrightWhite   = Color{"\033[97m"}

	// 别名

	Gray = BrightBlack
)

func GetNowTimeMs() string {
	return time.Now().Format("15:04:05.000")
}

func GetTrace(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	file = path.Base(file)
	return fmt.Sprintf("%s:%d", file, line)
}
