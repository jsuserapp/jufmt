package jufmt_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/jsuserapp/jufmt"
)

func setupBench(b *testing.B, printTrace, printTime bool) {
	b.Helper()
	jufmt.Output = io.Discard
	jufmt.SetLogOutputHook(nil)
	jufmt.SetPrintTrace(printTrace)
	jufmt.SetPrintTime(printTime)
}

// --- fmt (baseline) ---

func BenchmarkFmt_Fprintln(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		fmt.Fprintln(io.Discard, "hello", 42)
	}
}

func BenchmarkFmt_Fprintf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		fmt.Fprintf(io.Discard, "hello %d\n", 42)
	}
}

func BenchmarkFmt_Fprint(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		fmt.Fprint(io.Discard, "hello", 42)
	}
}

// --- jufmt Println / Printf / Print (no trace, no time) ---

func BenchmarkJufmt_Println(b *testing.B) {
	setupBench(b, false, false)
	b.ReportAllocs()
	for b.Loop() {
		jufmt.Println("hello", 42)
	}
}

func BenchmarkJufmt_Printf(b *testing.B) {
	setupBench(b, false, false)
	b.ReportAllocs()
	for b.Loop() {
		jufmt.Printf("hello %d\n", 42)
	}
}

func BenchmarkJufmt_Print(b *testing.B) {
	setupBench(b, false, false)
	b.ReportAllocs()
	for b.Loop() {
		jufmt.Print("hello", 42)
	}
}

// --- jufmt with optional prefixes (SetPrintTrace / SetPrintTime) ---

func BenchmarkJufmt_Println_Trace(b *testing.B) {
	setupBench(b, true, false)
	b.ReportAllocs()
	for b.Loop() {
		jufmt.Println("hello", 42)
	}
}

func BenchmarkJufmt_Println_TimeAndTrace(b *testing.B) {
	setupBench(b, true, true)
	b.ReportAllocs()
	for b.Loop() {
		jufmt.Println("hello", 42)
	}
}

// --- TracePrintln (always collects call-site; exStep adds upstream lines) ---

func BenchmarkJufmt_TracePrintln(b *testing.B) {
	for _, exStep := range []int{0, 1, 2} {
		b.Run(fmt.Sprintf("exStep%d", exStep), func(b *testing.B) {
			setupBench(b, false, false)
			b.ReportAllocs()
			for b.Loop() {
				benchTracePrintlnDepth(b, exStep)
			}
		})
	}
}

func benchTracePrintlnDepth(b *testing.B, exStep int) {
	b.Helper()
	if exStep <= 0 {
		jufmt.TracePrintln(0, "msg")
		return
	}
	benchTracePrintlnAtDepth(exStep, exStep)
}

func benchTracePrintlnAtDepth(exStep, depth int) {
	if depth > 0 {
		benchTracePrintlnAtDepth(exStep, depth-1)
		return
	}
	jufmt.TracePrintln(exStep, "msg")
}
