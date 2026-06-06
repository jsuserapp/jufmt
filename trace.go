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

func getTraceableFrameAt(depth int) (runtime.Frame, bool) {
	return frameAt(traceableFrames(), depth)
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

// traceableFrames collects user-relevant stack frames, skipping jufmt, runtime, and stdlib.
func traceableFrames() []runtime.Frame {
	pcs := make([]uintptr, 64)
	// Skip runtime.Callers and traceableFrames itself.
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	var out []runtime.Frame
	for {
		frame, more := frames.Next()
		if isTraceableFrame(frame) {
			out = append(out, frame)
		}
		if !more {
			break
		}
	}
	return out
}
