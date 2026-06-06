//go:build dev

package main

import (
	"fmt"
	"runtime"
)

// getCallerFrame 是替代 traceableFrames 的最高效版本。
// 参数 skip 表示你要跳过多少层“业务包装”函数。
func getCallerFrame(skip, extStep int) ([]runtime.Frame, bool) {
	// 【性能优化 1】：只需要定位用户调用的一行，因此只分配长度为 1 的数组。
	// 避免了分配大数组和后续的无效遍历。
	pc := make([]uintptr, 1+extStep)

	// 【性能优化 2】：精准跳过。
	// runtime.Callers 自身的调用为 skip 0
	// getCallerFrame 自身的调用为 skip 1
	// 传入的 skip 则是用来跳过你内部的库函数层级。
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		return nil, false
	}

	// 转换物理 PC 到逻辑 Frame
	calls := runtime.CallersFrames(pc)
	var frames []runtime.Frame
	for {
		frame, more := calls.Next()
		frames = append(frames, frame)
		if !more {
			break
		}
	}

	// 注意：虽然只需要一个 Frame，但依然必须调用 CallersFrames，
	// 否则如果用户的外部调用正好是被内联的函数，直接用 FuncForPC 会导致定位不准。
	return frames, true
}

// ---------------- 以下为模拟你的库内部结构 ----------------

// internalPrint 相当于你的 internal 内部打印逻辑
func internalPrint(extStep int, msg string) {
	// 【核心算数】：
	// 此时的调用栈是：
	// main -> JufmtColorPrintln -> internalPrint -> getCallerFrame
	// 我们在 internalPrint 里调用，需要跳过它自己 (1层) 和 JufmtColorPrintln (1层)
	// 所以传给 getCallerFrame 的 skip 值为 2。
	frames, _ := getCallerFrame(2, extStep)
	for _, frame := range frames {
		fmt.Printf("%s:%d,%s\n", frame.File, frame.Line, frame.Function)
	}

	fmt.Println(">> 实际打印内容:", msg)
}

// JufmtColorPrintln 相当于你的 jufmt.Color.Println，供外部使用的 API
func JufmtColorPrintln(extStep int, msg string) {
	internalPrint(extStep, msg)
}
func JufmtColorPrintTraceln(extStep int, msg string) {
	internalPrint(extStep, msg)
}

// ---------------- 以下为用户侧代码 ----------------

func main() {
	// 这是用户实际写代码的地方，我们希望日志准确指出这里的行号
	JufmtColorPrintln(0, "Hello, Lapintool!")
	extTest()
}
func extTest() {
	JufmtColorPrintTraceln(1, "Hello, Lapintool!")
}

//测试的输出：
//D:/Dev/Projects/Go/jufmt/main/dev.go:69,main.main
//>> 实际打印内容: Hello, Lapintool!
//D:/Dev/Projects/Go/jufmt/main/dev.go:73,main.extTest
//D:/Dev/Projects/Go/jufmt/main/dev.go:70,main.main
//>> 实际打印内容: Hello, Lapintool!
