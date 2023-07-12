package main

import (
	"fmt"
)

func main() {
	fmt.Println("example...")
	widthSys, heightSys := TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	ch_data_input := make(chan string)
	go UserInterfaceStart(ch_data_input, "\n")

	// useThisIdInReply: if the sender wants to get answer,
	// use this id in the reply message
	ch_data_input <- `	select:win:logs-left
						useThisIdInReply:1
						set:xLeft:5
						set:yTop:5
						set:width:22
						set:height:10
    `

	for i := 1; i < 5; i++ {
		// everything after the 'add:simpleText:' is the part of the text
		ch_data_input <- `select:win:logs-left
						add:simpleText:AddEverythingAfterColon` + fmt.Sprintf("%d\n", i)
		TimeSleep(2000)
	}
	UserInterfaceExit()
}
