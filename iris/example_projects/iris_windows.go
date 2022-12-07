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

	windows := iris.WindowsNewState(widthSys, heightSys)
	windows = iris.WinNew(windows, "Child", "0", "0", "8", "4", "C")
	windows["Child"][iris.KeyWinContentSrc] = "apple\norange\nbanana"
	windows["prgState"]["winActiveId"] = "Child"

	iris.UserInterfaceStart(windows)
}
