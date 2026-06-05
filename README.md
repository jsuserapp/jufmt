# jufmt

A partial replacement for the standard `fmt` and `log` packages. Most APIs are drop-in compatible.

Prints optional timestamps and call-site locations (`file:line`) alongside messages so you can jump to source in the IDE.

[中文文档](README.zh.md)

## Call-site tracing

- `Print`, `Println`, `Printf`, and the same methods on `Color`: call-site prefixes are controlled by `SetPrintTrace` and show the **direct caller** as `file:line`.
- `TracePrintln` / `TracePrintf`: always include call-site on the main line; `exStep` adds upstream lines before it. Each upstream line is `file:line` plus the **short function name** at that frame. Lines are printed in execution order (earliest call first).
- Locations use `runtime.Callers`, skipping jufmt internals, Go runtime (e.g. `proc.go`), and the standard library (GOROOT).
- `GetTrace(depth)` returns `file:line`; `GetCallerName(depth)` returns the short function name at the same depth.

Example (`TracePrintln(2, "traceTest2")` inside `traceTest2`, called from `traceTest1`, called from `main`):

```
main.go:19 main
main.go:27 traceTest1
main.go:30 traceTest2
```

## Redirecting output

`Output` is an `io.Writer` (default `os.Stdout`). Assign any writer to send jufmt output elsewhere.

**Log file**

```go
f, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
if err != nil { /* handle */ }
defer f.Close()
jufmt.Output = f
```

**Terminal and file together**

```go
jufmt.Output = io.MultiWriter(os.Stdout, f)
```

**Tests**

```go
var buf bytes.Buffer
jufmt.Output = &buf
jufmt.Println("capture me")
_ = buf.String()
```

**Database or remote log service**

Implement `io.Writer` and insert each `Write` payload as a log record:

```go
type dbLogWriter struct{ db *sql.DB }

func (w *dbLogWriter) Write(p []byte) (int, error) {
	_, err := w.db.Exec("INSERT INTO logs(body) VALUES (?)", string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// jufmt.Output = &dbLogWriter{db: conn}
```

Notes:

- jufmt writes **plain text** (with ANSI codes when using `Color`). A custom writer receives finished byte chunks, not structured fields.
- For databases, prefer one line per `Println`/`TracePrintln` call; avoid relying on partial `Print` fragments.
- Strip ANSI codes before storing if terminals are not the consumer (`strings` or a small strip helper).
- `SetPrintTime(false)` and bright colors are often disabled when logging to files or DB.

## Other

- `Log`: package-level default `Logger` instance.

Because each print adds a prefix, output is usually one line at a time. Unlike plain `fmt`, chaining partial-line writes without newlines tends to look messy when prefixes are concatenated.

## Example

```go
jufmt.SetPrintTime(true)
jufmt.SetPrintTrace(true)

jufmt.Println("hello")
jufmt.Green.Println("success")
jufmt.Log.Info("ready")

jufmt.Green.TracePrintln(2, "nested trace")
```
