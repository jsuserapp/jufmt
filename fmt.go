package jufmt

import gofmt "fmt"

func Printf(format string, a ...any) {
	BrightWhite.TracePrintf(1, format, a...)
}
func Println(a ...any) {
	BrightWhite.TracePrintln(1, a...)
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
