package main

import (
	"fmt"
)

func main() {
	fmt.Println("example...")
	widthSys, heightSys := TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	ch_data_input := make(chan MessageAndCharactersForWindowsUpdate)
	go UserInterfaceStart(ch_data_input, "\n")

	ch_data_input <- MessageAndCharactersForWindowsUpdate{
		msg: `	select:win:logs-left
						msgId:1
						set:yTop:5
						set:xLeft:4
						set:width:12
						set:height:11
`,
		addLine: LineChars{},
	}

	for i := 1; i < 5; i++ {
		// everything after the 'add:simpleText:' is the part of the text
		ch_data_input <- MessageAndCharactersForWindowsUpdate{
			msg:     `select:win:logs-left`,
			addLine: LineCharsFromStr(fmt.Sprintf("%d sample", i)),
		}
		TimeSleep(2000)
	}

	UserInterfaceExit()
}
