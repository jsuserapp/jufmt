package jufmt

import gofmt "fmt"

func Printf(format string, a ...any) {
	BrightWhite.printf(1, format, a...)
}
func Println(a ...any) {
	BrightWhite.println(1, a...)
}
func Sprintf(format string, a ...any) string {
	return gofmt.Sprintf(format, a...)
}
