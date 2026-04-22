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

var (
	printTime  = true
	printTrace = true
)

// SetPrintTime enables or disables the inclusion of timestamps in printed output based on the provided flag.
func SetPrintTime(enabled bool) {
	printTime = enabled
}

// SetPrintTrace enables or disables the inclusion of trace information in printed output based on the provided flag.
func SetPrintTrace(enabled bool) {
	printTrace = enabled
}

// Color 类型
type Color struct {
	code string
}

// Printf 输出带颜色的格式化文本
func (c Color) Printf(format string, a ...interface{}) {
	var prefix string
	if printTime {
		prefix += GetNowTimeMs() + " "
	}
	if printTrace {
		prefix += GetTrace(2) + " "
	}
	fmt.Fprint(os.Stdout, prefix+c.Sprintf(format, a...))
}

// Println 输出带颜色的文本并换行
func (c Color) Println(a ...interface{}) {
	var prefix string
	if printTime {
		prefix += GetNowTimeMs() + " "
	}
	if printTrace {
		prefix += GetTrace(2) + " "
	}
	fmt.Fprint(os.Stdout, prefix+c.Sprintln(a...))
}

// TracePrintf 以 format 格式打印参数，输出带颜色的文本并换行，skip < 0 相当于默认值 0
func (c Color) tracePrintf(printTrace bool, skip int, format string, a ...interface{}) {
	if skip < 0 {
		skip = 0
	}
	var prefix string
	if printTime {
		prefix += GetNowTimeMs() + " "
	}
	if printTrace {
		prefix += GetTrace(skip+2) + " "
	}
	fmt.Fprint(os.Stdout, prefix+c.Sprintf(format, a...))
}

// TracePrintln 以默认格式打印参数，输出带颜色的文本并换行，skip < 0 相当于默认值 0
func (c Color) tracePrintln(printTrace bool, skip int, a ...interface{}) {
	if skip < 0 {
		skip = 0
	}
	var prefix string
	if printTime {
		prefix += GetNowTimeMs() + " "
	}
	if printTrace {
		prefix += GetTrace(skip+2) + " "
	}
	fmt.Fprint(os.Stdout, prefix+c.Sprintln(a...))
}

// TracePrintln 打印多步调用位置，exStep 是额外打印多少步调用信息，exStep < 0 等同于默认值 0，这个函数无视全局的 PrintTrace 设置，总是打印调用信息
func (c Color) TracePrintln(exStep int, a ...any) {
	if exStep < 0 {
		exStep = 0
	}
	for i := exStep - 1; i >= 0; i-- {
		c.tracePrintf(true, i+2, "[call %d]\n", i+1)
	}
	c.tracePrintln(true, 1, a...)
}

// TracePrintf 打印多步调用位置，exStep 是额外打印多少步调用信息，exStep < 0 等同于默认值 0，这个函数无视全局的 PrintTrace 设置，总是打印调用信息
func (c Color) TracePrintf(exStep int, format string, a ...any) {
	if exStep < 0 {
		exStep = 0
	}
	for i := exStep - 2; i >= 0; i-- {
		c.tracePrintln(true, i+2, "call")
	}
	c.tracePrintf(true, 1, format, a...)
}

//下面的函数可以用于直接获取彩色字符串，可以使用原生 fmt.Print 函数打印多色字符

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
