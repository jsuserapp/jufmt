package jufmt

import (
	"fmt"
	"strings"
)

// Color wraps an ANSI escape sequence for terminal text styling.
type Color struct {
	code string
}

// Printf writes formatted, colored text to Output without a trailing newline.
func (c Color) Printf(format string, a ...interface{}) {
	emitMain(c, false, fmt.Sprintf(format, a...), false)
}

// Println writes colored text to Output with a trailing newline.
func (c Color) Println(a ...interface{}) {
	msg := fmt.Sprintln(a...)
	emitMain(c, false, trimNewline(msg), true)
}

// Print writes colored text to Output without a trailing newline.
func (c Color) Print(a ...interface{}) {
	emitMain(c, false, fmt.Sprint(a...), false)
}

func (c Color) tracePrintln(forceTrace bool, a ...interface{}) {
	msg := fmt.Sprintln(a...)
	emitMain(c, forceTrace, trimNewline(msg), true)
}

func (c Color) tracePrintf(forceTrace bool, format string, a ...interface{}) {
	emitMain(c, forceTrace, fmt.Sprintf(format, a...), false)
}

// TracePrintln writes a message with optional upstream call-site lines.
//
// exStep is the number of extra caller frames to print before the main message.
// Each upstream line shows file:line and the short function name at that frame.
// Lines are printed in execution order: the furthest upstream frame first,
// then nearer frames, then the main message. Negative exStep is treated as 0.
// Call-site prefixes are always included, regardless of SetPrintTrace.
// Upstream frames inside Go runtime or the standard library are omitted.
func (c Color) TracePrintln(exStep int, a ...any) {
	if exStep < 0 {
		exStep = 0
	}
	frames := traceableFrames()
	for i := exStep; i >= 1; i-- {
		if entry, ok := buildUpstreamEntry(frames, i); ok {
			entry.Color = c
			emit(entry)
		}
	}
	msg := fmt.Sprintln(a...)
	emitMainFromFrames(c, frames, true, trimNewline(msg), true)
}

// TracePrintf is like TracePrintln but with fmt.Printf-style formatting.
func (c Color) TracePrintf(exStep int, format string, a ...any) {
	if exStep < 0 {
		exStep = 0
	}
	frames := traceableFrames()
	for i := exStep; i >= 1; i-- {
		if entry, ok := buildUpstreamEntry(frames, i); ok {
			entry.Color = c
			emit(entry)
		}
	}
	emitMainFromFrames(c, frames, true, fmt.Sprintf(format, a...), false)
}

func trimNewline(s string) string {
	return strings.TrimSuffix(s, "\n")
}

// Sprintf returns a formatted string wrapped with ANSI color codes.
// The result can be passed to standard fmt print functions for multi-color output.
func (c Color) Sprintf(format string, a ...interface{}) string {
	return c.code + fmt.Sprintf(format, a...) + reset
}

// Sprintln returns a line-oriented string wrapped with ANSI color codes.
func (c Color) Sprintln(a ...interface{}) string {
	return c.code + fmt.Sprintln(a...) + reset
}

// Sprint returns a string wrapped with ANSI color codes.
func (c Color) Sprint(a ...interface{}) string {
	return c.code + fmt.Sprint(a...) + reset
}

const (
	reset = "\033[0m"
)

// Standard ANSI foreground colors.
var (
	Black   = Color{"\033[30m"}
	Red     = Color{"\033[31m"}
	Green   = Color{"\033[32m"}
	Yellow  = Color{"\033[33m"}
	Blue    = Color{"\033[34m"}
	Magenta = Color{"\033[35m"}
	Cyan    = Color{"\033[36m"}
	White   = Color{"\033[37m"}

	BrightBlack   = Color{"\033[90m"}
	BrightRed     = Color{"\033[91m"}
	BrightGreen   = Color{"\033[92m"}
	BrightYellow  = Color{"\033[93m"}
	BrightBlue    = Color{"\033[94m"}
	BrightMagenta = Color{"\033[95m"}
	BrightCyan    = Color{"\033[96m"}
	BrightWhite   = Color{"\033[97m"}

	Gray = BrightBlack
)