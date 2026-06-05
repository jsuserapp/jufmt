# jufmt

部分代替官方库的 `fmt` 和 `log`，大部分函数可以通用。

提供打印信息的同时，打印时间和调用位置信息（`file:line`），方便在 IDE 中点击跳转到源码。

[English README](README.md)

## 调用位置（trace）

- `Print` / `Println` / `Printf` 及 `Color` 的同名方法：可通过 `SetPrintTrace` 开关，显示**直接调用方**的位置。
- `TracePrintln` / `TracePrintf`：始终显示调用位置；`exStep` 额外打印上游 `[call N]` 行（N=1 为最近上游，N=2 再远一层）。按执行顺序输出：先 `[call exStep]`，再 … `[call 1]`，最后主消息。
- 位置通过 `runtime.Callers` 收集栈帧，跳过 jufmt 内部、Go runtime（如 `proc.go`）及标准库（GOROOT）；只保留用户/业务代码帧，不依赖固定 `skip`。
- `GetTrace(depth)`：`depth=0` 为直接调用方，`depth=1` 为再上一层，以此类推。

## 其他

- `Output`：可重定向输出（测试或写入文件），默认 `os.Stdout`。
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
