# jufmt

部分代替官方库的 `fmt` 和 `log`，大部分函数可以通用。

提供打印信息的同时，打印时间和调用位置信息（`file:line`），方便在 IDE 中点击跳转到源码。

[English README](README.md)

## 调用位置（trace）

- `Print` / `Println` / `Printf` 及 `Color` 的同名方法：可通过 `SetPrintTrace` 开关，显示**直接调用方**的 `file:line`（控制台为文件名）。
- `TracePrintln` / `TracePrintf`：主消息为普通文本；上游行为 `file:line` + **`[函数名]`**，按执行顺序输出。
- 位置通过 `runtime.Callers` 采集，由各 API 传入精准 `skip` 与帧数 `exStep`（`traceableFrames(skip, exStep)`）；上游行仍用 `isTraceableFrame` 跳过 runtime/stdlib。
- `GetTrace(depth)` 返回 `file:line`（basename）；`GetCallerName(depth)` 返回同一深度的短函数名。
- `SetPrintTrace(false)` 作用于 `Print`/`Println`/`Printf`：**不 walk 栈**，`File`/`Line`/`Func` 为空（与 `SetPrintTime(false)` 不写 `TimeText` 同理）。`TracePrintln`、`TracePrintf`、`Logger` 始终采集位置。

示例（`main` → `traceTest1` → `traceTest2` 内 `TracePrintln(2, "traceTest2")`）：

```
main.go:19 [main]
main.go:27 [traceTest1]
main.go:30 traceTest2
```

## LogOutputHook（结构化日志）

`SetLogOutputHook` 在每次打印时收到一条 `LogEntry`。库本身**不涉及数据库**；在 Hook 里用结构化字段自行入库：

| 字段 | 用途 |
|------|------|
| `File` | **runtime 源路径**（通常为绝对路径）；未采集位置时为空 |
| `Line`, `Func` | 调用位置；未采集时为空 |
| `Message` | 纯文本，**无 ANSI**；`LineMain` 为用户内容，`LineUpstream` 为 `[函数名]` |
| `Time` | 墙钟时间（始终有值） |
| `TimeText` | 格式化时间前缀；`SetPrintTime(false)` 时为空 |
| `Newline` | `true` 表示 Println 风格（`DefaultFormat` 末尾加 `\n`）；Print/Printf 为 `false` |
| `Kind` | `LineMain` = 用户消息行；`LineUpstream` = `TracePrintln` 额外 call-chain 行 |
| `ShowLocation` | 本行是否采集了 `File`/`Line`/`Func` |
| `Color` | 控制台 ANSI 参考 |

返回 `*LogHookResult`（多数情况直接 **`return nil`** 即可，保持默认控制台输出）：

| 返回值 | 行为 |
|--------|------|
| **`nil`** | 写入 `DefaultFormat(entry)` — 仅需回调、不改控制台格式时使用 |
| `&LogHookResult{WriteOutput: false}` | 不写 `Output`（仅入库） |
| `&LogHookResult{WriteOutput: true}` | 同 `nil`（`DefaultFormat`） |
| `&LogHookResult{WriteOutput: true, Output: "..."}` | 自定义字符串写入 `Output` |

**仅回调（默认控制台）**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) *jufmt.LogHookResult {
	go dbInsert(e) // e.File 为模块相对路径，e.Message 无 ANSI
	return nil
})
```

**仅入库（不打印控制台）**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) *jufmt.LogHookResult {
	go dbInsert(e)
	return &jufmt.LogHookResult{WriteOutput: false}
})
```

**自定义控制台格式**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) *jufmt.LogHookResult {
	persist(e)
	line := fmt.Sprintf("%s %s:%d %s\n", e.TimeText, e.File, e.Line, e.Message)
	return &jufmt.LogHookResult{WriteOutput: true, Output: line}
})
```

Hook 应尽量轻量；耗时 I/O 可在 goroutine 中异步执行。jufmt 在 `Println` 返回前**同步**调用 Hook。

`SetLogOutputHook(nil)` 恢复内置行为（与未设置 Hook 时一致）。

## Output

`Output`（默认 `os.Stdout`）接收 Hook 决定写入的**格式化字符串**。无 Hook 时可直接重定向到文件：

```go
f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
jufmt.Output = f
```

结构化入库优先使用 `SetLogOutputHook`，而不是在 `io.Writer` 里解析字节。

## 其他

- `Log`：包级默认 `Logger` 实例。

因为自动添加信息前缀，通常每次打印一行，不能像常规 fmt 不换行可以连续打印信息，否则直接连接的带前缀的信息会比较难看。

## 示例

```go
jufmt.SetPrintTime(true)
jufmt.SetPrintTrace(true)

jufmt.Println("hello")
jufmt.Green.Println("success")
jufmt.Log.Info("ready")

jufmt.Green.TracePrintln(2, "nested trace")
```

## 性能测试

与原生 `fmt.Fprint*`（同样写入 `io.Discard`）对比 jufmt 各 API 的开销。在模块根目录执行：

```bash
# 全部 benchmark，含内存分配
go test -bench=. -benchmem ./...

# 只对比 Println
go test -bench='BenchmarkFmt_Fprintln|BenchmarkJufmt_Println' -benchmem ./...

# 更长采样，结果更稳定
go test -bench=. -benchmem -benchtime=3s -count=5 ./...
```

| Benchmark | 含义 |
|-----------|------|
| `BenchmarkFmt_Fprint*` | 标准库基线（`io.Discard`） |
| `BenchmarkJufmt_Print*` | 无时间戳，`SetPrintTrace(false)` |
| `BenchmarkJufmt_Println_Trace` | `SetPrintTrace(true)`，采集一层栈 |
| `BenchmarkJufmt_Println_TimeAndTrace` | 时间戳 + 调用位置 |
| `BenchmarkJufmt_TracePrintln/exStepN` | 强制 trace；`exStep` 条上游 `[函数名]` 行 |

### 栈开销拆分（Callers vs CallersFrames）

与 `Println` + trace 相同 `callersSkip` 的包内 micro-benchmark：

| Benchmark | 含义 |
|-----------|------|
| `BenchmarkStack_Empty` | 空循环基线 |
| `BenchmarkStack_CallersOnly` | 仅 `runtime.Callers` |
| `BenchmarkStack_CallersFrames/firstFrame` | Callers + 一次 `Next()` |
| `BenchmarkStack_CallersFrames/allFrames` | Callers + 对所有 PC 做 `Next()` |
| `BenchmarkStack_TraceableFrames` | 完整 `traceableFrames` |

```bash
go test -bench=BenchmarkStack -benchmem -run=^$ .

# 对比栈拆分 vs 完整 jufmt Println+trace
go test -bench='BenchmarkStack|BenchmarkJufmt_Println_Trace' -benchmem -run=^$ .
```

jufmt 会附加 ANSI 颜色并可选 walk 栈；高频 `Print*` 若不需要位置前缀，请 `SetPrintTrace(false)`。
