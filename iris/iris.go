/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package iris

import (
	"fmt"
)

func UserInterfaceStart(ch_data_input chan MessageAndCharactersForWindowsUpdate, dataInputLineSeparator string) {
	userInterfaceInit()
	ch_user_input := make(chan string)
	go channelRead_userInput(ch_user_input)

	ch_terminal_size_change_detect := make(chan [2]int)
	go channelRead_terminalSizeChangeDetect(ch_terminal_size_change_detect)

	// windows is a read-only variable everywhere,
	windows := Windows{} // modified/updated ONLY here:
	go channelRead_dataInputInterpret(ch_data_input, &windows, dataInputLineSeparator)

	widthSysNow, heightSysNow := TerminalDimensionsWithSyscall()
	terminalSizeActual := [2]int{widthSysNow, heightSysNow}

	loopCounter := 0
	for {
		loopCounter++

		action := ""
		select { //                https://gobyexample.com/select
		case stdin, _ := <-ch_user_input: //  the message is coming...
			fmt.Println("Keys pressed:", stdin)
			action = action_of_user_input(stdin)

		case terminal_size_change, _ := <-ch_terminal_size_change_detect: //  the message is coming...
			terminalSizeActual = terminal_size_change
		default: //               or not coming
			_ = ""
		}

		if action == "quit" {
			UserInterfaceExit()
			break
		}
		fmt.Println("windows: ", windows)
		layers := ScreenLayersRenderFromWindows(windows, terminalSizeActual)
		// fmt.Println("layers:", layers)
		ScreenLayersDisplayAll(layers, dataInputLineSeparator, loopCounter)
		TimeSleep(TimeIntervalUserInterfaceRefreshTimeMillisec)
	}
}

func userInterfaceInit() {
	terminal_console_clear()
	terminal_console_input_buffering_disable()
	terminal_console_character_hide()
}

func UserInterfaceExit() {
	terminal_console_character_show()
	terminal_console_input_buffering_enable()
}
