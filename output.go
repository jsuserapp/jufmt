package jufmt

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// LineKind identifies which part of a TracePrintln/TracePrintf emission a LogEntry represents.
type LineKind int

const (
	// LineMain is the user-facing message line: Print/Println/Printf output, the final
	// line of TracePrintln/TracePrintf, or Logger level output. Message holds user content.
	LineMain LineKind = iota
	// LineUpstream is an extra call-chain line emitted before LineMain when exStep > 0.
	// Message holds the short function name at that upstream frame, wrapped in [brackets].
	LineUpstream
)

// LogEntry is structured data for one emitted log line passed to LogOutputHook.
//
// File is the source path from runtime.Frame (typically an absolute path on disk).
// Console output still uses the basename via DefaultFormat. Message is plain text without ANSI.
type LogEntry struct {
	Time         time.Time // wall time of the emit; set even when TimeText is empty
	TimeText     string    // formatted timestamp prefix; empty when SetPrintTime(false)
	File         string    // runtime source path; empty when location was not collected
	Line         int
	Func         string
	Message      string
	Color        Color
	Newline      bool     // append '\n' in DefaultFormat when true
	Kind         LineKind // LineMain vs LineUpstream (TracePrintln call-chain lines)
	ShowLocation bool     // false when SetPrintTrace(false) on non-trace APIs
}

// LogHookResult controls how a hook handles console Output.
//
// Hooks return *LogHookResult. Return nil to keep the default: write DefaultFormat(entry)
// to Output. Return a non-nil value to override (for example WriteOutput: false for
// database-only mode, or a custom Output string).
type LogHookResult struct {
	WriteOutput bool
	Output      string
}

// LogOutputHook is invoked for each line jufmt would print.
//
// The library does not perform database or network I/O. Use entry fields (runtime File path,
// Line, Func, plain Message, Kind, etc.) inside the hook for persistence.
//
// Return nil when you only need the callback and default console output is fine.
// Return &LogHookResult{WriteOutput: false} to skip Output entirely.
// Return &LogHookResult{WriteOutput: true, Output: "..."} for a custom console line;
// when WriteOutput is true and Output is empty, DefaultFormat(entry) is used.
//
// Keep the hook fast; defer heavy work (database inserts, RPC) with go func() if needed.
// Async hooks lose strict ordering relative to other goroutines—document and handle that
// in application code if ordering matters.
//
// Set h to nil to restore built-in console formatting without a hook.
type LogOutputHook func(entry LogEntry) *LogHookResult

var logOutputHook atomic.Value // stores LogOutputHook or nil

// SetLogOutputHook installs a hook called before writing to Output.
func SetLogOutputHook(h LogOutputHook) {
	if h == nil {
		logOutputHook.Store((LogOutputHook)(nil))
		return
	}
	logOutputHook.Store(h)
}

func currentLogOutputHook() LogOutputHook {
	v := logOutputHook.Load()
	if v == nil {
		return nil
	}
	return v.(LogOutputHook)
}

// DefaultFormat renders entry for terminal Output: basename file:line, optional time,
// ANSI color when entry.Color is set, and plain Message.
func DefaultFormat(entry LogEntry) string {
	var b strings.Builder
	if entry.TimeText != "" {
		b.WriteString(entry.TimeText)
		b.WriteByte(' ')
	}
	if entry.ShowLocation {
		if entry.File != "" {
			b.WriteString(formatFrameLocation(entry.File, entry.Line))
			b.WriteByte(' ')
		} else {
			b.WriteString("unknown:0 ")
		}
	}
	if entry.Color.code != "" {
		b.WriteString(entry.Color.code)
	}
	b.WriteString(entry.Message)
	if entry.Color.code != "" {
		b.WriteString(reset)
	}
	if entry.Newline {
		b.WriteByte('\n')
	}
	return b.String()
}

func formatFrameLocation(file string, line int) string {
	return path.Base(file) + ":" + strconv.Itoa(line)
}

func emit(entry LogEntry) {
	hook := currentLogOutputHook()
	if hook != nil {
		res := hook(entry)
		if res == nil {
			_, _ = fmt.Fprint(Output, DefaultFormat(entry))
			return
		}
		if !res.WriteOutput {
			return
		}
		out := res.Output
		if out == "" {
			out = DefaultFormat(entry)
		}
		_, _ = fmt.Fprint(Output, out)
		return
	}
	_, _ = fmt.Fprint(Output, DefaultFormat(entry))
}

func buildMainEntry(c Color, frames []runtime.Frame, forceTrace bool, plainMsg string, newline bool) LogEntry {
	entry := LogEntry{
		Time:    time.Now(),
		Message: plainMsg,
		Color:   c,
		Newline: newline,
		Kind:    LineMain,
	}
	if printTime.Load() {
		entry.TimeText = entry.Time.Format("15:04:05.000")
	}
	showLoc := forceTrace || printTrace.Load()
	entry.ShowLocation = showLoc
	if showLoc && len(frames) > 0 {
		frame := frames[0]
		if isTraceableFrame(frame) {
			entry.File = frame.File
			entry.Line = frame.Line
			entry.Func = frameFuncName(frame)
		}
	}
	return entry
}

func buildUpstreamEntryFromFrame(frame runtime.Frame) (LogEntry, bool) {
	if !isTraceableFrame(frame) {
		return LogEntry{}, false
	}
	entry := LogEntry{
		Time:         time.Now(),
		Message:      "[" + frameFuncName(frame) + "]",
		Newline:      true,
		Kind:         LineUpstream,
		ShowLocation: true,
		File:         frame.File,
		Line:         frame.Line,
		Func:         frameFuncName(frame),
	}
	if printTime.Load() {
		entry.TimeText = entry.Time.Format("15:04:05.000")
	}
	return entry, true
}

func emitMain(c Color, forceTrace bool, plainMsg string, newline bool) {
	emitMainWithSkip(c, callersSkipEmitMain, forceTrace, plainMsg, newline)
}

func emitMainWithSkip(c Color, callersSkip int, forceTrace bool, plainMsg string, newline bool) {
	var frames []runtime.Frame
	if forceTrace || printTrace.Load() {
		frames = traceableFrames(callersSkip, 1)
	}
	emit(buildMainEntry(c, frames, forceTrace, plainMsg, newline))
}

func emitMainFromFrames(c Color, frames []runtime.Frame, forceTrace bool, plainMsg string, newline bool) {
	emit(buildMainEntry(c, frames, forceTrace, plainMsg, newline))
}

func emitUpstreamFrames(c Color, frames []runtime.Frame) {
	for i := len(frames) - 1; i >= 1; i-- {
		if entry, ok := buildUpstreamEntryFromFrame(frames[i]); ok {
			entry.Color = c
			emit(entry)
		}
	}
}
