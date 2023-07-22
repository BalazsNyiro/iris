package main

import (
	"fmt"
)

func main() {
	fmt.Println("example...")
	widthSys, heightSys := TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	ch_data_input := make(chan MessageAndCharacters)
	go UserInterfaceStart(ch_data_input, "\n")

	ch_data_input <- MessageAndCharacters{
		msg: `	select:win:logs-left
						msgId:1
						set:top:5
						set:bottom:22
						set:left:4
						set:right:33`,
		addLine: Line{},
	}

	for i := 1; i < 5; i++ {
		// everything after the 'add:simpleText:' is the part of the text
		ch_data_input <- MessageAndCharacters{
			msg:     `select:win:logs-left`,
			addLine: LineFromStr(fmt.Sprintf("%d sample", i)),
		}
		TimeSleep(2000)
	}

	UserInterfaceExit()
}
