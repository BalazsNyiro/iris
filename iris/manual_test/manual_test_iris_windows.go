package main

import (
	"fmt"
)

func main() {
	fmt.Println("example...")
	widthSys, heightSys := TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	ch_data_input := make(chan MessageAndCharactersForWindowsUpdate)
	go UserInterfaceStartThread(ch_data_input, "\n")

	ch_data_input <- MessageAndCharactersForWindowsUpdate{
		msg: `	select:textBlock:logs-left
						msgId:1
`,
		addLine: LineChars{},
	}

	for i := 0; i < 8; i++ {
		// everything after the 'add:simpleText:' is the part of the text
		ch_data_input <- MessageAndCharactersForWindowsUpdate{
			msg:     `select:textBlock:logs-left`,
			addLine: LineCharsFromStr(fmt.Sprintf("%d sample", i), i),
		}
		TimeSleep(2000)
	}

	UserInterfaceExitThread()
}
