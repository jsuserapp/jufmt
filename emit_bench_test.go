package jufmt

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"
)

func benchSampleEntry() LogEntry {
	return LogEntry{
		Message:      "hello 42",
		Color:        BrightWhite,
		Newline:      true,
		ShowLocation: true,
		File:         "stack_bench_test.go",
		Line:         1,
		Func:         "benchSample",
	}
}

// --- emit path micro-benchmarks ---

func BenchmarkEmit_Sprintln(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = strings.TrimSuffix(fmt.Sprintln("hello", 42), "\n")
	}
}

func BenchmarkEmit_frameFuncName(b *testing.B) {
	frame := runtime.Frame{Function: "github.com/jsuserapp/jufmt_test.BenchmarkJufmt_Println_Trace"}
	b.ReportAllocs()
	for b.Loop() {
		_ = frameFuncName(frame)
	}
}

func BenchmarkEmit_DefaultFormat(b *testing.B) {
	entry := benchSampleEntry()
	b.ReportAllocs()
	for b.Loop() {
		_ = DefaultFormat(entry)
	}
}

func BenchmarkEmit_FprintDiscard(b *testing.B) {
	s := DefaultFormat(benchSampleEntry())
	b.ReportAllocs()
	for b.Loop() {
		_, _ = fmt.Fprint(io.Discard, s)
	}
}

func BenchmarkEmit_buildMainEntry(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		stackBenchPrintlnPath(func(skip int) {
			frames := traceableFrames(skip, 1)
			_ = buildMainEntry(BrightWhite, frames, false, "hello 42", true)
		})
	}
}

func BenchmarkEmit_emitOnly(b *testing.B) {
	entry := benchSampleEntry()
	old := Output
	Output = io.Discard
	b.Cleanup(func() { Output = old })
	b.ReportAllocs()
	for b.Loop() {
		emit(entry)
	}
}

func BenchmarkEmit_buildMainEntryPlusEmit(b *testing.B) {
	old := Output
	Output = io.Discard
	b.Cleanup(func() { Output = old })
	b.ReportAllocs()
	for b.Loop() {
		stackBenchPrintlnPath(func(skip int) {
			frames := traceableFrames(skip, 1)
			emit(buildMainEntry(BrightWhite, frames, false, "hello 42", true))
		})
	}
}
