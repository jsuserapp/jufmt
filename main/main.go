package main

import (
	"github.com/jsuserapp/jufmt"
)

func main() {
	jufmt.SetPrintTime(false)
	jufmt.SetPrintTrace(false)
	jufmt.Println("Hello World")
	jufmt.Printf("%s %s\n", "Hello", "World")
	jufmt.Green.Println("Hello World")
	jufmt.Cyan.Printf("%s %s\n", "Hello", "World")

	hello := jufmt.Yellow.Sprintf("%s", "Hello")
	world := jufmt.Red.Sprintf("%s", "World")
	jufmt.Printf("%s %s\n", hello, world)

	traceTest1()
}
func traceTest1() {
	jufmt.Green.TracePrintln(1, "traceTest1")
	traceTest2()
}
func traceTest2() {
	jufmt.Cyan.TracePrintln(2, "traceTest2")
}
