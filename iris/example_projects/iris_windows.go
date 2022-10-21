package main

// run: go run iris_windows.go

import "fmt"
import "github.com/BalazsNyiro/iris/iris"

// run: go run jyp_example_usage.go

func main() {
    fmt.Println(iris.TerminalDetect())
    fmt.Println(iris.TerminalDimensionsWithSyscall())
}
