package jufmt

import (
	gofmt "fmt"
	"os"
)

func Printf(format string, a ...any) {
	var prefix string
	if printTime {
		prefix += GetNowTimeMs() + " "
	}
	if printTrace {
		prefix += GetTrace(2) + " "
	}
	gofmt.Fprint(os.Stdout, prefix+BrightWhite.Sprintf(format, a...))
	//trace := GetTrace(2)
	//gofmt.Fprintf(os.Stdout, "%s %s %s", GetNowTimeMs(), trace, BrightWhite.Sprintf(format, a...))
}
func Println(a ...any) {
	var prefix string
	if printTime {
		prefix += GetNowTimeMs() + " "
	}
	if printTrace {
		prefix += GetTrace(2) + " "
	}
	gofmt.Fprint(os.Stdout, prefix+BrightWhite.Sprintln(a...))
	//trace := GetTrace(2)
	//gofmt.Fprintf(os.Stdout, "%s %s %s", GetNowTimeMs(), trace, BrightWhite.Sprintln(a...))
}
func Sprintf(format string, a ...any) string {
	return gofmt.Sprintf(format, a...)
}
func Errorf(format string, a ...any) error {
	return gofmt.Errorf(format, a...)
}

// TracePrintln 打印多步调用位置，exStep 是额外打印多少步调用信息，exStep < 0 等同于默认值 0
func TracePrintln(exStep int, a ...any) {
	BrightWhite.TracePrintln(exStep, a...)
}

// TracePrintf 打印多步调用位置，exStep 是额外打印多少步调用信息，exStep < 0 等同于默认值 0
func TracePrintf(exStep int, format string, a ...any) {
	BrightWhite.TracePrintf(exStep, format, a...)
}
