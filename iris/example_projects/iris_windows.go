package main

// run: go run iris_windows.go
import "fmt"
import "github.com/BalazsNyiro/iris/iris"

func main() {
	fmt.Println("iris example...")
	widthSys, heightSys := iris.TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	iris.UserInterfaceStart()
}
