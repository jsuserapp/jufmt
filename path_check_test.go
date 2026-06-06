package jufmt_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jsuserapp/jufmt"
)

func TestLogEntryFileFromRuntime(t *testing.T) {
	var file string
	jufmt.SetLogOutputHook(func(e jufmt.LogEntry) *jufmt.LogHookResult {
		file = e.File
		return nil
	})
	t.Cleanup(func() { jufmt.SetLogOutputHook(nil) })

	jufmt.SetPrintTrace(true)
	jufmt.SetPrintTime(false)
	jufmt.Println("probe")

	if file == "" {
		t.Fatal("expected File in LogEntry")
	}
	if !filepath.IsAbs(file) {
		t.Fatalf("LogEntry.File should be runtime absolute path, got %q", file)
	}
	if !strings.HasSuffix(file, "path_check_test.go") {
		t.Fatalf("unexpected File %q", file)
	}
}
