package jufmt

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const pkgPath = "github.com/jsuserapp/jufmt"

// Output is the writer used by Print/Println/Printf and Color print methods.
// Tests and integrations may redirect it to capture output.
var Output io.Writer = os.Stdout

var (
	printTime  atomic.Bool
	printTrace atomic.Bool
)

func init() {
	printTime.Store(true)
	printTrace.Store(true)
}

// SetPrintTime enables or disables timestamps in printed output.
func SetPrintTime(enabled bool) {
	printTime.Store(enabled)
}

// SetPrintTrace enables or disables call-site prefixes in Print/Println/Printf.
// TracePrintln and TracePrintf always include call-site information.
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
	frames := traceableFrames()
	if depth < 0 || depth >= len(frames) {
		return ""
	}
	return formatFrame(frames[depth])
}

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("%s:%d", path.Base(frame.File), frame.Line)
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

// formatCallPrefix builds the time + location prefix for a [call N] line at trace depth.
// Returns empty when that depth has no traceable frame (for example runtime above main).
func formatCallPrefix(depth int) string {
	loc := GetTrace(depth)
	if loc == "" {
		return ""
	}
	var prefix string
	if printTime.Load() {
		prefix += GetNowTimeMs() + " "
	}
	return prefix + loc + " "
}

// buildPrefix assembles the time and call-site prefix for a print call.
// forceTrace includes call-site even when the global printTrace flag is off.
// traceDepth selects which external frame to show (0 = direct caller).
func buildPrefix(traceDepth int, forceTrace bool) string {
	var prefix string
	if printTime.Load() {
		prefix += GetNowTimeMs() + " "
	}
	if forceTrace || printTrace.Load() {
		loc := GetTrace(traceDepth)
		if loc == "" && traceDepth == 0 {
			loc = "unknown:0"
		}
		if loc != "" {
			prefix += loc + " "
		}
	}
	return prefix
}
