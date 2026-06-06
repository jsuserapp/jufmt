package jufmt

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const pkgPath = "github.com/jsuserapp/jufmt"

// Output is the destination for Print/Println/Printf and Color print methods.
// It defaults to os.Stdout. Use SetLogOutputHook for structured logging or to
// suppress console output; assign Output to redirect formatted text elsewhere.
var Output io.Writer = os.Stdout

var (
	printTime  atomic.Bool
	printTrace atomic.Bool
)

func init() {
	printTime.Store(true)
	printTrace.Store(true)
}

// SetPrintTime enables or disables timestamp prefixes in printed output.
// When disabled, TimeText stays empty and timestamp formatting is skipped.
// LogEntry.Time is still set for hooks that need a wall-clock value.
func SetPrintTime(enabled bool) {
	printTime.Store(enabled)
}

// SetPrintTrace enables or disables call-site collection and prefixes on Print/Println/Printf.
// When disabled, stack frames are not collected for those APIs (same as SetPrintTime skipping
// timestamp formatting). TracePrintln, TracePrintf, and Logger always collect call-site data.
func SetPrintTrace(enabled bool) {
	printTrace.Store(enabled)
}

// GetNowTimeMs returns the current local time formatted as HH:MM:SS.mmm.
func GetNowTimeMs() string {
	return time.Now().Format("15:04:05.000")
}

// GetTrace returns file:line for a traceable caller frame.
//
// depth 0 is the direct call site (first frame outside jufmt, runtime, and Go stdlib).
// depth 1 is one user-code frame further up the call chain, and so on.
//
// Returns an empty string when depth is out of range or the frame is not traceable
// (for example runtime.main in proc.go). Unlike runtime.Caller(skip), depth counts
// only user-relevant frames, so it stays stable across refactors and inlining.
func GetTrace(depth int) string {
	frame, ok := getTraceableFrameAt(depth)
	if !ok {
		return ""
	}
	return formatFrame(frame)
}

// GetCallerName returns the short function name at trace depth (same semantics as GetTrace).
func GetCallerName(depth int) string {
	frame, ok := getTraceableFrameAt(depth)
	if !ok {
		return ""
	}
	return frameFuncName(frame)
}

func formatFrame(frame runtime.Frame) string {
	return formatFrameLocation(frame.File, frame.Line)
}

func frameAt(frames []runtime.Frame, depth int) (runtime.Frame, bool) {
	if depth < 0 || depth >= len(frames) {
		return runtime.Frame{}, false
	}
	return frames[depth], true
}

func frameFuncName(frame runtime.Frame) string {
	fn := frame.Function
	if fn == "" {
		return "unknown"
	}
	if i := strings.LastIndex(fn, "/"); i >= 0 {
		fn = fn[i+1:]
	}
	if i := strings.LastIndex(fn, "."); i >= 0 {
		fn = fn[i+1:]
	}
	return fn
}

// callersSkip constants are passed to runtime.Callers from traceableFrames.
// skip 0 is Callers itself; each value counts frames through traceableFrames up to user pc[0].
const (
	// traceableFrames → emitMainWithSkip → emitMain → Color.Print* / package Print* → user
	callersSkipEmitMain = 5
	// traceableFrames → emitMainWithSkip → tracePrintln → Logger.* → user
	callersSkipTracePrintln = 5
	// traceableFrames → Color.TracePrint* → user
	callersSkipTracePrint = 3
	// traceableFrames → getTraceableFrameAt → GetTrace/GetCallerName → user
	callersSkipGetTrace = 4
)

func getTraceableFrameAt(depth int) (runtime.Frame, bool) {
	if depth < 0 {
		return runtime.Frame{}, false
	}
	frames := traceableFrames(callersSkipGetTrace, depth+1)
	return frameAt(frames, depth)
}

func isInternalFrame(frame runtime.Frame) bool {
	return strings.HasPrefix(frame.Function, pkgPath+".")
}

var (
	gorootPath     string
	gorootPathOnce sync.Once
)

func goRoot() string {
	gorootPathOnce.Do(func() {
		gorootPath = filepath.Clean(runtime.GOROOT())
	})
	return gorootPath
}

// isTraceableFrame reports whether a frame should appear in trace output.
// jufmt internals, runtime, and Go standard library (GOROOT) are excluded.
func isTraceableFrame(frame runtime.Frame) bool {
	if isInternalFrame(frame) {
		return false
	}
	if strings.HasPrefix(frame.Function, "runtime.") {
		return false
	}
	if root := goRoot(); root != "" {
		file := filepath.Clean(frame.File)
		if file == root || strings.HasPrefix(file, root+string(os.PathSeparator)) {
			return false
		}
	}
	return true
}

// traceableFrames collects up to exStep stack frames starting at the user call site.
//
// callersSkip is the runtime.Callers skip counted at the call site that invokes this
// function (including traceableFrames itself). Each API passes its own skip so library
// wrappers do not hard-code frame depth here.
//
// pc[0] is the direct user call site; pc[1] is one caller further up, and so on.
func traceableFrames(callersSkip, exStep int) []runtime.Frame {
	if exStep < 1 {
		exStep = 1
	}
	pc := make([]uintptr, exStep)
	n := runtime.Callers(callersSkip, pc)
	if n == 0 {
		return nil
	}
	callerFrames := runtime.CallersFrames(pc[:n])
	out := make([]runtime.Frame, 0, n)
	for {
		frame, more := callerFrames.Next()
		out = append(out, frame)
		if !more {
			break
		}
	}
	return out
}
