package jufmt

import gofmt "fmt"

func Printf(format string, a ...any) {
	BrightWhite.Outputf(1, format, a...)
}
func Println(a ...any) {
	BrightWhite.Outputln(1, a...)
}
func Sprintf(format string, a ...any) string {
	return gofmt.Sprintf(format, a...)
}
