package main

import (
	"github.com/jsuserapp/jufmt"
)

func main() {
	jufmt.Println("Hello World")
	jufmt.Printf("%s %s\n", "Hello", "World")
	jufmt.Green.Println("Hello World")
	jufmt.Cyan.Printf("%s %s\n", "Hello", "World")

	hello := jufmt.Yellow.Sprintf("%s", "Hello")
	world := jufmt.Red.Sprintf("%s", "World")
	jufmt.Printf("%s %s\n", hello, world)
}
