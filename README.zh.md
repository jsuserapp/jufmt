# jufmt

部分代替官方库的 `fmt` 和 `log`，大部分函数可以通用。

提供打印信息的同时，打印时间和调用位置信息（`file:line`），方便在 IDE 中点击跳转到源码。

[English README](README.md)

## 调用位置（trace）

- `Print` / `Println` / `Printf` 及 `Color` 的同名方法：可通过 `SetPrintTrace` 开关，显示**直接调用方**的 `file:line`（控制台为文件名）。
- `TracePrintln` / `TracePrintf`：主消息始终带调用位置；`exStep` 在其之前额外打印上游帧，每行格式为 `file:line` + **函数名**（短名）。按执行顺序输出（先发生的调用在前）。
- 位置通过 `runtime.Callers` 收集，跳过 jufmt 内部、Go runtime（如 `proc.go`）及标准库（GOROOT）。
- `GetTrace(depth)` 返回 `file:line`（basename）；`GetCallerName(depth)` 返回同一深度的短函数名。
- `SetPrintTrace(false)` 作用于 `Print`/`Println`/`Printf`：**不 walk 栈**，`File`/`Line`/`Func` 为空（与 `SetPrintTime(false)` 不写 `TimeText` 同理）。`TracePrintln`、`TracePrintf`、`Logger` 始终采集位置。

示例（`main` → `traceTest1` → `traceTest2` 内 `TracePrintln(2, "traceTest2")`）：

```
main.go:19 main
main.go:27 traceTest1
main.go:30 traceTest2
```

## LogOutputHook（结构化日志）

`SetLogOutputHook` 在每次打印时收到一条 `LogEntry`。库本身**不涉及数据库**；在 Hook 里用结构化字段自行入库：

| 字段 | 用途 |
|------|------|
| `File` | **全路径**（已采集位置时）；`Print` 系列且 `SetPrintTrace(false)` 时为空 |
| `Line`, `Func` | 调用位置；未采集时为空 |
| `Message` | 纯文本，**无 ANSI**；`LineMain` 为用户内容，`LineUpstream` 为上游函数名 |
| `Time` | 墙钟时间（始终有值） |
| `TimeText` | 格式化时间前缀；`SetPrintTime(false)` 时为空 |
| `Newline` | `true` 表示 Println 风格（`DefaultFormat` 末尾加 `\n`）；Print/Printf 为 `false` |
| `Kind` | `LineMain` = 用户消息行；`LineUpstream` = `TracePrintln` 额外 call-chain 行 |
| `ShowLocation` | 本行是否采集了 `File`/`Line`/`Func` |
| `Color` | 控制台 ANSI 参考 |

返回 `LogHookResult`：

| `WriteOutput` | `Output` | 行为 |
|---------------|----------|------|
| `false` | （忽略） | 不写 `Output`（仅入库，适合服务器） |
| `true` | `""` | 使用 `DefaultFormat(entry)` → 控制台 basename + 颜色 |
| `true` | 自定义字符串 | 按你的格式写入 `Output` |

**仅入库（不打印控制台）**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) jufmt.LogHookResult {
	go func() { // 异步入库；顺序由业务自行保证
		db.Exec(`INSERT INTO logs(at,file,line,func,msg) VALUES (?,?,?,?,?)`,
			e.Time, e.File, e.Line, e.Func, e.Message)
	}()
	return jufmt.LogHookResult{WriteOutput: false}
})
```

**入库 + 默认控制台**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) jufmt.LogHookResult {
	go dbInsert(e) // e.File 为全路径
	return jufmt.LogHookResult{WriteOutput: true} // Output 空 → DefaultFormat
})
```

**自定义控制台格式**

```go
jufmt.SetLogOutputHook(func(e jufmt.LogEntry) jufmt.LogHookResult {
	persist(e)
	line := fmt.Sprintf("%s %s:%d %s\n", e.TimeText, e.File, e.Line, e.Message)
	return jufmt.LogHookResult{WriteOutput: true, Output: line}
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
