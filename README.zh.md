# jufmt

部分代替官方库的 `fmt` 和 `log`，大部分函数可以通用。

提供打印信息的同时，打印时间和调用位置信息（`file:line`），方便在 IDE 中点击跳转到源码。

[English README](README.md)

## 调用位置（trace）

- `Print` / `Println` / `Printf` 及 `Color` 的同名方法：可通过 `SetPrintTrace` 开关，显示**直接调用方**的 `file:line`。
- `TracePrintln` / `TracePrintf`：主消息始终带调用位置；`exStep` 在其之前额外打印上游帧，每行格式为 `file:line` + **函数名**（短名）。按执行顺序输出（先发生的调用在前）。
- 位置通过 `runtime.Callers` 收集，跳过 jufmt 内部、Go runtime（如 `proc.go`）及标准库（GOROOT）。
- `GetTrace(depth)` 返回 `file:line`；`GetCallerName(depth)` 返回同一深度的短函数名。

示例（`main` → `traceTest1` → `traceTest2` 内 `TracePrintln(2, "traceTest2")`）：

```
main.go:19 main
main.go:27 traceTest1
main.go:30 traceTest2
```

## 重定向输出

`Output` 是 `io.Writer`（默认 `os.Stdout`），赋值为任意 Writer 即可把输出写到别处。

**写入日志文件**

```go
f, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
if err != nil { /* handle */ }
defer f.Close()
jufmt.Output = f
```

**同时输出到终端和文件**

```go
jufmt.Output = io.MultiWriter(os.Stdout, f)
```

**测试捕获**

```go
var buf bytes.Buffer
jufmt.Output = &buf
jufmt.Println("capture me")
_ = buf.String()
```

**数据库或远程日志**

实现 `io.Writer`，在 `Write` 里把字节写入存储：

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

说明：

- jufmt 输出的是**纯文本**（使用 `Color` 时会带 ANSI 转义码），自定义 Writer 收到的是字节块，不是结构化字段。
- 入库时建议以 `Println` / `TracePrintln` 等「一行一次」的调用为单位；少依赖多次 `Print` 拼一行。
- 若下游不是终端，可先去掉 ANSI 再存。
- 写文件/数据库时，常配合 `SetPrintTime(false)` 并少用亮色。

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
