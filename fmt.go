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
