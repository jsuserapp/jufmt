package jufmt

import (
	"runtime"
	"testing"
)

// stackBenchPCCap matches traceableFrames: exStep(1) + slack for runtime/stdlib frames.
const stackBenchPCCap = 12

// stackBenchPrintlnPath mirrors Print/Println → emitMain → emitMainWithSkip → traceableFrames
// so runtime.Callers skip matches callersSkipEmitMain in production code.
func stackBenchPrintlnPath(walk func(int)) {
	stackBenchEmitMain(walk)
}

func stackBenchEmitMain(walk func(int)) {
	stackBenchEmitMainWithSkip(callersSkipEmitMain, walk)
}

func stackBenchEmitMainWithSkip(callersSkip int, walk func(int)) {
	walk(callersSkip)
}

func stackWalkCallersOnly(callersSkip int) {
	pc := make([]uintptr, stackBenchPCCap)
	runtime.Callers(callersSkip, pc)
}

func stackWalkCallersFrames(callersSkip int, resolveAll bool) {
	pc := make([]uintptr, stackBenchPCCap)
	n := runtime.Callers(callersSkip, pc)
	if n == 0 {
		return
	}
	frames := runtime.CallersFrames(pc[:n])
	if resolveAll {
		for {
			if _, more := frames.Next(); !more {
				break
			}
		}
		return
	}
	frames.Next()
}

func stackWalkTraceableFrames(callersSkip int) {
	traceableFrames(callersSkip, 1)
}

// --- stack micro-benchmarks (package jufmt, same skip as Println+trace) ---

func BenchmarkStack_Empty(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
	}
}

func BenchmarkStack_CallersOnly(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		stackBenchPrintlnPath(stackWalkCallersOnly)
	}
}

func BenchmarkStack_CallersFrames(b *testing.B) {
	for _, name := range []struct {
		label      string
		resolveAll bool
	}{
		{"firstFrame", false},
		{"allFrames", true},
	} {
		b.Run(name.label, func(b *testing.B) {
			resolveAll := name.resolveAll
			b.ReportAllocs()
			for b.Loop() {
				stackBenchPrintlnPath(func(skip int) {
					stackWalkCallersFrames(skip, resolveAll)
				})
			}
		})
	}
}

func BenchmarkStack_TraceableFrames(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		stackBenchPrintlnPath(stackWalkTraceableFrames)
	}
}
