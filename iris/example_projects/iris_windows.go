package main

// run: go run iris_windows.go
import "fmt"
import "github.com/BalazsNyiro/iris/iris"

func main() {
	fmt.Println("iris example...")
	widthSys, heightSys := iris.TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	ch_data_input := make(chan iris.MessageAndCharactersForWindowsUpdate)
	go iris.UserInterfaceStart(ch_data_input)

	// msgId: if the sender wants to get answer,
	// use this id in the reply message
	ch_data_input <- iris.MessageAndCharactersForWindowsUpdate{
		msg: `	select:win:logs-left
						msgId:1
						set:yTop:5
						set:height:22
						set:xLeft:4
						set:width:33`,
		addLine: Line{},
	}

	for i := 1; i < 5; i++ {
		// everything after the 'add:simpleText:' is the part of the text
		ch_data_input <- iris.MessageAndCharactersForWindowsUpdate{
			msg: `select:win:logs-left`,
		}
		// addLine: iris.LineChars  FromStr(fmt.Sprintf("%d sample", i)),
		iris.TimeSleep(2000)
	}
	iris.UserInterfaceExit()
}

}