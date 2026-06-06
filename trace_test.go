package jufmt_test

import (
	"bytes"
	"os"
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

func withHook(t *testing.T, h jufmt.LogOutputHook) {
	t.Helper()
	old := jufmt.Output
	jufmt.SetLogOutputHook(h)
	t.Cleanup(func() {
		jufmt.SetLogOutputHook(nil)
		jufmt.Output = old
	})
}

func stripANSI(s string) string {
	return regexp.MustCompile("\033\\[[0-9;]*m").ReplaceAllString(s, "")
}

func TestPrintlnCallerIsTestSite(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(true)

	var buf bytes.Buffer
	jufmt.Output = &buf
	t.Cleanup(func() { jufmt.Output = os.Stdout })

	jufmt.Println("msg")

	plain := stripANSI(buf.String())
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

	var buf bytes.Buffer
	jufmt.Output = &buf
	t.Cleanup(func() { jufmt.Output = os.Stdout })

	(&jufmt.Logger{}).Info("msg")

	plain := stripANSI(buf.String())
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

	idxTest := strings.Index(plain, "[TestTracePrintlnCallOrder]\n")
	idxLevel1 := strings.Index(plain, "[level1]\n")
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

func markerLevel1(t *testing.T) {
	t.Helper()
	markerLevel2(t)
}

func markerLevel1Printf(t *testing.T) {
	t.Helper()
	markerLevel2Printf(t)
}

func markerLevel2(t *testing.T) {
	t.Helper()
	jufmt.Cyan.TracePrintln(2, "marker")
}

func markerLevel2Printf(t *testing.T) {
	t.Helper()
	jufmt.Cyan.TracePrintf(2, "marker")
}

func TestTracePrintlnMatchesTracePrintf(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)

	var printlnBuf, printfBuf bytes.Buffer
	old := jufmt.Output
	t.Cleanup(func() { jufmt.Output = old })

	jufmt.Output = &printlnBuf
	markerLevel1(t)

	jufmt.Output = &printfBuf
	markerLevel1Printf(t)

	printlnOut := printlnBuf.String()
	printfOut := printfBuf.String()

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

	var depth0Buf bytes.Buffer
	jufmt.Output = &depth0Buf
	jufmt.Println(jufmt.GetTrace(0))

	var depth1Buf bytes.Buffer
	jufmt.Output = &depth1Buf
	captureDepth1(t)

	d0 := strings.TrimSpace(stripANSI(depth0Buf.String()))
	d1 := strings.TrimSpace(stripANSI(depth1Buf.String()))
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

func TestLogOutputHookDBOnly(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(true)

	var got jufmt.LogEntry
	var buf bytes.Buffer
	jufmt.Output = &buf
	withHook(t, func(entry jufmt.LogEntry) *jufmt.LogHookResult {
		got = entry
		return &jufmt.LogHookResult{WriteOutput: false}
	})

	jufmt.Println("persist")

	if buf.Len() != 0 {
		t.Fatalf("expected no Output, got %q", buf.String())
	}
	if got.Message != "persist" {
		t.Fatalf("unexpected message %q", got.Message)
	}
	if got.File == "" || !strings.HasSuffix(got.File, "trace_test.go") {
		t.Fatalf("expected runtime path ending in trace_test.go, got %q", got.File)
	}
	if strings.Contains(got.Message, "\033") {
		t.Fatal("hook message should be plain text without ANSI")
	}
}

func TestNoStackWhenPrintTraceDisabled(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)

	var got jufmt.LogEntry
	withHook(t, func(entry jufmt.LogEntry) *jufmt.LogHookResult {
		got = entry
		return &jufmt.LogHookResult{WriteOutput: false}
	})

	jufmt.Println("no trace")

	if got.ShowLocation {
		t.Fatalf("ShowLocation should be false, got %+v", got)
	}
	if got.File != "" || got.Line != 0 || got.Func != "" {
		t.Fatalf("location fields should be empty, got %+v", got)
	}
}

func TestLogOutputHookCustomOutput(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)

	var buf bytes.Buffer
	jufmt.Output = &buf
	withHook(t, func(entry jufmt.LogEntry) *jufmt.LogHookResult {
		return &jufmt.LogHookResult{WriteOutput: true, Output: "CUSTOM:" + entry.Message + "\n"}
	})

	jufmt.Println("x")

	if buf.String() != "CUSTOM:x\n" {
		t.Fatalf("unexpected output %q", buf.String())
	}
}

func TestLogOutputHookNilUsesDefault(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(true)

	var buf bytes.Buffer
	jufmt.Output = &buf
	withHook(t, func(entry jufmt.LogEntry) *jufmt.LogHookResult {
		_ = entry
		return nil
	})

	jufmt.Println("default")

	want := stripANSI(buf.String())
	if !strings.Contains(want, "trace_test.go:") || !strings.Contains(want, "default") {
		t.Fatalf("nil hook result should use DefaultFormat, got %q", want)
	}
}

func TestLogOutputHookDefaultFormat(t *testing.T) {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(true)

	var buf bytes.Buffer
	jufmt.Output = &buf
	withHook(t, func(entry jufmt.LogEntry) *jufmt.LogHookResult {
		return &jufmt.LogHookResult{WriteOutput: true}
	})

	jufmt.Println("y")

	want := stripANSI(buf.String())
	if !strings.Contains(want, "trace_test.go:") || !strings.Contains(want, "y") {
		t.Fatalf("expected default format, got %q", want)
	}
}
