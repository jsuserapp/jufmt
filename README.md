# jufmt

A partial replacement for the standard `fmt` and `log` packages. Most APIs are drop-in compatible.

Prints optional timestamps and call-site locations (`file:line`) alongside messages so you can jump to source in the IDE.

[中文文档](README.zh.md)

## Call-site tracing

- `Print`, `Println`, `Printf`, and the same methods on `Color`: call-site prefixes are controlled by `SetPrintTrace` and show the **direct caller**.
- `TracePrintln` / `TracePrintf`: always include call-site information; `exStep` adds upstream `[call N]` lines (N=1 is nearest upstream, N=2 is further up). Lines are printed in execution order: `[call exStep]` first, then ... `[call 1]`, then the main message.
- Locations are collected with `runtime.Callers`, skipping jufmt internals, the Go runtime (e.g. `proc.go`), and the standard library (GOROOT). Only user/application frames are kept; no fixed `skip` values.
- `GetTrace(depth)`: `depth=0` is the direct caller, `depth=1` is one level up, and so on.

## Other

- `Output`: redirect output for tests or files; defaults to `os.Stdout`.
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
