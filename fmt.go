package jufmt

import gofmt "fmt"

// Printf writes formatted text to Output in bright white without a trailing newline.
// Timestamps and call-site prefixes follow SetPrintTime and SetPrintTrace.
func Printf(format string, a ...any) {
	gofmt.Fprint(Output, buildPrefix(0, false)+BrightWhite.Sprintf(format, a...))
}

// Println writes text to Output in bright white with a trailing newline.
// Timestamps and call-site prefixes follow SetPrintTime and SetPrintTrace.
func Println(a ...any) {
	gofmt.Fprint(Output, buildPrefix(0, false)+BrightWhite.Sprintln(a...))
}

// Print writes text to Output in bright white without a trailing newline.
// Timestamps and call-site prefixes follow SetPrintTime and SetPrintTrace.
func Print(a ...any) {
	gofmt.Fprint(Output, buildPrefix(0, false)+BrightWhite.Sprint(a...))
}

// Sprintf formats text without writing to Output or adding color prefixes.
func Sprintf(format string, a ...any) string {
	return gofmt.Sprintf(format, a...)
}

// Errorf builds an error with fmt.Errorf; it does not write to Output.
func Errorf(format string, a ...any) error {
	return gofmt.Errorf(format, a...)
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
