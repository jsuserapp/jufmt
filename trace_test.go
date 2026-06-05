package jufmt_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/jsuserapp/jufmt"
)

func withOutput(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	old := jufmt.Output
	jufmt.Output = &buf
	t.Cleanup(func() { jufmt.Output = old })
	fn()
	return buf.String()
}

func stripANSI(s string) string {
	return regexp.MustCompile("\033\\[[0-9;]*m").ReplaceAllString(s, "")
}

func TestPrintlnCallerIsTestSite(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(true)

	out := withOutput(t, func() {
		jufmt.Println("msg")
	})
	plain := stripANSI(out)
	if strings.Contains(plain, "fmt.go:") {
		t.Fatalf("Println should not report internal fmt.go, got %q", plain)
	}
	if !strings.Contains(plain, "trace_test.go:") {
		t.Fatalf("Println should report test caller, got %q", plain)
	}
}

func TestLoggerInfoCallerIsNotLogger(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(true)

	out := withOutput(t, func() {
		(&jufmt.Logger{}).Info("msg")
	})
	plain := stripANSI(out)
	if strings.Contains(plain, "logger.go:") {
		t.Fatalf("Logger.Info should not report logger.go, got %q", plain)
	}
	if !strings.Contains(plain, "trace_test.go:") {
		t.Fatalf("Logger.Info should report test caller, got %q", plain)
	}
	if strings.Contains(plain, "[call") {
		t.Fatalf("Logger.Info with exStep=0 should not print [call N] lines, got %q", plain)
	}
}

func TestTracePrintlnCallOrder(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)

	out := withOutput(t, func() {
		level1(t)
	})
	plain := stripANSI(out)

	call1 := strings.Index(plain, "[call 1]")
	call2 := strings.Index(plain, "[call 2]")
	mainMsg := strings.Index(plain, "trace leaf")
	if call1 < 0 || call2 < 0 || mainMsg < 0 {
		t.Fatalf("missing expected lines in %q", plain)
	}
	if !(call2 < call1 && call1 < mainMsg) {
		t.Fatalf("expected [call 2] before [call 1] before main message, got %q", plain)
	}
	if strings.Contains(plain, "proc.go:") {
		t.Fatalf("should not trace into runtime, got %q", plain)
	}
}

func level1(t *testing.T) {
	t.Helper()
	level2(t)
}

func level2(t *testing.T) {
	t.Helper()
	jufmt.Green.TracePrintln(2, "trace leaf")
}

func TestTracePrintlnMatchesTracePrintf(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)

	printlnOut := withOutput(t, func() {
		jufmt.Cyan.TracePrintln(2, "marker")
	})
	printfOut := withOutput(t, func() {
		jufmt.Cyan.TracePrintf(2, "marker")
	})

	countCalls := func(s string) int {
		return strings.Count(stripANSI(s), "[call ")
	}

	if countCalls(printlnOut) != 2 || countCalls(printfOut) != 2 {
		t.Fatalf("expected 2 call lines each, got println=%d printf=%d", countCalls(printlnOut), countCalls(printfOut))
	}
}

func TestGetTraceDepth(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)

	depth0 := withOutput(t, func() {
		jufmt.Println(jufmt.GetTrace(0))
	})
	depth1 := withOutput(t, func() {
		captureDepth1(t)
	})

	d0 := strings.TrimSpace(stripANSI(depth0))
	d1 := strings.TrimSpace(stripANSI(depth1))
	if d0 == d1 {
		t.Fatalf("GetTrace(0) and GetTrace(1) should differ, both got %q", d0)
	}
	if !strings.Contains(d0, "trace_test.go:") {
		t.Fatalf("GetTrace(0) unexpected: %q", d0)
	}
}

func captureDepth1(t *testing.T) {
	t.Helper()
	jufmt.Println(jufmt.GetTrace(1))
}
