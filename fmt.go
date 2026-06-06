package jufmt

import "fmt"

// Printf writes formatted text to Output in bright white without a trailing newline.
// Timestamps and call-site prefixes follow SetPrintTime and SetPrintTrace.
func Printf(format string, a ...any) {
	emitMain(BrightWhite, false, fmt.Sprintf(format, a...), false)
}

// Println writes text to Output in bright white with a trailing newline.
// Timestamps and call-site prefixes follow SetPrintTime and SetPrintTrace.
func Println(a ...any) {
	msg := fmt.Sprintln(a...)
	emitMain(BrightWhite, false, trimNewline(msg), true)
}

// Print writes text to Output in bright white without a trailing newline.
// Timestamps and call-site prefixes follow SetPrintTime and SetPrintTrace.
func Print(a ...any) {
	emitMain(BrightWhite, false, fmt.Sprint(a...), false)
}

// Sprintf formats text without writing to Output or adding color prefixes.
func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

// Errorf builds an error with fmt.Errorf; it does not write to Output.
func Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// TracePrintln writes a bright-white message with optional upstream call-site lines.
// See Color.TracePrintln for exStep semantics.
func TracePrintln(exStep int, a ...any) {
	BrightWhite.TracePrintln(exStep, a...)
}

// TracePrintf writes a bright-white formatted message with optional upstream call-site lines.
// See Color.TracePrintf for exStep semantics.
func TracePrintf(exStep int, format string, a ...any) {
	BrightWhite.TracePrintf(exStep, format, a...)
}
