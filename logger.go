package jufmt

import "os"

type Logger struct{}

func (log *Logger) Debug(a ...any) {
	Cyan.TracePrintln(1, a...)
}
func (log *Logger) Info(a ...any) {
	Green.TracePrintln(1, a...)
}
func (log *Logger) Warn(a ...any) {
	Yellow.TracePrintln(1, a...)
}
func (log *Logger) Error(a ...any) {
	Red.TracePrintln(1, a...)
}
func (log *Logger) With(_ any) *Logger {
	return log
}
func (log *Logger) Fatal(a ...any) {
	Red.TracePrintln(1, a...)
	os.Exit(1)
}
