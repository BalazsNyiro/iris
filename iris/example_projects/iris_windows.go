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

	ch_data_input <- `	select:win:logs-left
						set:top:5
						set:bottom:22
						set:left:4
						set:right:33`

	for i := 1; i < 5; i++ {
		// everything after the 'addSimpleText:' is the part of the text
		ch_data_input <- `select:win:logs-left
						add:simpleText:AddEverythingAfterColon` + fmt.Sprintf("%d", i)
		iris.TimeSleep(2000)
	}

}
