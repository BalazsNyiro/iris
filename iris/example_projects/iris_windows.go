package main

// run: go run iris_windows.go

import "fmt"
import "github.com/BalazsNyiro/iris/iris"

// run: go run jyp_example_usage.go

func main() {
	widthStty, heightStty := iris.TerminalDimensionsSttySize()
	fmt.Println("stty:", widthStty, heightStty)
	widthSys, heightSys := iris.TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	iris.UserInterfaceStart()
}
