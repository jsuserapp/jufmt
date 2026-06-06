package jufmt

import "os"

// Logger writes leveled, color-coded messages with call-site prefixes.
type Logger struct{}

// Log is the default package-level logger.
var Log = &Logger{}

// Debug writes a cyan message with the caller's file:line.
func (log *Logger) Debug(a ...any) {
	Cyan.tracePrintln(true, a...)
}

// Info writes a green message with the caller's file:line.
func (log *Logger) Info(a ...any) {
	Green.tracePrintln(true, a...)
}

// Warn writes a yellow message with the caller's file:line.
func (log *Logger) Warn(a ...any) {
	Yellow.tracePrintln(true, a...)
}

// Error writes a red message with the caller's file:line.
func (log *Logger) Error(a ...any) {
	Red.tracePrintln(true, a...)
}

// With is a no-op adapter for components that expect a chainable With method.
func (log *Logger) With(_ any) *Logger {
	return log
}

// Fatal writes a red message with the caller's file:line and exits with code 1.
func (log *Logger) Fatal(a ...any) {
	Red.tracePrintln(true, a...)
	os.Exit(1)
}
