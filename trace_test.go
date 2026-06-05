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
	lines := strings.Split(strings.TrimSpace(plain), "\n")
	if len(lines) != 1 {
		t.Fatalf("Logger.Info with exStep=0 should print one line, got %q", plain)
	}
}

func TestTracePrintlnCallOrder(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)

	var buf bytes.Buffer
	old := jufmt.Output
	jufmt.Output = &buf
	t.Cleanup(func() { jufmt.Output = old })

	level1(t)
	plain := stripANSI(buf.String())

	idxTest := strings.Index(plain, "TestTracePrintlnCallOrder\n")
	idxLevel1 := strings.Index(plain, "level1\n")
	mainMsg := strings.Index(plain, "trace leaf")
	if idxTest < 0 || idxLevel1 < 0 || mainMsg < 0 {
		t.Fatalf("missing expected lines in %q", plain)
	}
	if !(idxTest < idxLevel1 && idxLevel1 < mainMsg) {
		t.Fatalf("expected upstream funcs before main message, got %q", plain)
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

	upstreamLines := func(s string) int {
		lines := strings.Split(strings.TrimSpace(stripANSI(s)), "\n")
		if len(lines) <= 1 {
			return 0
		}
		return len(lines) - 1
	}

	if upstreamLines(printlnOut) == 0 || upstreamLines(printfOut) == 0 {
		t.Fatalf("expected upstream call lines")
	}
	if upstreamLines(printlnOut) != upstreamLines(printfOut) {
		t.Fatalf("TracePrintln and TracePrintf differ in upstream line count")
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
