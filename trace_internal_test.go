package jufmt

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsTraceableFrame(t *testing.T) {
	userMain := filepath.Join("D:", "project", "main.go")
	runtimeMain := filepath.Join(runtime.GOROOT(), "src", "runtime", "proc.go")

	tests := []struct {
		name  string
		frame runtime.Frame
		want  bool
	}{
		{
			name:  "jufmt internal",
			frame: runtime.Frame{Function: pkgPath + ".Println", File: "fmt.go", Line: 10},
			want:  false,
		},
		{
			name:  "runtime.main",
			frame: runtime.Frame{Function: "runtime.main", File: runtimeMain, Line: 290},
			want:  false,
		},
		{
			name:  "user main",
			frame: runtime.Frame{Function: "main.main", File: userMain, Line: 23},
			want:  true,
		},
		{
			name:  "user test",
			frame: runtime.Frame{Function: "github.com/foo/bar_test.TestFoo", File: "bar_test.go", Line: 12},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTraceableFrame(tt.frame); got != tt.want {
				t.Fatalf("isTraceableFrame() = %v, want %v", got, tt.want)
			}
		})
	}
}
