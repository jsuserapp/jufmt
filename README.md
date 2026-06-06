# jufmt

A partial replacement for the standard `fmt` and `log` packages. Most APIs are drop-in compatible.

Prints optional timestamps and call-site locations (`file:line`) alongside messages so you can jump to source in the IDE.

[中文文档](README.zh.md)

## Call-site tracing

- `Print`, `Println`, `Printf`, and the same methods on `Color`: call-site prefixes are controlled by `SetPrintTrace` and show the **direct caller** as `file:line` (basename on console).
- `TracePrintln` / `TracePrintf`: always include call-site on the main line; `exStep` adds upstream lines before it. Each upstream line is `file:line` plus the **short function name** at that frame. Lines are printed in execution order (earliest call first).
- Locations use `runtime.Callers`, skipping jufmt internals, Go runtime (e.g. `proc.go`), and the standard library (GOROOT).
- `GetTrace(depth)` returns `file:line` (basename); `GetCallerName(depth)` returns the short function name at the same depth.
- `SetPrintTrace(false)` on `Print`/`Println`/`Printf`: **no stack walk**; location fields stay empty (same idea as `SetPrintTime(false)` skipping `TimeText`). `TracePrintln`, `TracePrintf`, and `Logger` always collect call-site data.

Example (`TracePrintln(2, "traceTest2")` inside `traceTest2`, called from `traceTest1`, called from `main`):

```
main.go:19 main
main.go:27 traceTest1
main.go:30 traceTest2
```

## LogOutputHook (structured logging)

`SetLogOutputHook` receives a `LogEntry` for every line jufmt would print. The library does **not** talk to databases; use `LogEntry` fields inside the hook:

| Field | Use |
|-------|-----|
| `File` | **Full path** when location was collected; empty if `SetPrintTrace(false)` on Print APIs |
| `Line`, `Func` | Call-site; empty when location not collected |
| `Message` | Plain text, **no ANSI**; user content on `LineMain`, function name on `LineUpstream` |
| `Time` | Wall-clock time (always set) |
| `TimeText` | Formatted prefix; empty when `SetPrintTime(false)` |
| `Newline` | `true` for Println-style (append `\n` in `DefaultFormat`); `false` for Print/Printf |
| `Kind` | `LineMain` = user message line; `LineUpstream` = extra call-chain line from `TracePrintln` |
| `ShowLocation` | Whether `File`/`Line`/`Func` were collected for this line |
| `Color` | ANSI hint for console formatting |

Return `*LogHookResult` (or **`nil`** for default console output):

| Return | Behavior |
|--------|----------|
| **`nil`** | Write `DefaultFormat(entry)` to `Output` — use when you only need the callback |
| `&LogHookResult{WriteOutput: false}` | Nothing written to `Output` (database-only server mode) |
| `&LogHookResult{WriteOutput: true}` | Same as `nil` (`DefaultFormat`) |
| `&LogHookResult{WriteOutput: true, Output: "..."}` | Custom string written to `Output` |

**Callback only (default console)**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) *jufmt.LogHookResult {
	go dbInsert(e) // e.File is full path, e.Message has no ANSI
	return nil
})
```

**Database only (no console)**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) *jufmt.LogHookResult {
	go dbInsert(e)
	return &jufmt.LogHookResult{WriteOutput: false}
})
```

**Custom console format**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) *jufmt.LogHookResult {
	persist(e)
	line := fmt.Sprintf("%s %s:%d %s\n", e.TimeText, e.File, e.Line, e.Message)
	return &jufmt.LogHookResult{WriteOutput: true, Output: line}
})
```

Keep the hook fast; run heavy I/O in a goroutine if needed. jufmt calls the hook synchronously before returning from `Println`.

`SetLogOutputHook(nil)` restores built-in behavior (same as before the hook existed).

## Output

`Output` (default `os.Stdout`) receives the **formatted string** from the hook path above. For simple file redirect without a hook:

```go
f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
jufmt.Output = f
```

For structured persistence, prefer `SetLogOutputHook` over parsing bytes from a custom `io.Writer`.

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
