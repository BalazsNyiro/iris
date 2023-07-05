package main

// run: go run iris_windows.go
import "fmt"
import "github.com/BalazsNyiro/iris/iris"

func main() {
	fmt.Println("iris example...")
	widthSys, heightSys := iris.TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	ch_data_input := make(chan string)
	go iris.UserInterfaceStart(ch_data_input)

	ch_data_input <- `	obj:win_1
						top:5
						bottom:22
						left:4
						right:33`

	for i := 1; i < 5; i++ {
		ch_data_input <- ` obj:win_1
						addSimpleText:AddEverythingAfter\tfirstColon
This is the second line`
		iris.TimeSleep(2000)
	}

	ch_data_input <- ` 	obj:win_manager
						cmd:exit
`
}
